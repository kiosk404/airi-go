package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/infra/dao"
)

func NewSingleAgentVersionRepo(rdb rdb.Provider, idGen idgen.IDGenerator) SingleAgentVersionRepo {
	return dao.NewSingleAgentVersion(rdb.NewSession(context.Background()).DB(), idGen)
}

//go:generate mockgen -destination=mocks/single_agent_version.go -package=mocks . SingleAgentVersionRepo
type SingleAgentVersionRepo interface {
	GetLatest(ctx context.Context, agentID int64) (*entity.SingleAgent, error)
	Get(ctx context.Context, agentID int64, version string) (*entity.SingleAgent, error)
	List(ctx context.Context, agentID int64, pageIndex, pageSize int32) ([]*entity.SingleAgentPublish, error)
	SavePublishRecord(ctx context.Context, p *entity.SingleAgentPublish, e *entity.SingleAgent) (err error)
	Create(ctx context.Context, version string, e *entity.SingleAgent) (int64, error)
}
