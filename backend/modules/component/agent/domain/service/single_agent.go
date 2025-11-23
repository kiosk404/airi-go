package service

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
)

type SingleAgent interface {
	// draft agent
	CreateSingleAgentDraft(ctx context.Context, creatorID int64, draft *entity.SingleAgent) (agentID int64, err error)
	CreateSingleAgentDraftWithID(ctx context.Context, agentID, creator int64, draft *entity.SingleAgent) (int64, error)
	MGetSingleAgentDraft(ctx context.Context, agentIDs []int64) (agents []*entity.SingleAgent, err error)
	GetSingleAgentDraft(ctx context.Context, agentID int64) (agentInfo *entity.SingleAgent, err error)
	UpdateSingleAgentDraft(ctx context.Context, agentInfo *entity.SingleAgent) (err error)
	DeleteAgentDraft(ctx context.Context, agentID int64) (err error)
	UpdateAgentDraftDisplayInfo(ctx context.Context, userID int64, e *entity.AgentDraftDisplayInfo) error
	GetAgentDraftDisplayInfo(ctx context.Context, userID, agentID int64) (*entity.AgentDraftDisplayInfo, error)

	// online agent
	CreateSingleAgent(ctx context.Context, version string, e *entity.SingleAgent) (int64, error)
	DuplicateInMemory(ctx context.Context, req *entity.DuplicateInfo) (newAgent *entity.SingleAgent, err error)
	StreamExecute(ctx context.Context, req *entity.ExecuteRequest) (events *schema.StreamReader[*entity.AgentEvent], err error)
	GetSingleAgent(ctx context.Context, agentID int64, version string) (botInfo *entity.SingleAgent, err error)
	ListAgentPublishHistory(ctx context.Context, agentID int64, pageIndex, pageSize int32) ([]*entity.SingleAgentPublish, error)
	// ObtainAgentByIdentity support obtain agent by agentID
	ObtainAgentByIdentity(ctx context.Context, identity *entity.AgentIdentity) (*entity.SingleAgent, error)

	// Publish
	GetPublishedTime(ctx context.Context, agentID int64) (int64, error)
	GetPublishedInfo(ctx context.Context, agentID int64) (*entity.PublishInfo, error)
	SavePublishRecord(ctx context.Context, p *entity.SingleAgentPublish, e *entity.SingleAgent) error
}
