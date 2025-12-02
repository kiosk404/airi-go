package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/dao"
	"gorm.io/gorm"
)

func NewModelRunRecordRepo(db *gorm.DB) ModelRunRecordRepository {
	return dao.NewModelRunRecordDao(db)
}

type ModelRunRecordRepository interface {
	Create(ctx context.Context, runRecord *entity.ModelRequestRecord) (err error)
	List(ctx context.Context, modelID string, limit int) (modelRunRecordList []*entity.ModelRequestRecord)
}
