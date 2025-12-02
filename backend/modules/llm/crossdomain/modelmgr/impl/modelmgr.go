package impl

import (
	"context"

	crossmodelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
)

var defaultSVC crossmodelmgr.ModelManager

type impl struct {
	DomainSVC modelmgr.ModelManager
}

func InitDomainService(m modelmgr.ModelManager) crossmodelmgr.ModelManager {
	defaultSVC = &impl{
		DomainSVC: m,
	}
	return defaultSVC
}

func (i impl) CreateModel(ctx context.Context, request model.CreateModelRequest) (int64, error) {
	modelClassEntity := entity.ModelClass{ModelClass: ptr.Of(request.ModelClass)}
	modelConnEntity := &entity.Connection{Connection: ptr.Of(request.Conn)}
	extraEntity := &entity.ModelExtra{ModelExtra: ptr.Of(request.Extra)}

	modelID, err := i.DomainSVC.CreateLLMModel(ctx, modelClassEntity, request.ModelShowName, modelConnEntity, extraEntity)
	if err != nil {
		return 0, err
	}
	return modelID, nil
}

func (i impl) UpdateModel(ctx context.Context, request model.UpdateModelRequest) error {
	modelConnEntity := &entity.Connection{Connection: ptr.Of(request.Conn)}
	extraEntity := &entity.ModelExtra{ModelExtra: ptr.Of(request.Extra)}

	return i.DomainSVC.UpdateLLMModel(ctx, request.ID, modelConnEntity, extraEntity)
}

func (i impl) DeleteModel(ctx context.Context, modelID int64) error {
	return i.DomainSVC.DeleteModelByID(ctx, modelID)
}

func (i impl) GetModelByID(ctx context.Context, modelID int64) (*model.Model, error) {
	instance, err := i.DomainSVC.GetModelByID(ctx, modelID)
	if err != nil {
		return nil, err
	}
	return &model.Model{
		ID:              instance.ID,
		Provider:        ptr.Of(instance.Provider),
		DisplayInfo:     ptr.Of(instance.DisplayInfo),
		Capability:      ptr.Of(instance.Capability),
		Connection:      ptr.Of(instance.Connection.Model()),
		Type:            instance.Type.Model(),
		Parameters:      instance.Parameters,
		EnableBase64URL: instance.Extra.EnableBase64URL,
	}, nil
}

func (i impl) MGetModelByID(ctx context.Context, ids []int64) ([]*model.Model, error) {
	models := make([]*model.Model, 0, len(ids))
	for _, modelID := range ids {
		instance, err := i.DomainSVC.GetModelByID(ctx, modelID)
		if err != nil {
			return nil, err
		}
		models = append(models, &model.Model{
			ID:              instance.ID,
			Provider:        ptr.Of(instance.Provider),
			DisplayInfo:     ptr.Of(instance.DisplayInfo),
			Capability:      ptr.Of(instance.Capability),
			Connection:      ptr.Of(instance.Connection.Model()),
			Type:            instance.Type.Model(),
			Parameters:      instance.Parameters,
			EnableBase64URL: instance.Extra.EnableBase64URL,
		})
	}
	return models, nil
}

func (i impl) GetOnlineModelListWithLimit(ctx context.Context, limit int) ([]*model.Model, error) {
	modelList, err := i.DomainSVC.ListModelByType(ctx, entity.LLMModelType(), limit)
	if err != nil {
		return nil, err
	}
	models := make([]*model.Model, 0, len(modelList))
	slices.ForEach(modelList, func(instance *entity.ModelInstance, i int) {
		models = append(models, &model.Model{
			ID:              instance.ID,
			Provider:        ptr.Of(instance.Provider),
			DisplayInfo:     ptr.Of(instance.DisplayInfo),
			Capability:      ptr.Of(instance.Capability),
			Connection:      ptr.Of(instance.Connection.Model()),
			Type:            instance.Type.Model(),
			Parameters:      instance.Parameters,
			EnableBase64URL: instance.Extra.EnableBase64URL,
		})
	})

	return models, nil
}

func (i impl) GetOnlineModelList(ctx context.Context) ([]*model.Model, error) {
	//TODO implement me
	panic("implement me")
}

func (i impl) GetAllModelList(ctx context.Context) ([]*model.Model, error) {
	//TODO implement me
	panic("implement me")
}
