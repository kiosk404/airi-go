package modelbuilder

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	ollamaapi "github.com/ollama/ollama/api"
)

type ollamaModelBuilder struct {
	cfg *modelmgr.Model
}

func newOllamaModelBuilder(cfg *modelmgr.Model) Service {
	return &ollamaModelBuilder{
		cfg: cfg,
	}
}

func (o *ollamaModelBuilder) getDefaultOllamaConfig() *ollama.ChatModelConfig {
	return &ollama.ChatModelConfig{
		Options: &ollamaapi.Options{},
		Model:   o.cfg.Name,
		BaseURL: "http://127.0.0.1:11434",
	}
}

func (o *ollamaModelBuilder) applyParamsToOllamaConfig(conf *ollama.ChatModelConfig, params *modelmgr.LLMParams) {
	if params == nil {
		return
	}

	if params.Temperature != nil {
		conf.Options.Temperature = float32(*params.Temperature)
	}

	if params.TopP != nil {
		conf.Options.TopP = float32(*params.TopP)
	}

	if params.TopK != nil {
		conf.Options.TopK = int(*params.TopK)
	}

	if params.FrequencyPenalty != 0 {
		conf.Options.FrequencyPenalty = float32(params.FrequencyPenalty)
	}

	if params.PresencePenalty != 0 {
		conf.Options.PresencePenalty = float32(params.PresencePenalty)
	}

	conf.Thinking = &ollamaapi.ThinkValue{
		Value: true,
	}

}

func (o *ollamaModelBuilder) Build(ctx context.Context, params *modelmgr.LLMParams) (ToolCallingChatModel, error) {

	conf := o.getDefaultOllamaConfig()

	o.applyParamsToOllamaConfig(conf, params)

	return ollama.NewChatModel(ctx, conf)
}
