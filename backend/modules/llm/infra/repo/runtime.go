package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	dboption "github.com/kiosk404/airi-go/backend/infra/impl/rdb/common/option"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/repo"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/repo/convertor"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/repo/dao"
)

type RuntimeRepoImpl struct {
	db                rdb.Provider
	modelReqRecordDao dao.IModelRequestRecordDao
}

func NewRuntimeRepo(db rdb.Provider, modelReqRecordDao dao.IModelRequestRecordDao) repo.IRuntimeRepository {
	return &RuntimeRepoImpl{
		db:                db,
		modelReqRecordDao: modelReqRecordDao,
	}
}

func (r *RuntimeRepoImpl) CreateModelRequestRecord(ctx context.Context, record *entity.ModelRequestRecord) (err error) {
	return r.db.Transaction(ctx, func(tx rdb.RDB) error {
		option := dboption.Option{}
		opt := option.WithTransaction(tx)
		return r.modelReqRecordDao.Create(ctx, convertor.ModelReqRecordDO2PO(record), opt)
	})
}
