package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/cache"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/infra/dao"
)

func NewSingleAgentRepo(rdb rdb.Provider, idGen idgen.IDGenerator, cli cache.Cmdable) SingleAgentDraftRepo {
	return dao.NewSingleAgentDraftDAO(rdb.NewSession(context.Background()).DB(), idGen, cli)
}

//go:generate mockgen -destination=mocks/single_agent_draft.go -package=mocks . SingleAgentDraftRepo
type SingleAgentDraftRepo interface {
	Create(ctx context.Context, creatorID int64, draft *entity.SingleAgent) (draftID int64, err error)
	CreateWithID(ctx context.Context, agentID, creator int64, draft *entity.SingleAgent) (draftID int64, err error)
	Get(ctx context.Context, agentID int64) (*entity.SingleAgent, error)
	MGet(ctx context.Context, agentIDs []int64) ([]*entity.SingleAgent, error)
	Delete(ctx context.Context, agentID int64) (err error)
	Update(ctx context.Context, agentInfo *entity.SingleAgent) (err error)
	Save(ctx context.Context, agentInfo *entity.SingleAgent) (err error)
	GetDisplayInfo(ctx context.Context, userID, agentID int64) (*entity.AgentDraftDisplayInfo, error)
	UpdateDisplayInfo(ctx context.Context, userID int64, e *entity.AgentDraftDisplayInfo) error
	List(ctx context.Context, page, pageSize int) ([]*entity.SingleAgent, int64, error)
	ListByCreator(ctx context.Context, creatorID int64, page, pageSize int) ([]*entity.SingleAgent, int64, error)
}
