package singleagent

import (
	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/api/crossdomain/agentrun"
	"github.com/kiosk404/airi-go/backend/api/crossdomain/plugin"
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	"gorm.io/gorm"
)

type EventType string

const (
	EventTypeOfChatModelAnswer        EventType = "chatmodel_answer"
	EventTypeOfToolsAsChatModelStream EventType = "tools_as_chatmodel_answer"
	EventTypeOfToolMidAnswer          EventType = "tool_mid_answer"
	EventTypeOfToolsMessage           EventType = "tools_message"
	EventTypeOfFuncCall               EventType = "func_call"
	EventTypeOfSuggest                EventType = "suggest"
	EventTypeOfKnowledge              EventType = "knowledge"
	EventTypeOfInterrupt              EventType = "interrupt"
)

type AgentEvent struct {
	EventType EventType

	ToolMidAnswer         *schema.StreamReader[*schema.Message]
	ToolAsChatModelAnswer *schema.StreamReader[*schema.Message]

	ChatModelAnswer *schema.StreamReader[*schema.Message]
	ToolsMessage    []*schema.Message
	FuncCall        *schema.Message
	Suggest         *schema.Message
	Knowledge       []*schema.Document
	Interrupt       *InterruptInfo
}

type SingleAgent struct {
	AgentID   int64
	Name      string
	Desc      string
	IconURI   string
	CreatedAt int64
	UpdatedAt int64
	Version   string
	DeletedAt gorm.DeletedAt

	Variables               []*bot_common.Variable            // 上下文变量
	OnboardingInfo          *bot_common.OnboardingInfo        // 开场白
	ModelInfo               *bot_common.ModelInfo             // 模型信息
	Prompt                  *bot_common.PromptInfo            // 提示词
	Plugin                  []*bot_common.PluginInfo          // 插件（如查询天气）
	Knowledge               *bot_common.Knowledge             // 知识库（Rag）
	Workflow                []*bot_common.WorkflowInfo        // 工作流 (流程)
	SuggestReply            *bot_common.SuggestReplyInfo      // 用户问题建议
	JumpConfig              *bot_common.JumpConfig            // 跳转配置
	BackgroundImageInfoList []*bot_common.BackgroundImageInfo // 聊天背景图
	Database                []*bot_common.Database            // 数据库
	BotMode                 bot_common.BotMode                // 机器人模式
	LayoutInfo              *bot_common.LayoutInfo
	ShortcutCommand         []string
}

type InterruptEventType int64

const (
	InterruptEventType_LocalPlugin         InterruptEventType = 1
	InterruptEventType_Question            InterruptEventType = 2
	InterruptEventType_RequireInfos        InterruptEventType = 3
	InterruptEventType_SceneChat           InterruptEventType = 4
	InterruptEventType_InputNode           InterruptEventType = 5
	InterruptEventType_WorkflowLocalPlugin InterruptEventType = 6
	InterruptEventType_OauthPlugin         InterruptEventType = 7
	InterruptEventType_WorkflowLLM         InterruptEventType = 100
)

// 存档
type InterruptInfo struct {
	AllToolInterruptData map[string]*plugin.ToolInterruptEvent
	ToolCallID           string
	InterruptType        InterruptEventType
	InterruptID          string
}

type ExecuteRequest struct {
	Identity *AgentIdentity
	UserID   string

	Input           *schema.Message
	History         []*schema.Message
	ResumeInfo      *InterruptInfo
	PreCallTools    []*agentrun.ToolsRetriever
	CustomVariables map[string]string
	ConversationID  int64
}

type AgentIdentity struct {
	AgentID int64
	// State   AgentState
	Version string
	IsDraft bool
}
