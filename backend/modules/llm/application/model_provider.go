package application

import (
	"github.com/kiosk404/airi-go/backend/api/model/app/developer_api"
	"github.com/kiosk404/airi-go/backend/api/model/modelapi"
)

func getModelProviderList() []*modelapi.ModelProvider {
	return []*modelapi.ModelProvider{
		{
			Name: &modelapi.I18nText{
				ZhCn: "Claude 模型",
				EnUs: "Claude Model",
			},
			IconURI: "default_icon/claude_v2.png",
			Description: &modelapi.I18nText{
				ZhCn: "Claude 模型家族",
				EnUs: "claude model family",
			},
			ModelClass: developer_api.ModelClass_Claude,
		},
		{
			Name: &modelapi.I18nText{
				ZhCn: "Deepseek 模型",
				EnUs: "Deepseek Model",
			},
			IconURI: "default_icon/deepseek_v2.png",
			Description: &modelapi.I18nText{
				ZhCn: "Deepseek 模型家族",
				EnUs: "deepseek model family",
			},
			ModelClass: developer_api.ModelClass_DeepSeek,
		},
		{
			Name: &modelapi.I18nText{
				ZhCn: "Gemini 模型",
				EnUs: "Gemini Model",
			},
			IconURI: "default_icon/gemini_v2.png",
			Description: &modelapi.I18nText{
				ZhCn: "Gemini 模型家族",
				EnUs: "gemini model family",
			},
			ModelClass: developer_api.ModelClass_Gemini,
		},
		{
			Name: &modelapi.I18nText{
				ZhCn: "Ollama 模型",
				EnUs: "Ollama Model",
			},
			IconURI: "default_icon/ollama.png",
			Description: &modelapi.I18nText{
				ZhCn: "Ollama 模型家族",
				EnUs: "ollama model family",
			},
			ModelClass: developer_api.ModelClass_Ollama,
		},
		{
			Name: &modelapi.I18nText{
				ZhCn: "OpenAI 模型",
				EnUs: "OpenAI Model",
			},
			IconURI: "default_icon/openai_v2.png",
			Description: &modelapi.I18nText{
				ZhCn: "OpenAI 模型家族",
				EnUs: "openai model family",
			},
			ModelClass: developer_api.ModelClass_GPT,
		},
		{
			Name: &modelapi.I18nText{
				ZhCn: "Qwen 模型",
				EnUs: "Qwen Model",
			},
			IconURI: "default_icon/qwen_v2.png",
			Description: &modelapi.I18nText{
				ZhCn: "Qwen 模型家族",
				EnUs: "qwen model family",
			},
			ModelClass: developer_api.ModelClass_QWen,
		},
	}
}

func SupportProtocol(class developer_api.ModelClass) bool {
	_, ok := GetModelProvider(class)

	return ok
}

func GetModelProvider(class developer_api.ModelClass) (*modelapi.ModelProvider, bool) {
	modelProviders := getModelProviderList()
	for _, modelProvider := range modelProviders {
		if modelProvider.ModelClass == class {
			return modelProvider, true
		}
	}

	return nil, false
}
