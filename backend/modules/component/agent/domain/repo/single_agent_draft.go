package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
)

//go:generate mockgen -destination=mocks/single_agent_draft.go -package=mocks . SingleAgentDraftRepo
type SingleAgentDraftRepo interface {
	Create(ctx context.Context, draft *entity.SingleAgent) (draftID int64, err error)
	CreateWithID(ctx context.Context, agentID int64, draft *entity.SingleAgent) (draftID int64, err error)
	Get(ctx context.Context, agentID int64) (*entity.SingleAgent, error)
	MGet(ctx context.Context, agentIDs []int64) ([]*entity.SingleAgent, error)
	Delete(ctx context.Context, agentID int64) (err error)
	Update(ctx context.Context, agentInfo *entity.SingleAgent) (err error)
	Save(ctx context.Context, agentInfo *entity.SingleAgent) (err error)
	GetDisplayInfo(ctx context.Context, userID, agentID int64) (*entity.AgentDraftDisplayInfo, error)
	UpdateDisplayInfo(ctx context.Context, userID int64, e *entity.AgentDraftDisplayInfo) error
}
