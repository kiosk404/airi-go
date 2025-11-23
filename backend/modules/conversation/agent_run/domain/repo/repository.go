package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/infra/repo/dao"
)

func NewRunRecordRepo(rdb rdb.Provider, idGen idgen.IDGenerator) RunRecordRepo {
	return dao.NewRunRecordDAO(rdb.NewSession(context.Background()).DB(), idGen)
}

type RunRecordRepo interface {
	Create(ctx context.Context, runMeta *entity.AgentRunMeta) (*entity.RunRecordMeta, error)
	GetByID(ctx context.Context, id int64) (*entity.RunRecordMeta, error)
	Cancel(ctx context.Context, req *entity.CancelRunMeta) (*entity.RunRecordMeta, error)
	Delete(ctx context.Context, id []int64) error
	UpdateByID(ctx context.Context, id int64, update *entity.UpdateMeta) error
	List(ctx context.Context, meta *entity.ListRunRecordMeta) ([]*entity.RunRecordMeta, error)
}
