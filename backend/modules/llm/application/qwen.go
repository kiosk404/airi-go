package application

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

type qwenModelBuilder struct {
	cfg *model.Model
}

func newQwenModelBuilder(cfg *model.Model) Service {
	return &qwenModelBuilder{
		cfg: cfg,
	}
}

func (q *qwenModelBuilder) getDefaultQwenConfig() *qwen.ChatModelConfig {
	return &qwen.ChatModelConfig{
		Temperature: ptr.Of(float32(0.7)),
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type:       "text",
			JSONSchema: nil,
		},
	}
}

func (q *qwenModelBuilder) applyParamsToQwenConfig(conf *qwen.ChatModelConfig, params *model.LLMParams) {
	if params == nil {
		return
	}

	conf.TopP = params.TopP

	if params.Temperature != nil {
		conf.Temperature = ptr.Of(*params.Temperature)
	}

	if params.MaxTokens != 0 {
		conf.MaxTokens = ptr.Of(params.MaxTokens)
	}

	if params.FrequencyPenalty != 0 {
		conf.FrequencyPenalty = ptr.Of(params.FrequencyPenalty)
	}

	if params.PresencePenalty != 0 {
		conf.PresencePenalty = ptr.Of(params.PresencePenalty)
	}

	if params.EnableThinking != nil {
		conf.EnableThinking = params.EnableThinking
	}
}

func (q *qwenModelBuilder) Build(ctx context.Context, params *model.LLMParams) (ToolCallingChatModel, error) {
	base := q.cfg.Connection.BaseConnInfo

	conf := q.getDefaultQwenConfig()
	conf.APIKey = base.APIKey
	conf.BaseURL = base.BaseURL
	conf.Model = base.Model

	switch base.ThinkingType {
	case model.ThinkingType_Enable:
		conf.EnableThinking = ptr.Of(true)
	case model.ThinkingType_Disable:
		conf.EnableThinking = ptr.Of(false)
	}

	q.applyParamsToQwenConfig(conf, params)

	return qwen.NewChatModel(ctx, conf)
}
