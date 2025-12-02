package application

import (
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

func newLLMParamsWithSettings(appSettings *bot_common.ModelInfo) *model.LLMParams {
	if appSettings == nil {
		return nil
	}

	l := &model.LLMParams{}

	if appSettings.Temperature != nil {
		t := float32(*appSettings.Temperature)
		l.Temperature = &t
	}
	if appSettings.FrequencyPenalty != nil {
		f := float32(*appSettings.FrequencyPenalty)
		l.FrequencyPenalty = f
	}
	if appSettings.PresencePenalty != nil {
		p := float32(*appSettings.PresencePenalty)
		l.PresencePenalty = p
	}
	if appSettings.MaxTokens != nil {
		l.MaxTokens = int(*appSettings.MaxTokens)
	}
	if appSettings.TopP != nil {
		t := float32(*appSettings.TopP)
		l.TopP = &t
	}
	if appSettings.TopK != nil {
		k := int32(*appSettings.TopK)
		l.TopK = &k
	}
	if appSettings.ResponseFormat != nil {
		l.ResponseFormat = model.ModelResponseFormat(ptr.From(appSettings.ResponseFormat))
	}

	return l
}
