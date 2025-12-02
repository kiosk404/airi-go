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
	Identity *entity.AgentIdentity
	CPStore  compose.CheckPointStore

	CustomVariables map[string]string
	ConversationID  int64
}

const (
	keyOfPersonRender           = "persona_render"
	keyOfKnowledgeRetriever     = "knowledge_retriever"
	keyOfKnowledgeRetrieverPack = "knowledge_retriever_pack"
	keyOfPromptVariables        = "prompt_variables"
	keyOfPromptTemplate         = "prompt_template"
	keyOfReActAgent             = "react_agent"
	keyOfReActAgentToolsNode    = "agent_tool"
	keyOfReActAgentChatModel    = "re_act_chat_model"
	keyOfLLM                    = "llm"
	keyOfToolsPreRetriever      = "tools_pre_retriever"
)

func BuildAgent(ctx context.Context, conf *Config) (r *AgentRunner, err error) {
	persona := conf.Agent.Prompt.GetPrompt()
	// 将用户信息编辑进 上下文变量中
	avConf := &variableConf{
		Agent:  conf.Agent,
		UserID: conf.UserID,
	}
	avs, err := loadAgentVariables(ctx, avConf)
	if err != nil {
		return nil, err
	}
	if conf.CustomVariables != nil {
		for k, v := range conf.CustomVariables {
			avs[k] = v
		}
	}

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

	var isReActAgent bool
	if len(agentTools) > 0 {
		isReActAgent = true
		requireCheckpoint = true
		if modelInfo.Capability != nil && !modelInfo.Capability.GetFunctionCall() {
			return nil, fmt.Errorf("model %v does not support function call", modelInfo.DisplayInfo.GetName())
		}
	}

	var agentGraph compose.AnyGraph
	var agentNodeOpts []compose.GraphAddNodeOpt
	var agentNodeName string
	if isReActAgent {
		agent, err := react.NewAgent(ctx, &react.AgentConfig{
			ToolCallingModel: chatModel,
			ToolsConfig: compose.ToolsNodeConfig{
				Tools: agentTools,
			},
			ToolReturnDirectly: returnDirectlyTools,
			ModelNodeName:      keyOfReActAgentChatModel,
			ToolsNodeName:      keyOfReActAgentToolsNode,
		})
		if err != nil {
			return nil, err
		}
		agentGraph, agentNodeOpts = agent.ExportGraph()

		agentNodeName = keyOfReActAgent
	} else {
		agentNodeName = keyOfLLM
	}

	// 生成问题建议，比如在回答完成任务后，建议生成后续的问题
	suggestGraph, nsg := newSuggestGraph(ctx, conf, chatModel)

	g := compose.NewGraph[*AgentRequest, *schema.Message](
		compose.WithGenLocalState(func(ctx context.Context) (state *AgentState) {
			return &AgentState{}
		}))

	_ = g.AddLambdaNode(keyOfPersonRender,
		compose.InvokableLambda[*AgentRequest, string](personaVars.RenderPersona),
		compose.WithStatePreHandler(func(ctx context.Context, ar *AgentRequest, state *AgentState) (*AgentRequest, error) {
			state.UserInput = ar.Input
			return ar, nil
		}),
		compose.WithOutputKey(placeholderOfPersona))

	_ = g.AddLambdaNode(keyOfPromptVariables,
		compose.InvokableLambda[*AgentRequest, map[string]any](promptVars.AssemblePromptVariables))

	// 知识库检索
	_ = g.AddLambdaNode(keyOfKnowledgeRetriever,
		compose.InvokableLambda[*AgentRequest, []*schema.Document](kr.Retrieve),
		compose.WithNodeName(keyOfKnowledgeRetriever))

	_ = g.AddLambdaNode(keyOfToolsPreRetriever,
		compose.InvokableLambda[*AgentRequest, []*schema.Message](tr.toolPreRetrieve),
		compose.WithOutputKey(keyOfToolsPreRetriever),
		compose.WithNodeName(keyOfToolsPreRetriever),
	)
	_ = g.AddLambdaNode(keyOfKnowledgeRetrieverPack,
		compose.InvokableLambda[[]*schema.Document, string](kr.PackRetrieveResultInfo),
		compose.WithOutputKey(placeholderOfKnowledge),
	)
	_ = g.AddChatTemplateNode(keyOfPromptTemplate, chatPrompt)

	agentNodeOpts = append(agentNodeOpts, compose.WithNodeName(agentNodeName))

	if isReActAgent {
		_ = g.AddGraphNode(agentNodeName, agentGraph, agentNodeOpts...)
	} else {
		_ = g.AddChatModelNode(agentNodeName, chatModel, agentNodeOpts...)
	}

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
