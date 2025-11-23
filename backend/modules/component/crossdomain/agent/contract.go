package agent

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/agent/model"
	agentrun "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/agentrun/model"
)

// SingleAgent Requests and responses must not reference domain entities and can only use models under api/model/crossdomain.
type SingleAgent interface {
	StreamExecute(ctx context.Context, agentRuntime *AgentRuntime) (*schema.StreamReader[*model.AgentEvent], error)
	ObtainAgentByIdentity(ctx context.Context, identity *model.AgentIdentity) (*model.SingleAgent, error)
	GetSingleAgentDraft(ctx context.Context, agentID int64) (agentInfo *model.SingleAgent, err error)
}

type AgentRuntime struct {
	AgentVersion     string
	UserID           string
	AgentID          int64
	ConversationId   int64
	IsDraft          bool
	PreRetrieveTools []*agentrun.Tool
	CustomVariables  map[string]string

	HistoryMsg []*schema.Message
	Input      *schema.Message
	ResumeInfo *ResumeInfo
}

type ResumeInfo = model.InterruptInfo

type AgentEvent = model.AgentEvent

var defaultSVC SingleAgent

func DefaultSVC() SingleAgent {
	return defaultSVC
}

func SetDefaultSVC(svc SingleAgent) {
	defaultSVC = svc
}

type ShortcutCommandComponentType string

const (
	ShortcutCommandComponentTypeText   ShortcutCommandComponentType = "text"
	ShortcutCommandComponentTypeSelect ShortcutCommandComponentType = "select"
	ShortcutCommandComponentTypeFile   ShortcutCommandComponentType = "file"
)

type ShortcutCommandComponentFileType string

const (
	ShortcutCommandComponentFileTypeImage ShortcutCommandComponentFileType = "image"
	ShortcutCommandComponentFileTypeDoc   ShortcutCommandComponentFileType = "doc"
	ShortcutCommandComponentFileTypeTable ShortcutCommandComponentFileType = "table"
	ShortcutCommandComponentFileTypeAudio ShortcutCommandComponentFileType = "audio"
	ShortcutCommandComponentFileTypeVideo ShortcutCommandComponentFileType = "video"
	ShortcutCommandComponentFileTypeZip   ShortcutCommandComponentFileType = "zip"
	ShortcutCommandComponentFileTypeCode  ShortcutCommandComponentFileType = "code"
	ShortcutCommandComponentFileTypeTxt   ShortcutCommandComponentFileType = "txt"
	ShortcutCommandComponentFileTypePPT   ShortcutCommandComponentFileType = "ppt"
)
