package conf

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
)

//go:generate mockgen -destination=mocks/manage.go -package=mocks . IConfigManage
type IConfigManage interface {
	ListModels(ctx context.Context, req entity.ListModelsRequest) (models []*entity.Model, total int64, hasMore bool, nextPageToken int64, err error)
	GetModel(ctx context.Context, id int64) (model *entity.Model, err error)
	GetModelSet(ctx context.Context) (map[int64]*entity.Model, error)
}
