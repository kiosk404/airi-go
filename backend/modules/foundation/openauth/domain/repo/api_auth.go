package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/foundation/openauth/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/foundation/openauth/infra/repo/gorm_gen/model"
)

//go:generate mockgen -destination=mocks/api_auth.go -package=mocks . ApiAuthRepository
type ApiAuthRepository interface {
	Create(ctx context.Context, do *entity.CreateApiKey) (*entity.ApiKey, error)
	Delete(ctx context.Context, id int64, userID int64) error
	Get(ctx context.Context, id int64) (*model.APIKey, error)
	FindByKey(ctx context.Context, key string) (*model.APIKey, error)
	List(ctx context.Context, userID int64, limit int, page int) ([]*model.APIKey, bool, error)
	Update(ctx context.Context, id int64, userID int64, columnData map[string]any) error
}
