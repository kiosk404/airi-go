package convertor

import (
	druntime "github.com/kiosk404/airi-go/backend/api/model/llm/domain/runtime"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
)

func ModelAndTools2OptionDOs(modelCfg *druntime.ModelConfig, tools []*druntime.Tool) []model.Option {
	var opts []model.Option
	if modelCfg != nil {
		if modelCfg.Temperature != nil {
			opts = append(opts, model.WithTemperature(float32(*modelCfg.Temperature)))
		}
		if modelCfg.MaxTokens != nil {
			opts = append(opts, model.WithMaxTokens(int(*modelCfg.MaxTokens)))
		}
		if modelCfg.TopP != nil {
			opts = append(opts, model.WithTopP(float32(*modelCfg.TopP)))
		}
		if len(modelCfg.Stop) > 0 {
			opts = append(opts, model.WithStop(modelCfg.Stop))
		}
		if modelCfg.ToolChoice != nil {
			opts = append(opts, model.WithToolChoice(ToolChoiceDTO2DO(modelCfg.ToolChoice)))
		}
		if modelCfg.ResponseFormat != nil {
			opts = append(opts, model.WithResponseFormat(ResponseFormatDTO2DO(modelCfg.ResponseFormat)))
		}
		if modelCfg.TopK != nil {
			opts = append(opts, model.WithTopK(ptr.Of(int(ptr.From(modelCfg.TopK)))))
		}
		if modelCfg.PresencePenalty != nil {
			opts = append(opts, model.WithPresencePenalty(float32(*modelCfg.PresencePenalty)))
		}
		if modelCfg.FrequencyPenalty != nil {
			opts = append(opts, model.WithFrequencyPenalty(float32(*modelCfg.FrequencyPenalty)))
		}
	}
	if len(tools) > 0 {
		toolsDTO := slices.Map(tools, func(t *druntime.Tool, _ int) *entity.ToolInfo {
			return ToolDTO2DO(t)
		})
		opts = append(opts, model.WithTools(toolsDTO))
	}
	return opts
}

func ResponseFormatDTO2DO(r *druntime.ResponseFormat) *model.ResponseFormat {
	if r == nil {
		return nil
	}

	return &model.ResponseFormat{
		Type: convertResponseFormat(r.GetType()),
	}
}

func ToolsDTO2DO(ts []*druntime.Tool) []*entity.ToolInfo {
	return slices.Map(ts, func(t *druntime.Tool, _ int) *entity.ToolInfo {
		return ToolDTO2DO(t)
	})
}

func ToolDTO2DO(t *druntime.Tool) *entity.ToolInfo {
	if t == nil {
		return nil
	}
	return &entity.ToolInfo{
		Name:        t.GetName(),
		Desc:        t.GetDesc(),
		ToolDefType: model.ToolDefType(t.GetDefType()),
		Def:         t.GetDef(),
	}
}

func ToolChoiceDTO2DO(tc *druntime.ToolChoice) *entity.ToolChoice {
	if tc == nil {
		return nil
	}
	return ptr.Of(entity.ToolChoice(*tc))
}

func convertResponseFormat(format druntime.ResponseFormatType) model.ResponseFormatType {
	var modelResponseFormat model.ResponseFormatType
	switch format {
	case druntime.ResponseFormatText:
		modelResponseFormat = model.ResponseFormatText
	case druntime.ResponseFormatJSONObject:
		modelResponseFormat = model.ResponseFormatJSON
	case druntime.ResourceFormatMarkdown:
		modelResponseFormat = model.ResponseFormatMarkdown
	default:
		modelResponseFormat = model.ResponseFormatText
	}

	return modelResponseFormat
}
