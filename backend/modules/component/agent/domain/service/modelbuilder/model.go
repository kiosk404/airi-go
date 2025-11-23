package modelbuilder

import (
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
)

type Model struct {
	*model.Model
}

const (
	temperature      = "temperature"
	maxTokens        = "max_tokens"
	topP             = "top_p"
	topK             = "top_k"
	responseFormat   = "response_format"
	frequencyPenalty = "frequency_penalty"
	presencePenalty  = "presence_penalty"
)

func (m *Model) GetDefaultTemperature() *float32 {
	cfg := m.ParamConfig.GetCommonParamDefaultVal()
	return cfg.Temperature
}

func (m *Model) GetDefaultMaxTokens() *int {
	cfg := m.ParamConfig.GetCommonParamDefaultVal()
	return cfg.MaxTokens
}

func (m *Model) GetDefaultTopP() *float32 {
	cfg := m.ParamConfig.GetCommonParamDefaultVal()
	return cfg.TopP
}

func (m *Model) GetDefaultFrequencyPenalty() *float32 {
	cfg := m.ParamConfig.GetCommonParamDefaultVal()
	return cfg.FrequencyPenalty
}

func (m *Model) GetDefaultPresencePenalty() *float32 {
	cfg := m.ParamConfig.GetCommonParamDefaultVal()
	return cfg.PresencePenalty
}

func (m *Model) GetDefaultTopK() *int {
	cfg := m.ParamConfig.GetCommonParamDefaultVal()
	return cfg.TopK
}
