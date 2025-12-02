package application

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	"github.com/kiosk404/airi-go/backend/modules/llm/application/convert"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/pkg"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

type BaseChatModel = model.BaseChatModel

type ToolCallingChatModel = model.ToolCallingChatModel

type Service interface {
	Build(ctx context.Context, params *modelmgr.LLMParams) (ToolCallingChatModel, error)
}

var modelClass2NewModelBuilder = map[modelmgr.ModelClass]func(m *modelmgr.Model) Service{
	modelmgr.ModelClass_Ollama:   newOllamaModelBuilder,
	modelmgr.ModelClass_GPT:      newOpenaiModelBuilder,
	modelmgr.ModelClass_DeepSeek: newDeepseekModelBuilder,
	modelmgr.ModelClass_Gemini:   newGeminiModelBuilder,
	modelmgr.ModelClass_QWen:     newQwenModelBuilder,
}

func NewModelBuilder(modelClass modelmgr.ModelClass, cfg *modelmgr.Model) (Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("model config is nil")
	}

	if cfg.Connection == nil {
		return nil, fmt.Errorf("model connection is nil")
	}

	if cfg.Connection.BaseConnInfo == nil {
		return nil, fmt.Errorf("model base connection is nil")
	}

	buildFn, ok := modelClass2NewModelBuilder[modelClass]
	if !ok {
		return nil, fmt.Errorf("model class %v not supported", modelClass)
	}

	return buildFn(cfg), nil
}

func BuildModelByID(ctx context.Context, modelID int64, params *modelmgr.LLMParams) (bcm ToolCallingChatModel, info *modelmgr.Model, err error) {
	m, err := ModelMgrSVC.DomainSVC.GetModelByID(ctx, modelID)
	if err != nil {
		return nil, nil, fmt.Errorf("get model by id failed: %w", err)
	}
	mm := convert.ModelInstance(m)
	bcm, err = buildModelWithConfParams(ctx, mm, params)
	if err != nil {
		return nil, nil, fmt.Errorf("build model failed: %w", err)
	}

	return bcm, mm, nil
}

func BuildModelBySettings(ctx context.Context, appSettings *bot_common.ModelInfo) (bcm ToolCallingChatModel, info *modelmgr.Model, err error) {
	if appSettings == nil {
		return nil, nil, fmt.Errorf("model settings is nil")
	}

	if appSettings.ModelId == nil {
		logs.DebugX(pkg.ModelName, "model id is nil, app settings: %v", conv.DebugJsonToStr(appSettings))
		return nil, nil, fmt.Errorf("model id is nil")
	}

	params := newLLMParamsWithSettings(appSettings)

	return BuildModelByID(ctx, *appSettings.ModelId, params)
}

func buildModelWithConfParams(ctx context.Context, m *modelmgr.Model, params *modelmgr.LLMParams) (bcm ToolCallingChatModel, err error) {
	modelBuilder, err := NewModelBuilder(m.Provider.ModelClass, m)
	if err != nil {
		return nil, fmt.Errorf("new model builder failed: %w", err)
	}

	bcm, err = modelBuilder.Build(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("build model failed: %w", err)
	}

	return bcm, nil
}
