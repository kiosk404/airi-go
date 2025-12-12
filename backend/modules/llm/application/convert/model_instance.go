package convert

import (
	"context"

	"github.com/kiosk404/airi-go/backend/api/model/app/developer_api"
	"github.com/kiosk404/airi-go/backend/api/model/modelapi"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/llm/pkg"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ternary"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

func ModelInstance(model *entity.ModelInstance) *modelmgr.Model {
	return &modelmgr.Model{
		ID:              model.ID,
		Provider:        ptr.Of(model.Provider),
		DisplayInfo:     ptr.Of(model.DisplayInfo),
		Capability:      ptr.Of(model.Capability),
		Connection:      ptr.Of(model.Connection.Model()),
		Type:            model.Type.Model(),
		Parameters:      model.Parameters,
		EnableBase64URL: model.Extra.EnableBase64URL,
	}
}

func ModelClassDao(modelClass modelmgr.ModelClass) entity.ModelClass {
	return entity.ModelClass{ModelClass: ptr.Of(modelClass)}
}

func ConnectionDao(connection *modelmgr.Connection) *entity.Connection {
	return &entity.Connection{Connection: connection}
}

func ToModel(ctx context.Context, oss storage.Storage, model *entity.ModelInstance) *modelapi.Model {
	if model.Provider.IconURL != "" {
		url, err := oss.GetObjectUrl(ctx, model.Provider.IconURL)
		if err != nil {
			logs.WarnX(pkg.ModelName, "get model icon url failed, err: %v", err)
		} else {
			model.Provider.IconURL = url
		}
	}

	conn, err := decryptConn(ctx, model.Connection)
	if err != nil {
		logs.WarnX(pkg.ModelName, "decrypt model connection failed, err: %v", err)
	}

	m := &modelapi.Model{
		ID:              model.ID,
		Provider:        convertModelProvider(model.Provider),
		DisplayInfo:     convertModelDisplayInfo(model.DisplayInfo),
		Capability:      convertModelCapability(model.Capability),
		Connection:      convertModelConnection(conn),
		Type:            modelapi.ModelType(model.Type.Model()),
		Parameters:      convertModelParameters(model.Parameters),
		EnableBase64URL: model.Extra.EnableBase64URL,
	}

	m.Status = ternary.IFElse(model.DeletedAt.IsZero(), modelapi.ModelStatus_StatusInUse, modelapi.ModelStatus_StatusDeleted)

	return m
}

func convertModelProvider(provider modelmgr.ModelProvider) *modelapi.ModelProvider {
	return &modelapi.ModelProvider{
		Name: &modelapi.I18nText{
			ZhCn: provider.Name.ZhCn,
			EnUs: provider.Name.EnUs,
		},
		IconURI: provider.IconURI,
		IconURL: provider.IconURL,
		Description: &modelapi.I18nText{
			ZhCn: provider.Description.ZhCn,
			EnUs: provider.Description.EnUs,
		},
		ModelClass: convertModelClass(provider.ModelClass),
	}
}

func convertModelDisplayInfo(displayInfo modelmgr.DisplayInfo) *modelapi.DisplayInfo {
	if displayInfo.Description == nil {
		return &modelapi.DisplayInfo{
			Name:         displayInfo.Name,
			OutputTokens: displayInfo.OutputTokens,
			MaxTokens:    displayInfo.MaxTokens,
		}
	}
	return &modelapi.DisplayInfo{
		Name: displayInfo.Name,
		Description: &modelapi.I18nText{
			ZhCn: displayInfo.Description.ZhCn,
			EnUs: displayInfo.Description.EnUs,
		},
		OutputTokens: displayInfo.OutputTokens,
		MaxTokens:    displayInfo.MaxTokens,
	}
}

func convertModelCapability(capability modelmgr.ModelAbility) *developer_api.ModelAbility {
	return &developer_api.ModelAbility{
		CotDisplay:         ptr.Of(capability.CotDisplay),
		FunctionCall:       ptr.Of(capability.FunctionCall),
		ImageUnderstanding: ptr.Of(capability.ImageUnderstanding),
		VideoUnderstanding: ptr.Of(capability.VideoUnderstanding),
		AudioUnderstanding: ptr.Of(capability.AudioUnderstanding),
		SupportMultiModal:  ptr.Of(capability.SupportMultiModal),
		PrefillResp:        ptr.Of(capability.PrefillResp),
	}
}

func convertModelConnection(connection entity.Connection) *modelapi.Connection {
	return &modelapi.Connection{
		BaseConnInfo: convertBaseConnectionInfo(connection.BaseConnInfo),
		Openai:       convertOpenAIConnInfo(connection.Openai),
		Deepseek:     convertDeepseekConnInfo(connection.Deepseek),
		Gemini:       convertGeminiConnInfo(connection.Gemini),
		Qwen:         convertQwenConnInfo(connection.Qwen),
		Ollama:       convertOllamaConnInfo(connection.Ollama),
		Claude:       convertClaudeConnInfo(connection.Claude),
	}
}
func convertModelParameters(parameters []modelmgr.ModelParameter) []*developer_api.ModelParameter {
	var modelParameterList []*developer_api.ModelParameter
	for _, parameter := range parameters {
		modelParameterList = append(modelParameterList, &developer_api.ModelParameter{
			Name:       parameter.Name,
			Label:      parameter.Label,
			Desc:       parameter.Desc,
			Min:        parameter.Min,
			Max:        parameter.Max,
			Precision:  parameter.Precision,
			DefaultVal: convertModelParamDefaultValue(parameter.DefaultVal),
			Options:    convertOptions(parameter.Options),
			ParamClass: convertModelParamClass(parameter.ParamClass),
			Type:       developer_api.ModelParamType(parameter.Type),
		})
	}
	return modelParameterList
}

func convertModelParamDefaultValue(defaultVal *modelmgr.ModelParamDefaultValue) *developer_api.ModelParamDefaultValue {
	if defaultVal == nil {
		return nil
	}
	return &developer_api.ModelParamDefaultValue{
		DefaultVal: defaultVal.DefaultVal,
		Creative:   ptr.Of(defaultVal.Creative),
		Balance:    ptr.Of(defaultVal.Balance),
		Precise:    ptr.Of(defaultVal.Precise),
	}
}

func convertOptions(options []*modelmgr.Option) []*developer_api.Option {
	optionList := make([]*developer_api.Option, 0, len(options))
	for _, option := range options {
		optionList = append(optionList, &developer_api.Option{
			Value: option.Value,
			Label: option.Label,
		})
	}
	return optionList
}

func convertModelParamClass(paramClass *modelmgr.ModelParamClass) *developer_api.ModelParamClass {
	if paramClass == nil {
		return nil
	}
	return &developer_api.ModelParamClass{
		ClassID: paramClass.ClassID,
		Label:   paramClass.Label,
	}
}

func convertModelClass(modelClass modelmgr.ModelClass) developer_api.ModelClass {
	switch modelClass {
	case modelmgr.ModelClass_Ollama:
		return developer_api.ModelClass_Ollama
	case modelmgr.ModelClass_GPT:
		return developer_api.ModelClass_GPT
	case modelmgr.ModelClass_DeepSeek:
		return developer_api.ModelClass_DeepSeek
	case modelmgr.ModelClass_Gemini:
		return developer_api.ModelClass_Gemini
	case modelmgr.ModelClass_QWen:
		return developer_api.ModelClass_QWen
	case modelmgr.ModelClass_Claude:
		return developer_api.ModelClass_Claude
	default:
		return developer_api.ModelClass_Other
	}
}

func convertBaseConnectionInfo(baseConnInfo *modelmgr.BaseConnectionInfo) *modelapi.BaseConnectionInfo {
	if baseConnInfo == nil {
		return nil
	}
	return &modelapi.BaseConnectionInfo{
		BaseURL:      baseConnInfo.BaseURL,
		APIKey:       baseConnInfo.APIKey,
		Model:        baseConnInfo.Model,
		ThinkingType: modelapi.ThinkingType(baseConnInfo.ThinkingType),
	}
}

func convertOpenAIConnInfo(openaiConnInfo *modelmgr.OpenAIConnInfo) *modelapi.OpenAIConnInfo {
	if openaiConnInfo == nil {
		return nil
	}
	return &modelapi.OpenAIConnInfo{
		ByAzure:    openaiConnInfo.ByAzure,
		APIVersion: openaiConnInfo.APIVersion,
	}
}

func convertDeepseekConnInfo(deepseekConnInfo *modelmgr.DeepseekConnInfo) *modelapi.DeepseekConnInfo {
	if deepseekConnInfo == nil {
		return nil
	}
	return &modelapi.DeepseekConnInfo{}
}

func convertGeminiConnInfo(geminiConnInfo *modelmgr.GeminiConnInfo) *modelapi.GeminiConnInfo {
	if geminiConnInfo == nil {
		return nil
	}
	return &modelapi.GeminiConnInfo{
		Backend:  geminiConnInfo.Backend,
		Project:  geminiConnInfo.Project,
		Location: geminiConnInfo.Location,
	}
}

func convertQwenConnInfo(qwenConnInfo *modelmgr.QwenConnInfo) *modelapi.QwenConnInfo {
	if qwenConnInfo == nil {
		return nil
	}
	return &modelapi.QwenConnInfo{}
}

func convertOllamaConnInfo(ollamaConnInfo *modelmgr.OllamaConnInfo) *modelapi.OllamaConnInfo {
	if ollamaConnInfo == nil {
		return nil
	}
	return &modelapi.OllamaConnInfo{}
}

func convertClaudeConnInfo(claudeConnInfo *modelmgr.ClaudeConnInfo) *modelapi.ClaudeConnInfo {
	if claudeConnInfo == nil {
		return nil
	}
	return &modelapi.ClaudeConnInfo{}
}

func encryptConn(ctx context.Context, conn entity.Connection) (entity.Connection, error) {
	// encrypt conn if you need
	return conn, nil
}

func decryptConn(ctx context.Context, conn entity.Connection) (entity.Connection, error) {
	return conn, nil
}
