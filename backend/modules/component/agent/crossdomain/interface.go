package crossdomain

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/api/crossdomain/agentrun"
	"github.com/kiosk404/airi-go/backend/api/crossdomain/singleagent"
)

type SingleAgent interface {
	StreamExecute(ctx context.Context,
		agentRuntime *AgentRuntime) (*schema.StreamReader[*singleagent.AgentEvent], error)
	ObtainAgentByIdentity(ctx context.Context, identity *singleagent.AgentIdentity) (*singleagent.SingleAgent, error)
}

type AgentRuntime struct {
	AgentVersion     string
	UserID           string
	AgentID          int64
	IsDraft          bool
	SpaceID          int64
	ConnectorID      int64
	PreRetrieveTools []*agentrun.Tool

	HistoryMsg []*schema.Message
	Input      *schema.Message
	ResumeInfo *ResumeInfo
}

type ResumeInfo = singleagent.InterruptInfo

type AgentEvent = singleagent.AgentEvent

var defaultSVC SingleAgent

func DefaultSVC() SingleAgent {
	return defaultSVC
}

func SetDefaultSVC(svc SingleAgent) {
	defaultSVC = svc
}
