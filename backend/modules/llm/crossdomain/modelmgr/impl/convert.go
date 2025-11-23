package impl

import (
	"strconv"

	"github.com/bytedance/gg/gslice"
	"github.com/kiosk404/airi-go/backend/api/model/llm/domain/common"
	modelmgr "github.com/kiosk404/airi-go/backend/api/model/llm/domain/manage"
	"github.com/kiosk404/airi-go/backend/api/model/llm/manage"
	"github.com/kiosk404/airi-go/backend/api/model/llm/runtime"
	"github.com/kiosk404/airi-go/backend/modules/llm/application/convertor"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

func convertListModelsRequest(req *model.ListModelsRequest) *manage.ListModelsRequest {
	return &manage.ListModelsRequest{
		Scenario:  (*common.Scenario)(req.Scenario),
		PageSize:  ptr.Of(ptr.OfConvert[int64, int32](req.PageSize)),
		PageToken: ptr.Of(strconv.FormatInt(req.PageToken, 10)),
	}
}

func convertListModelsResponse(resp *manage.ListModelsResponse) *model.ListModelsResponse {
	return &model.ListModelsResponse{
		Models:        gslice.Map(resp.Models, convertModel),
		HasMore:       resp.HasMore,
		NextPageToken: resp.NextPageToken,
		Total:         resp.Total,
	}
}

func convertGetModelRequest(req *model.GetModelRequest) *manage.GetModelRequest {
	return &manage.GetModelRequest{
		ModelID: ptr.Of(req.ModelID),
	}
}

func convertGetModelResponse(resp *manage.GetModelResponse) *model.GetModelResponse {
	return &model.GetModelResponse{
		Model: convertModel(resp.Model),
	}
}

func convertChatRequest(req *model.ChatRequest) *runtime.ChatRequest {
	return &runtime.ChatRequest{}
}

func convertChatResponse(resp *runtime.ChatResponse) *model.ChatResponse {
	return &model.ChatResponse{}
}

func convertModel(m *modelmgr.Model) *entity.Model {
	return &entity.Model{
		ID:      ptr.From(m.ModelID),
		Name:    ptr.From(m.Name),
		Desc:    ptr.From(m.Desc),
		IconURI: ptr.From(m.IconURI),
		IconURL: ptr.From(m.IconURL),

		Ability:         convertor.AbilityDTO2DO(m.Ability),
		Protocol:        convertor.ProtocolDTO2DO(m.Protocol),
		ProtocolConfig:  convertor.ProtocolConfigDTO2DO(m.ProtocolConfig),
		ScenarioConfigs: convertor.ScenarioConfigMapDTO2DO(m.ScenarioConfigs),
		ParamConfig:     convertor.ParamConfigDTO2DO(m.ParamConfig),
	}
}
