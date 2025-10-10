package agentflow

import (
	"context"
	"github.com/cloudwego/eino/compose"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	llm "github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llmimpl/chatmodel"
)

type Config struct {
	Agent        *entity.SingleAgent
	UserID       string
	Identity     *entity.AgentIdentity
	ModelMgr     llm.IManage
	ModelFactory chatmodel.Factory
	CPStore      compose.CheckPointStore

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
	panic("BuildAgent not implemented")
}
