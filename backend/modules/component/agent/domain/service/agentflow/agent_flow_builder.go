package agentflow

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/llm/application"
	"github.com/kiosk404/airi-go/backend/pkg/lang/maps"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
)

type Config struct {
	Agent    *entity.SingleAgent
	UserID   string
	Identity *entity.AgentIdentity   // 身份标识，用于工具调用时的身份识别
	CPStore  compose.CheckPointStore // 用于 ReAct Agent 的状态持久化，支持断点续传

	CustomVariables map[string]string // 用户自定义变量，会被注入到提示词模板中
	ConversationID  int64
}

const (
	// keyOfPersonRender 人格/角色 渲染节点，负责将变量注入到 persona 提示词中
	keyOfPersonRender = "persona_render"
	// keyOfKnowledgeRetriever 知识库检索节点
	keyOfKnowledgeRetriever = "knowledge_retriever"
	// keyOfKnowledgeRetrieverPack 知识库检索节点，用于将检索到的知识库文档打包到上下文变量中
	keyOfKnowledgeRetrieverPack = "knowledge_retriever_pack"
	// keyOfPromptVariables 提示词变量组装节点
	keyOfPromptVariables = "prompt_variables"
	// keyOfPromptTemplate 提示词模板节点
	keyOfPromptTemplate      = "prompt_template"
	keyOfReActAgent          = "react_agent"
	keyOfReActAgentToolsNode = "agent_tool"
	keyOfReActAgentChatModel = "re_act_chat_model"
	keyOfLLM                 = "llm"
	keyOfToolsPreRetriever   = "tools_pre_retriever"
)

// BuildAgent 是 Agent 构建的核心函数，负责组装完整的 Agent 执行图。
//
// 该函数实现了一个基于 DAG (有向无环图) 的 Agent 执行流程，主要包含以下步骤：
//  1. 加载并处理变量（用户变量、自定义变量）
//  2. 构建人格/角色渲染器（处理 Jinja2 风格的模板变量）
//  3. 初始化知识库检索器
//  4. 构建 LLM 聊天模型
//  5. 加载各类工具（插件工具、数据库工具、变量工具）
//  6. 根据工具配置决定使用 ReAct Agent 还是普通 LLM
//  7. 构建并编译执行图
//
// 执行图的拓扑结构如下：
//
//	                 ┌─────────────────────────────────────────────┐
//	                 │                   START                     │
//	                 └─────────────────────────────────────────────┘
//	                                      │
//	      ┌───────────────┬───────────────┼───────────────┬───────────────┐
//	      ▼               ▼               ▼               ▼               │
//	┌───────────┐  ┌───────────┐  ┌───────────────┐  ┌─────────────┐      │
//	│  persona  │  │  prompt   │  │  knowledge    │  │ tools_pre   │      │
//	│  render   │  │ variables │  │  retriever    │  │ retriever   │      │
//	└───────────┘  └───────────┘  └───────────────┘  └─────────────┘      │
//	      │               │               │                   │           │
//	      │               │               ▼                   │           │
//	      │               │        ┌───────────────┐          │           │
//	      │               │        │  knowledge    │          │           │
//	      │               │        │retriever_pack │          │           │
//	      │               │        └───────────────┘          │           │
//	      │               │               │                   │           │
//	      └───────────────┴───────────────┼───────────────────┘           │
//	                                      ▼                               │
//	                            ┌───────────────────┐                     │
//	                            │  prompt_template  │                     │
//	                            └───────────────────┘                     │
//	                                      │                               │
//	                                      ▼                               │
//	                        ┌───────────────────────────┐                 │
//	                        │  ReAct Agent 或 LLM 节点   │                 │
//	                        │  (根据是否有工具决定)       │                 │
//	                        └───────────────────────────┘                 │
//	                                      │                               │
//	                  ┌───────────────────┴───────────────────┐           │
//	                  │         (如果启用建议功能)             │           │
//	                  ▼                                       ▼           │
//	       ┌───────────────────┐                     ┌─────────────┐      │
//	       │ suggest_pre_parse │                     │     END     │◄─────┘
//	       └───────────────────┘                     └─────────────┘
//	                  │
//	                  ▼
//	       ┌───────────────────┐
//	       │   suggest_graph   │
//	       └───────────────────┘
//	                  │
//	                  ▼
//	            ┌───────────┐
//	            │    END    │
//	            └───────────┘
//
// 参数：
//   - ctx: 上下文，用于控制超时和取消
//   - conf: Agent 配置，包含模型、工具、知识库等所有必要配置
//
// 返回值：
//   - r: AgentRunner 实例，可用于执行 Agent
//   - err: 构建过程中的错误
func BuildAgent(ctx context.Context, conf *Config) (r *AgentRunner, err error) {
	persona := conf.Agent.Prompt.GetPrompt()
	// 获取 Agent 的人格/角色提示词（包含变量）将用户信息编辑进 上下文变量中
	avConf := &variableConf{
		Agent:  conf.Agent,
		UserID: conf.UserID,
	}
	// 从配置中加载 Agent 变量 （如系统预定义的用户信息）
	avs, err := loadAgentVariables(ctx, avConf)
	if err != nil {
		return nil, err
	}
	if conf.CustomVariables != nil {
		for k, v := range conf.CustomVariables {
			avs[k] = v
		}
	}

	// 创建提示词变量组装器
	promptVars := &promptVariables{
		Agent: conf.Agent,
		avs:   avs,
	}

	// 用户的提示词
	personaVars := &personaRender{
		personaVariableNames: extractJinja2Placeholder(persona),
		persona:              persona,
		variables:            avs,
	}

	// 加载知识库
	kr, err := newKnowledgeRetriever(ctx, &retrieverConfig{
		knowledgeConfig: conf.Agent.Knowledge,
	})
	if err != nil {
		return nil, err
	}

	// 生成LLM模型 (聊天模型)
	chatModel, modelInfo, err := application.BuildModelBySettings(ctx, conf.Agent.ModelInfo)
	if err != nil {
		return nil, err
	}
	requireCheckpoint := false
	pluginTools, err := newPluginTools(ctx, &toolConfig{
		userID:        conf.UserID,
		agentIdentity: conf.Identity,
		toolConf:      conf.Agent.Plugin,

		conversationID: conf.ConversationID,
	})
	if err != nil {
		return nil, err
	}
	// 预处理词加载注入到最终的提示词
	tr := newPreToolRetriever(&toolPreCallConf{})

	// todo: 加载工作流 workflow
	returnDirectlyToolSets := mapset.NewSet[string]()
	returnDirectlyTools := maps.MapFromSet(returnDirectlyToolSets)

	// 加载数据库工具
	// 允许 Agent 直接查询和操作数据库
	var dbTools []tool.InvokableTool
	if len(conf.Agent.Database) > 0 {
		dbTools, err = newDatabaseTools(ctx, &databaseConfig{
			userID:        conf.UserID,
			agentIdentity: conf.Identity,
			databaseConf:  conf.Agent.Database,
		})
		if err != nil {
			return nil, err
		}
	}

	var avTools []tool.InvokableTool
	if len(avs) > 0 {
		avTools, err = newAgentVariableTools(ctx, avConf)
		if err != nil {
			return nil, err
		}
	}
	containWfTool := false

	agentTools := make([]tool.BaseTool, 0, len(pluginTools)+len(dbTools)+len(avTools))
	agentTools = append(agentTools, slices.Transform(pluginTools, func(a tool.InvokableTool) tool.BaseTool {
		return a
	})...)
	agentTools = append(agentTools, slices.Transform(dbTools, func(a tool.InvokableTool) tool.BaseTool {
		return a
	})...)
	agentTools = append(agentTools, slices.Transform(avTools, func(a tool.InvokableTool) tool.BaseTool {
		return a
	})...)

	// 根据是否有可用的工具决定使用 ReAct Agent 还是普通 LLM
	// 如果有工具，则使用 ReAct Agent，否则使用普通 LLM
	var isReActAgent bool
	if len(agentTools) > 0 {
		isReActAgent = true
		requireCheckpoint = true

		// 验证模型是否支持函数调用
		// 如果模型不支持函数调用，则返回错误
		if modelInfo.Capability != nil && !modelInfo.Capability.GetFunctionCall() {
			return nil, fmt.Errorf("model %v does not support function call", modelInfo.DisplayInfo.GetName())
		}
	}

	var agentGraph compose.AnyGraph
	var agentNodeOpts []compose.GraphAddNodeOpt
	var agentNodeName string
	if isReActAgent {
		// 构建 ReAct Agent 节点
		// ReAct (Reasoning and Action) Agent 是一个基于 ReAct 框架的 Agent，它可以通过工具来实现复杂的任务
		// 它会执行： 思考 -> 决定行动 -> 执行工具 -> 观察结果 -> 继续思考
		agent, err := react.NewAgent(ctx, &react.AgentConfig{
			ToolCallingModel: chatModel,
			ToolsConfig: compose.ToolsNodeConfig{
				Tools: agentTools,
			},
			// ToolReturnDirectly 指定哪些工具可以直接返回结果，不经过 ReAct 思考和行动
			ToolReturnDirectly: returnDirectlyTools,
			// ModelNodeName 和 ToolsNodeName 分别指定模型节点和工具节点的名称
			ModelNodeName: keyOfReActAgentChatModel,
			ToolsNodeName: keyOfReActAgentToolsNode,
		})
		if err != nil {
			return nil, err
		}
		// 导出 ReAct Agent 的图形和节点选项，将其作为子图嵌入主图
		agentGraph, agentNodeOpts = agent.ExportGraph()
		agentNodeName = keyOfReActAgent
	} else {
		// 无工具时，使用普通 LLM 模型
		agentNodeName = keyOfLLM
	}

	// 生成问题建议，比如在回答完成任务后，建议生成后续的问题
	suggestGraph, nsg := newSuggestGraph(ctx, conf, chatModel)

	// 创建主执行图，输入类型为 AgentRequest， 输出类型为 schema.Message
	g := compose.NewGraph[*AgentRequest, *schema.Message](
		compose.WithGenLocalState(func(ctx context.Context) (state *AgentState) {
			return &AgentState{}
		}))

	// 渲染人格
	_ = g.AddLambdaNode(keyOfPersonRender,
		compose.InvokableLambda[*AgentRequest, string](personaVars.RenderPersona),
		compose.WithStatePreHandler(func(ctx context.Context, ar *AgentRequest, state *AgentState) (*AgentRequest, error) {
			state.UserInput = ar.Input
			return ar, nil
		}),
		compose.WithOutputKey(placeholderOfPersona))

	// 收集并组装所有提示词变量 (如对话历史，系统信息等)
	_ = g.AddLambdaNode(keyOfPromptVariables,
		compose.InvokableLambda[*AgentRequest, map[string]any](promptVars.AssemblePromptVariables))

	// 知识库检索节点 (根据用户输入检索知识库文档，为 LLM 提供背景知识)
	_ = g.AddLambdaNode(keyOfKnowledgeRetriever,
		compose.InvokableLambda[*AgentRequest, []*schema.Document](kr.Retrieve),
		compose.WithNodeName(keyOfKnowledgeRetriever))

	// 工具预检索节点 (加载工具相关信息，注入到提示词中)
	_ = g.AddLambdaNode(keyOfToolsPreRetriever,
		compose.InvokableLambda[*AgentRequest, []*schema.Message](tr.toolPreRetrieve),
		compose.WithOutputKey(keyOfToolsPreRetriever),
		compose.WithNodeName(keyOfToolsPreRetriever),
	)
	// 知识库文档打包节点 (将知识库文档打包为字符串，为 LLM 提供背景知识)
	_ = g.AddLambdaNode(keyOfKnowledgeRetrieverPack,
		compose.InvokableLambda[[]*schema.Document, string](kr.PackRetrieveResultInfo),
		compose.WithOutputKey(placeholderOfKnowledge),
	)
	// 提示词模板节点 (根据变量生成提示词模板)
	_ = g.AddChatTemplateNode(keyOfPromptTemplate, chatPrompt)

	agentNodeOpts = append(agentNodeOpts, compose.WithNodeName(agentNodeName))

	if isReActAgent {
		_ = g.AddGraphNode(agentNodeName, agentGraph, agentNodeOpts...)
	} else {
		_ = g.AddChatModelNode(agentNodeName, chatModel, agentNodeOpts...)
	}

	// 添加问题的后续建议
	if nsg {
		_ = g.AddLambdaNode(keyOfSuggestPreInputParse, compose.ToList[*schema.Message](),
			compose.WithStatePostHandler(func(ctx context.Context, out []*schema.Message, state *AgentState) ([]*schema.Message, error) {
				out = append(out, state.UserInput)
				return out, nil
			}),
		)
		_ = g.AddGraphNode(keyOfSuggestGraph, suggestGraph)
	}

	_ = g.AddEdge(compose.START, keyOfPersonRender)
	_ = g.AddEdge(compose.START, keyOfPromptVariables)
	_ = g.AddEdge(compose.START, keyOfKnowledgeRetriever)
	_ = g.AddEdge(compose.START, keyOfToolsPreRetriever)

	_ = g.AddEdge(keyOfPersonRender, keyOfPromptTemplate)
	_ = g.AddEdge(keyOfPromptVariables, keyOfPromptTemplate)
	_ = g.AddEdge(keyOfKnowledgeRetriever, keyOfKnowledgeRetrieverPack)
	_ = g.AddEdge(keyOfKnowledgeRetrieverPack, keyOfPromptTemplate)
	_ = g.AddEdge(keyOfToolsPreRetriever, keyOfPromptTemplate)

	// 提示词模板 -> Agent 执行
	_ = g.AddEdge(keyOfPromptTemplate, agentNodeName)

	if nsg {
		_ = g.AddEdge(agentNodeName, keyOfSuggestPreInputParse)
		_ = g.AddEdge(keyOfSuggestPreInputParse, keyOfSuggestGraph)
		_ = g.AddEdge(keyOfSuggestGraph, compose.END)
	} else {
		_ = g.AddEdge(agentNodeName, compose.END)
	}

	var opts []compose.GraphCompileOption
	if requireCheckpoint {
		opts = append(opts, compose.WithCheckPointStore(conf.CPStore))
	}
	opts = append(opts, compose.WithNodeTriggerMode(compose.AllPredecessor))
	runner, err := g.Compile(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &AgentRunner{
		runner:              runner,
		requireCheckpoint:   requireCheckpoint,
		modelInfo:           modelInfo,
		containWfTool:       containWfTool,
		returnDirectlyTools: returnDirectlyToolSets,
	}, nil
}

// extractJinja2Placeholder 提取Jinja2占位符的变量，返回变量名列表
func extractJinja2Placeholder(persona string) (variableNames []string) {
	re := regexp.MustCompile(`{{([^}]*)}}`)
	matches := re.FindAllStringSubmatch(persona, -1)
	variables := make([]string, 0, len(matches))
	for _, match := range matches {
		val := strings.TrimSpace(match[1])
		if val != "" {
			variables = append(variables, match[1])
		}
	}
	return variables
}
