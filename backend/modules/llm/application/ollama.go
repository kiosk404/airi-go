package application

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/ollama/ollama/api"
)

type ollamaModelBuilder struct {
	cfg *model.Model
}

func newOllamaModelBuilder(cfg *model.Model) Service {
	return &ollamaModelBuilder{
		cfg: cfg,
	}
}

func (o *ollamaModelBuilder) getDefaultOllamaConfig() *ollama.ChatModelConfig {
	return &ollama.ChatModelConfig{
		Options: &api.Options{},
		BaseURL: "http://127.0.0.1:11434",
	}
}

func (o *ollamaModelBuilder) applyParamsToOllamaConfig(conf *ollama.ChatModelConfig, params *model.LLMParams) {
	if params == nil {
		return
	}

	if params.Temperature != nil {
		conf.Options.Temperature = *params.Temperature
	}

	if params.TopP != nil {
		conf.Options.TopP = *params.TopP
	}

	if params.TopK != nil {
		conf.Options.TopK = int(*params.TopK)
	}

	if params.FrequencyPenalty != 0 {
		conf.Options.FrequencyPenalty = params.FrequencyPenalty
	}

	if params.PresencePenalty != 0 {
		conf.Options.PresencePenalty = params.PresencePenalty
	}

	if params.EnableThinking != nil {
		conf.Thinking = &api.ThinkValue{
			Value: ptr.From(params.EnableThinking),
		}
	}
}

func (o *ollamaModelBuilder) Build(ctx context.Context, params *model.LLMParams) (ToolCallingChatModel, error) {
	base := o.cfg.Connection.BaseConnInfo

	conf := o.getDefaultOllamaConfig()
	if base.BaseURL != "" {
		conf.BaseURL = base.BaseURL
	}
	conf.Model = base.Model

	switch base.ThinkingType {
	case model.ThinkingType_Enable:
		conf.Thinking = &api.ThinkValue{
			Value: ptr.Of(true),
		}
	case model.ThinkingType_Disable:
		conf.Thinking = &api.ThinkValue{
			Value: ptr.Of(false),
		}
	}

	o.applyParamsToOllamaConfig(conf, params)

	return ollama.NewChatModel(ctx, conf)
}
