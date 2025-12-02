package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
)

type ModelManager interface {
	CreateLLMModel(ctx context.Context, modelClass entity.ModelClass, modelShowName string, conn *entity.Connection, extra *entity.ModelExtra) (int64, error)
	UpdateLLMModel(ctx context.Context, modelID int64, conn *entity.Connection, extra *entity.ModelExtra) error
	GetModelByID(ctx context.Context, id int64) (*entity.ModelInstance, error)
	ListModelByType(ctx context.Context, modelType entity.ModelType, limit int) ([]*entity.ModelInstance, error)
	DeleteModelByID(ctx context.Context, id int64) error
}
