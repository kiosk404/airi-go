package service

import (
	"context"
	"fmt"

	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/repo"
	"github.com/kiosk404/airi-go/backend/pkg/conf"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
)

type modelManageImpl struct {
	oss             storage.Storage
	ModelManageRepo repo.ModelMgrRepository
	ModelMeta       *ModelMetaConf
}

func NewService(oss storage.Storage, modelManageRepo repo.ModelMgrRepository, configFactory conf.IConfigLoaderFactory) ModelManager {
	modelMeta, err := initModelMetaConf(context.Background(), configFactory)
	if err != nil {
		panic(fmt.Sprintf("init model meta conf failed, err=%v", err))
	}
	return &modelManageImpl{
		oss:             oss,
		ModelManageRepo: modelManageRepo,
		ModelMeta:       modelMeta,
	}
}

func (m modelManageImpl) CreateLLMModel(ctx context.Context, modelClass entity.ModelClass, modelShowName string, conn *entity.Connection, extra *entity.ModelExtra) (int64, error) {
	if conn == nil {
		return 0, fmt.Errorf("connection is nil")
	}

	if conn.BaseConnInfo == nil {
		return 0, fmt.Errorf("base conn info is nil")
	}

	provider, ok := GetModelProvider(modelClass)
	if !ok {
		return 0, fmt.Errorf("model class %s not supported", modelClass)
	}

	conn, err := encryptConn(ctx, conn)
	if err != nil {
		return 0, err
	}
	modelName := conn.BaseConnInfo.Model
	modelMeta, err := m.ModelMeta.GetModelMeta(modelClass.Model(), modelName)
	if err != nil {
		return 0, fmt.Errorf("get model meta failed, err: %w", err)
	}
	if modelMeta.Connection != nil {
		conn.Openai = modelMeta.Connection.Openai
		conn.Deepseek = modelMeta.Connection.Deepseek
		conn.Gemini = modelMeta.Connection.Gemini
		conn.Qwen = modelMeta.Connection.Qwen
		conn.Ollama = modelMeta.Connection.Ollama
		conn.Claude = modelMeta.Connection.Claude
	}

	mInstance := &entity.ModelInstance{
		Type:       entity.ModelType{ModelType: ptr.Of(model.ModelType_LLM)},
		Provider:   ptr.From(provider),
		Connection: ptr.From(conn),
		Capability: ptr.From(modelMeta.Capability),
		Parameters: slices.Transform(modelMeta.Parameters, func(p *model.ModelParameter) model.ModelParameter {
			return ptr.From(p)
		}),
		DisplayInfo: ptr.From(modelMeta.DisplayInfo),
		Extra:       ptr.From(extra),
	}

	id, err := m.ModelManageRepo.CreateModel(ctx, mInstance)
	if err != nil {
		return 0, fmt.Errorf("create model failed, err: %w", err)
	}
	return id, nil
}

func (m modelManageImpl) UpdateLLMModel(ctx context.Context, modelID int64, conn *entity.Connection, extra *entity.ModelExtra) error {
	conn, err := encryptConn(ctx, conn)
	if err != nil {
		return err
	}
	instance, err := m.GetModelByID(ctx, modelID)
	if err != nil {
		return err
	}
	instance.Connection = ptr.From(conn)
	instance.Extra = ptr.From(extra)
	return m.ModelManageRepo.UpdateModel(ctx, instance)
}

func encryptConn(ctx context.Context, conn *entity.Connection) (*entity.Connection, error) {
	// encrypt conn if you need
	return conn, nil
}

func (m modelManageImpl) GetModelByID(ctx context.Context, id int64) (*entity.ModelInstance, error) {
	return m.ModelManageRepo.GetModel(ctx, id)
}

func (m modelManageImpl) ListModelByType(ctx context.Context, modelType entity.ModelType, limit int) ([]*entity.ModelInstance, error) {
	return m.ModelManageRepo.ListModelByType(ctx, modelType, limit)
}

func (m modelManageImpl) DeleteModelByID(ctx context.Context, id int64) error {
	return m.ModelManageRepo.DeleteModel(ctx, id)
}
