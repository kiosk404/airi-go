package convertor

import (
	"github.com/kiosk404/airi-go/backend/api/model/llm/domain/common"
	"github.com/kiosk404/airi-go/backend/api/model/llm/domain/manage"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
)

func ModelsDO2DTO(models []*entity.Model, mask bool) []*manage.Model {
	return slices.Map(models, func(model *entity.Model, _ int) *manage.Model {
		return ModelDO2DTO(model, mask)
	})
}

func ModelDO2DTO(m *entity.Model, mask bool) *manage.Model {
	if m == nil {
		return nil
	}

	var pc *manage.ProtocolConfig
	if !mask {
		pc = ptr.Of(m.GetProtocolConfig().ProtocolConfig)
	}
	return &manage.Model{
		ModelID:         ptr.Of(m.ID),
		Name:            ptr.Of(m.Name),
		Desc:            ptr.Of(m.Desc),
		Ability:         AbilityDO2DTO(m.Ability),
		Protocol:        ProtocolDO2DTO(m.Protocol),
		ProtocolConfig:  pc,
		ScenarioConfigs: ScenarioConfigMapDO2DTO(m.ScenarioConfigs),
		ParamConfig:     ParamConfigDO2DTO(m.ParamConfig),
	}
}

func ModelDTO2DO(m *manage.Model) *model.Model {
	if m == nil {
		return nil
	}

	return &model.Model{
		ID:       ptr.From(m.ModelID),
		Name:     ptr.From(m.Name),
		Desc:     ptr.From(m.Desc),
		Ability:  AbilityDTO2DO(m.Ability),
		Protocol: ProtocolDTO2DO(m.Protocol),
		ProtocolConfig: &model.ProtocolConfig{
			ProtocolConfig: ptr.From(m.ProtocolConfig),
		},
		ParamConfig: ParamConfigDTO2DO(m.ParamConfig),
	}
}

func AbilityDO2DTO(a *model.Ability) *manage.Ability {
	if a == nil {
		return nil
	}
	return &manage.Ability{
		MaxContextTokens:  a.MaxContextTokens,
		MaxInputTokens:    a.MaxInputTokens,
		MaxOutputTokens:   a.MaxOutputTokens,
		FunctionCall:      a.FunctionCall,
		JSONMode:          a.JSONMode,
		MultiModal:        a.MultiModal,
		AbilityMultiModal: AbilityMultiModalDO2DTO(a.AbilityMultiModal),
	}
}

func AbilityDTO2DO(a *manage.Ability) *model.Ability {
	if a == nil {
		return nil
	}

	return &model.Ability{
		MaxContextTokens:  a.MaxContextTokens,
		MaxInputTokens:    a.MaxInputTokens,
		MaxOutputTokens:   a.MaxOutputTokens,
		FunctionCall:      a.FunctionCall,
		JSONMode:          a.JSONMode,
		MultiModal:        a.MultiModal,
		AbilityMultiModal: AbilityMultiModalDTO2DO(a.AbilityMultiModal),
	}
}

func AbilityMultiModalDO2DTO(a *entity.AbilityMultiModal) *manage.AbilityMultiModal {
	if a == nil {
		return nil
	}
	return &manage.AbilityMultiModal{
		Image:        a.Image,
		Video:        a.Video,
		Audio:        a.Audio,
		FunctionCall: a.FunctionCall,
		PrefillResp:  a.PrefillResp,
	}
}

func AbilityMultiModalDTO2DO(a *manage.AbilityMultiModal) *entity.AbilityMultiModal {
	if a == nil {
		return nil
	}
	return &entity.AbilityMultiModal{
		Image:        a.Image,
		Video:        a.Video,
		Audio:        a.Audio,
		FunctionCall: a.FunctionCall,
		PrefillResp:  a.PrefillResp,
	}
}

func ProtocolDO2DTO(p *model.Protocol) *manage.Protocol {
	return ptr.Of(manage.Protocol(*p))
}

func ProtocolDTO2DO(a *manage.Protocol) *model.Protocol {
	return ptr.Of(model.Protocol(*a))
}

func ProtocolConfigDO2DTO(p *entity.ProtocolConfig) *model.ProtocolConfig {
	if p == nil {
		return nil
	}

	return &model.ProtocolConfig{
		ProtocolConfig: ptr.From(p),
	}
}

func ProtocolConfigDTO2DO(p *manage.ProtocolConfig) *model.ProtocolConfig {
	if p == nil {
		return &model.ProtocolConfig{}
	}

	return &model.ProtocolConfig{
		ProtocolConfig: manage.ProtocolConfig{
			BaseURL:                p.BaseURL,
			APIKey:                 p.APIKey,
			Model:                  p.Model,
			ProtocolConfigOpenai:   ProtocolConfigOpenaiDTO2DO(p.ProtocolConfigOpenai),
			ProtocolConfigClaude:   ProtocolConfigClaudeDTO2DO(p.ProtocolConfigClaude),
			ProtocolConfigDeepseek: ProtocolConfigDeepSeekDTO2DO(p.ProtocolConfigDeepseek),
			ProtocolConfigGemini:   ProtocolConfigGeminiDTO2DO(p.ProtocolConfigGemini),
			ProtocolConfigQwen:     ProtocolConfigQwenDTO2DO(p.ProtocolConfigQwen),
			ProtocolConfigOllama:   ProtocolConfigOllamaDTO2DO(p.ProtocolConfigOllama),
		},
	}
}

func ProtocolConfigOpenaiDO2DTO(p *entity.ProtocolConfigOpenAI) *manage.ProtocolConfigOpenAI {
	if p == nil {
		return nil
	}
	return &manage.ProtocolConfigOpenAI{
		ByAzure:                  p.ByAzure,
		APIVersion:               p.APIVersion,
		ResponseFormatType:       p.ResponseFormatType,
		ResponseFormatJSONSchema: p.ResponseFormatJSONSchema,
	}
}

func ProtocolConfigOpenaiDTO2DO(p *manage.ProtocolConfigOpenAI) *entity.ProtocolConfigOpenAI {
	if p == nil {
		return nil
	}
	return &entity.ProtocolConfigOpenAI{
		ByAzure:                  p.ByAzure,
		APIVersion:               p.APIVersion,
		ResponseFormatType:       p.ResponseFormatType,
		ResponseFormatJSONSchema: p.ResponseFormatJSONSchema,
	}
}

func ProtocolConfigClaudeDO2DTO(p *entity.ProtocolConfigClaude) *manage.ProtocolConfigClaude {
	if p == nil {
		return nil
	}
	return &manage.ProtocolConfigClaude{
		ByBedrock:       p.ByBedrock,
		AccessKey:       p.AccessKey,
		SecretAccessKey: p.SecretAccessKey,
		SessionToken:    p.SessionToken,
		Region:          p.Region,
	}
}

func ProtocolConfigClaudeDTO2DO(p *manage.ProtocolConfigClaude) *entity.ProtocolConfigClaude {
	if p == nil {
		return nil
	}
	return &entity.ProtocolConfigClaude{
		ByBedrock:       p.ByBedrock,
		AccessKey:       p.AccessKey,
		SecretAccessKey: p.SecretAccessKey,
		SessionToken:    p.SessionToken,
		Region:          p.Region,
	}
}

func ProtocolConfigDeepSeekDTO2DO(p *manage.ProtocolConfigDeepSeek) *entity.ProtocolConfigDeepSeek {
	if p == nil {
		return nil
	}
	return &entity.ProtocolConfigDeepSeek{ResponseFormatType: p.ResponseFormatType}
}

func ProtocolConfigGeminiDTO2DO(p *manage.ProtocolConfigGemini) *entity.ProtocolConfigGemini {
	if p == nil {
		return nil
	}
	return &entity.ProtocolConfigGemini{
		Backend:         p.Backend,
		Project:         p.Project,
		Location:        p.Location,
		APIVersion:      p.APIVersion,
		TimeoutMs:       p.TimeoutMs,
		IncludeThoughts: p.IncludeThoughts,
		ThinkingBudget:  p.ThinkingBudget,
	}
}

func ProtocolConfigQwenDTO2DO(p *manage.ProtocolConfigQwen) *entity.ProtocolConfigQwen {
	if p == nil {
		return nil
	}
	return &entity.ProtocolConfigQwen{
		ResponseFormatType:       p.ResponseFormatType,
		ResponseFormatJSONSchema: p.ResponseFormatJSONSchema,
	}
}

func ProtocolConfigOllamaDTO2DO(p *manage.ProtocolConfigOllama) *entity.ProtocolConfigOllama {
	if p == nil {
		return nil
	}
	return &entity.ProtocolConfigOllama{
		Format:      p.Format,
		KeepAliveMs: p.KeepAliveMs,
	}
}

func ScenarioConfigMapDO2DTO(s map[model.Scenario]*model.ScenarioConfig) map[common.Scenario]*manage.ScenarioConfig {
	if s == nil {
		return nil
	}
	res := make(map[common.Scenario]*manage.ScenarioConfig)
	for k, v := range s {
		res[ScenarioDO2DTO(k)] = ScenarioConfigDO2DTO(v)
	}
	return res
}

func ScenarioConfigMapDTO2DO(s map[common.Scenario]*manage.ScenarioConfig) map[model.Scenario]*model.ScenarioConfig {
	if s == nil {
		return nil
	}
	res := make(map[model.Scenario]*model.ScenarioConfig)
	for k, v := range s {
		res[ScenarioDTO2DO(k)] = ScenarioConfigDTO2DO(v)
	}
	return res
}

func ScenarioConfigDO2DTO(s *model.ScenarioConfig) *manage.ScenarioConfig {
	if s == nil {
		return nil
	}
	return &manage.ScenarioConfig{
		Scenario:    s.Scenario,
		Quota:       s.Quota,
		Unavailable: s.Unavailable,
	}
}

func ScenarioConfigDTO2DO(s *manage.ScenarioConfig) *model.ScenarioConfig {
	if s == nil {
		return nil
	}
	return &model.ScenarioConfig{
		Scenario:    s.Scenario,
		Quota:       s.Quota,
		Unavailable: s.Unavailable,
	}
}

func ParamConfigDO2DTO(p model.ParamConfig) *manage.ParamConfig {
	return &manage.ParamConfig{
		ParamSchemas: slices.Map(p.ParamSchemas, func(s *entity.ParamSchema, _ int) *manage.ParamSchema {
			return ParamSchemaDO2DTO(s)
		}),
	}
}

func ParamConfigDTO2DO(p *manage.ParamConfig) model.ParamConfig {
	if p == nil {
		return model.ParamConfig{}
	}
	return model.ParamConfig{
		ParamSchemas: slices.Map(p.ParamSchemas, func(s *manage.ParamSchema, _ int) *entity.ParamSchema {
			return ParamSchemaDTO2DO(s)
		}),
	}
}

func ParamSchemaDO2DTO(ps *entity.ParamSchema) *manage.ParamSchema {
	if ps == nil {
		return nil
	}

	// Convert map[DefaultType]string to map[string]string
	var defaultValue map[string]string
	if ps.DefaultValue != nil {
		defaultValue = make(map[string]string)
		for k, v := range ps.DefaultValue {
			defaultValue[string(k)] = v
		}
	}

	return &manage.ParamSchema{
		Name:         ps.Name,
		Label:        ps.Label,
		Desc:         ps.Desc,
		Type:         ps.Type,
		Min:          ps.Min,
		Max:          ps.Max,
		DefaultValue: defaultValue,
		Options:      ParamOptionsDO2DTO(ps.Options),
	}
}

func ParamSchemaDTO2DO(ps *manage.ParamSchema) *entity.ParamSchema {
	if ps == nil {
		return nil
	}

	// Convert map[string]string to map[DefaultType]string
	var defaultValue map[entity.DefaultType]string
	if ps.DefaultValue != nil {
		var defaultValue = make(map[model.DefaultType]string)
		for k, v := range ps.DefaultValue {
			defaultValue[model.DefaultType(k)] = v
		}
	}

	return &entity.ParamSchema{
		Name:         ps.Name,
		Label:        ps.Label,
		Desc:         ps.Desc,
		Type:         ps.Type,
		Min:          ps.Min,
		Max:          ps.Max,
		DefaultValue: defaultValue,
		Options:      ParamOptionsDTO2DO(ps.Options),
	}
}

func ParamOptionsDO2DTO(os []*entity.ParamOption) []*manage.ParamOption {
	return slices.Map(os, func(o *entity.ParamOption, _ int) *manage.ParamOption {
		return ParamOptionDO2DTO(o)
	})
}

func ParamOptionsDTO2DO(os []*manage.ParamOption) []*entity.ParamOption {
	return slices.Map(os, func(o *manage.ParamOption, _ int) *entity.ParamOption {
		return ParamOptionDTO2DO(o)
	})
}

func ParamOptionDO2DTO(o *entity.ParamOption) *manage.ParamOption {
	if o == nil {
		return nil
	}
	return &manage.ParamOption{
		Value: o.Value,
		Label: o.Label,
	}
}

func ParamOptionDTO2DO(o *manage.ParamOption) *entity.ParamOption {
	if o == nil {
		return nil
	}
	return &entity.ParamOption{
		Value: o.Value,
		Label: o.Label,
	}
}
