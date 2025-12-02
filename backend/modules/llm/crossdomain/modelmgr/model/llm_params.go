package model

import (
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

type LLMParams struct {
	Temperature      *float32            `json:"temperature"`
	FrequencyPenalty float32             `json:"frequencyPenalty"`
	PresencePenalty  float32             `json:"presencePenalty"`
	MaxTokens        int                 `json:"maxTokens"`
	TopP             *float32            `json:"topP"`
	TopK             *int32              `json:"topK"`
	ResponseFormat   ModelResponseFormat `json:"responseFormat"`
	EnableThinking   *bool               `json:"enable_thinking,omitempty" yaml:"enable_thinking,omitempty"`
}

type ModelResponseFormat int64

const (
	ModelResponseFormat_Text     ModelResponseFormat = 0
	ModelResponseFormat_Markdown ModelResponseFormat = 1
	ModelResponseFormat_JSON     ModelResponseFormat = 2
)

func (p ModelResponseFormat) String() string {
	switch p {
	case ModelResponseFormat_Text:
		return "Text"
	case ModelResponseFormat_Markdown:
		return "Markdown"
	case ModelResponseFormat_JSON:
		return "JSON"
	}
	return "<UNSET>"
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

func (m *Model) GetDefaultTemperature() *float64 {
	for _, param := range m.Parameters {
		if param.Name == temperature && param.DefaultVal != nil {
			t, err := conv.StrToFloat64(param.DefaultVal.DefaultVal)
			if err != nil {
				return nil
			}

			return ptr.Of(t)
		}
	}

	return nil
}

func (m *Model) GetDefaultMaxTokens() *int32 {
	for _, param := range m.Parameters {
		if param.Name == maxTokens && param.DefaultVal != nil {
			t, err := conv.StrToInt64(param.DefaultVal.DefaultVal)
			if err != nil {
				return nil
			}

			return ptr.Of(int32(t))
		}
	}

	return nil
}

func (m *Model) GetDefaultTopP() *float64 {
	for _, param := range m.Parameters {
		if param.Name == topP && param.DefaultVal != nil {
			t, err := conv.StrToFloat64(param.DefaultVal.DefaultVal)
			if err != nil {
				return nil
			}

			return ptr.Of(t)
		}
	}

	return nil
}

func (m *Model) GetDefaultFrequencyPenalty() *float64 {
	for _, param := range m.Parameters {
		if param.Name == frequencyPenalty && param.DefaultVal != nil {
			t, err := conv.StrToFloat64(param.DefaultVal.DefaultVal)
			if err != nil {
				return nil
			}

			return ptr.Of(t)
		}
	}

	return nil
}

func (m *Model) GetDefaultPresencePenalty() *float64 {
	for _, param := range m.Parameters {
		if param.Name == presencePenalty && param.DefaultVal != nil {
			t, err := conv.StrToFloat64(param.DefaultVal.DefaultVal)
			if err != nil {
				return nil
			}

			return ptr.Of(t)
		}
	}

	return nil
}

func (m *Model) GetDefaultTopK() *int32 {
	for _, param := range m.Parameters {
		if param.Name == topK && param.DefaultVal != nil {
			t, err := conv.StrToInt64(param.DefaultVal.DefaultVal)
			if err != nil {
				return nil
			}

			return ptr.Of(int32(t))
		}
	}

	return nil
}
