package modelmgr

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
)

type ModelManager interface {
	CreateModel(ctx context.Context, request model.CreateModelRequest) (int64, error)
	UpdateModel(ctx context.Context, request model.UpdateModelRequest) error
	DeleteModel(ctx context.Context, modelID int64) error
	GetModelByID(ctx context.Context, modelID int64) (*model.Model, error)
	MGetModelByID(ctx context.Context, ids []int64) ([]*model.Model, error)
	GetOnlineModelListWithLimit(ctx context.Context, limit int) ([]*model.Model, error)
	GetOnlineModelList(ctx context.Context) ([]*model.Model, error)
	GetAllModelList(ctx context.Context) ([]*model.Model, error)
}

type ModelManagerApp = model.Model

var defaultSVC ModelManager

func DefaultSVC() ModelManager {
	return defaultSVC
}

func SetDefaultSVC(svc ModelManager) {
	defaultSVC = svc
}
