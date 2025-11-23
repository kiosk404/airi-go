package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/modules/component/prompt/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/prompt/infra/dao"
	"gorm.io/gorm"
)

func NewPromptRepo(db *gorm.DB, generator idgen.IDGenerator) PromptRepository {
	return dao.NewPromptDAO(db, generator)
}

type PromptRepository interface {
	CreatePromptResource(ctx context.Context, do *entity.PromptResource) (int64, error)
	GetPromptResource(ctx context.Context, promptID int64) (*entity.PromptResource, error)
	UpdatePromptResource(ctx context.Context, promptID int64, name, description, promptText *string) error
	DeletePromptResource(ctx context.Context, ID int64) error
}
