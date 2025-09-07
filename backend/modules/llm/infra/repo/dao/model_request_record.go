package dao

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/repo/gorm_gen/query"
)

type IModelRequestRecordDao interface {
	Create(ctx context.Context, modelPO *model.ModelRequestRecord, opts ...rdb.Option) (err error)
}

type ModelRequestRecordDaoImpl struct {
	db rdb.Provider
}

func NewModelRequestRecordDao(db rdb.Provider) IModelRequestRecordDao {
	return &ModelRequestRecordDaoImpl{db: db}
}

func (m *ModelRequestRecordDaoImpl) Create(ctx context.Context, modelPO *model.ModelRequestRecord, opts ...rdb.Option) (err error) {
	q := query.Use(m.db.NewSession(ctx, opts...).DB()).WithContext(ctx)
	err = q.ModelRequestRecord.Create(modelPO)
	if err != nil {
		return err
	}
	return nil
}
