package modelbuilder

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	crossmodelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

type BaseChatModel = model.BaseChatModel

type ToolCallingChatModel = model.ToolCallingChatModel

type Service interface {
	Build(ctx context.Context, params *modelmgr.LLMParams) (ToolCallingChatModel, error)
}

var modelClass2NewModelBuilder = map[modelmgr.Protocol]func(*modelmgr.Model) Service{
	modelmgr.ProtocolOllama: newOllamaModelBuilder,
}

func NewModelBuilder(modelProtocol modelmgr.Protocol, cfg *modelmgr.Model) (Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("model config is nil")
	}

	buildFn, ok := modelClass2NewModelBuilder[modelProtocol]
	if !ok {
		return nil, fmt.Errorf("model class %v not supported", modelProtocol)
	}
	return buildFn(cfg), nil
}

func BuildModelBySettings(ctx context.Context, appSettings *bot_common.ModelInfo) (bcm ToolCallingChatModel, info *Model, err error) {
	if appSettings == nil {
		return nil, nil, fmt.Errorf("model settings is nil")
	}

	if appSettings.ModelId == nil {
		logs.Debug("model id is nil, app settings: %v", conv.DebugJsonToStr(appSettings))
		return nil, nil, fmt.Errorf("model id is nil")
	}

	params := newLLMParamsWithSettings(appSettings)

	return BuildModelByID(ctx, *appSettings.ModelId, params)
}

func BuildModelByID(ctx context.Context, modelID int64, params *modelmgr.LLMParams) (bcm ToolCallingChatModel, info *Model, err error) {
	req := &modelmgr.GetModelRequest{ModelID: modelID}
	m, err := crossmodelmgr.DefaultSVC().GetModel(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("get model by id failed: %w", err)
	}

	bcm, err = buildModelWithConfParams(ctx, m.Model, params)
	if err != nil {
		return nil, nil, fmt.Errorf("build model failed: %w", err)
	}

	return bcm, &Model{Model: m.Model}, err
}

func buildModelWithConfParams(ctx context.Context, m *modelmgr.Model, params *modelmgr.LLMParams) (bcm ToolCallingChatModel, err error) {
	modelBuilder, err := NewModelBuilder(m.GetProtocol(), m)
	if err != nil {
		return nil, fmt.Errorf("new model builder failed: %w", err)
	}

	bcm, err = modelBuilder.Build(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("build model failed: %w", err)
	}

	return bcm, nil
}
