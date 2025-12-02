package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/dao"
)

func NewModelMgrRepo(db rdb.Provider, idGen idgen.IDGenerator) ModelMgrRepository {
	return dao.NewModelMgrDao(db.NewSession(context.Background()).DB(), idGen)
}

type ModelMgrRepository interface {
	CreateModel(ctx context.Context, instance *entity.ModelInstance) (id int64, err error)
	GetModel(ctx context.Context, id int64) (do *entity.ModelInstance, err error)
	DeleteModel(ctx context.Context, id int64) (err error)
	ListModels(ctx context.Context) (do []*entity.ModelInstance, err error)
	ListModelByType(ctx context.Context, modelClass entity.ModelType, limit int) (do []*entity.ModelInstance, err error)
	UpdateModel(ctx context.Context, instance *entity.ModelInstance) (err error)
}
