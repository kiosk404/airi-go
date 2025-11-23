package modelbuilder

import (
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

type LLMParams struct {
	Temperature      *float32                       `json:"temperature"`
	FrequencyPenalty float32                        `json:"frequencyPenalty"`
	PresencePenalty  float32                        `json:"presencePenalty"`
	MaxTokens        int                            `json:"maxTokens"`
	TopP             *float32                       `json:"topP"`
	TopK             *int32                         `json:"topK"`
	ResponseFormat   bot_common.ModelResponseFormat `json:"responseFormat"`
	EnableThinking   *bool                          `json:"enable_thinking,omitempty" yaml:"enable_thinking,omitempty"`
}

func newLLMParamsWithSettings(appSettings *bot_common.ModelInfo) *modelmgr.LLMParams {
	if appSettings == nil {
		return nil
	}

	l := &modelmgr.LLMParams{}

	if appSettings.Temperature != nil {
		l.Temperature = appSettings.Temperature
	}
	if appSettings.FrequencyPenalty != nil {
		l.FrequencyPenalty = ptr.From(appSettings.FrequencyPenalty)
	}
	if appSettings.PresencePenalty != nil {
		l.PresencePenalty = ptr.From(appSettings.PresencePenalty)
	}
	if appSettings.MaxTokens != nil {
		l.MaxTokens = int(*appSettings.MaxTokens)
	}
	if appSettings.TopP != nil {
		l.TopP = appSettings.TopP
	}
	if appSettings.TopK != nil {
		k := int(*appSettings.TopK)
		l.TopK = &k
	}

	if appSettings.ResponseFormat != nil {
		switch *appSettings.ResponseFormat {
		case bot_common.ModelResponseFormat_Text:
			l.ResponseFormat = modelmgr.ResponseFormatText
		case bot_common.ModelResponseFormat_Markdown:
			l.ResponseFormat = modelmgr.ResponseFormatMarkdown
		case bot_common.ModelResponseFormat_JSON:
			l.ResponseFormat = modelmgr.ResponseFormatJSON
		}
	}

	return l
}
