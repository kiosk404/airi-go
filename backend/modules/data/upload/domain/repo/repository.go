package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/data/upload/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/data/upload/infra/dao"
	"gorm.io/gorm"
)

func NewFilesRepo(db *gorm.DB) FilesRepo {
	return dao.NewFilesDAO(db)
}

type FilesRepo interface {
	Create(ctx context.Context, file *entity.File) error
	BatchCreate(ctx context.Context, files []*entity.File) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*entity.File, error)
	MGetByIDs(ctx context.Context, ids []int64) ([]*entity.File, error)
}
