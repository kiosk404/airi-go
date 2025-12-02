package service

import (
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

func getModelProviderList() []*model.ModelProvider {
	return []*model.ModelProvider{
		{
			Name:        &model.I18nText{ZhCn: "OpenAI", EnUs: "OpenAI"},
			IconURI:     "openai-icon",
			IconURL:     "https://openai.com/favicon.ico",
			Description: &model.I18nText{ZhCn: "OpenAI models", EnUs: "OpenAI models"},
			ModelClass:  model.ModelClass_GPT,
		},
		{
			Name:        &model.I18nText{ZhCn: "Ollama", EnUs: "Ollama"},
			IconURI:     "Ollama",
			IconURL:     "https://ollama.com/favicon.ico",
			Description: &model.I18nText{ZhCn: "Ollama models", EnUs: "Ollama models"},
			ModelClass:  model.ModelClass_Ollama,
		},
	}
}

func GetModelProvider(class entity.ModelClass) (*model.ModelProvider, bool) {
	modelProviders := getModelProviderList()
	for _, modelProvider := range modelProviders {
		if modelProvider.ModelClass == ptr.From(class.ModelClass) {
			return modelProvider, true
		}
	}

	return nil, false
}
