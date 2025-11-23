package entity

import (
	"fmt"

	"github.com/kiosk404/airi-go/backend/api/model/llm/domain/manage"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
)

type Model = model.Model

type Ability = manage.Ability
type AbilityMultiModal = manage.AbilityMultiModal
type ProtocolConfig = manage.ProtocolConfig
type CommonParam = model.CommonParam

var ProtocolOpenAI = model.ProtocolOpenAI
var ProtocolDeepseek = model.ProtocolDeepseek
var ProtocolClaude = model.ProtocolClaude
var ProtocolOllama = model.ProtocolOllama
var ProtocolGemini = model.ProtocolGemini
var ProtocolQwen = model.ProtocolQwen

type ProtocolConfigOpenAI = manage.ProtocolConfigOpenAI
type ProtocolConfigClaude = manage.ProtocolConfigClaude
type ProtocolConfigDeepSeek = manage.ProtocolConfigDeepSeek
type ProtocolConfigGemini = manage.ProtocolConfigGemini
type ProtocolConfigQwen = manage.ProtocolConfigQwen
type ProtocolConfigOllama = manage.ProtocolConfigOllama

type ScenarioConfig = manage.ScenarioConfig
type ParamSchema = manage.ParamSchema
type ParamOption = manage.ParamOption
type ParamConfig = manage.ParamConfig
type DefaultType = manage.DefaultType

type ListModelsRequest struct {
	Scenario  *model.Scenario
	PageToken int64
	PageSize  int64
}

func NewListModelsRequest() *ListModelsRequest {
	return &ListModelsRequest{}
}

func (p *ListModelsRequest) InitDefault() {
}

func (p *ListModelsRequest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ListModelsRequest(%+v)", *p)
}

type GetModelRequest struct {
	ModelID int64
}

func NewGetModelRequest() *GetModelRequest {
	return &GetModelRequest{}
}

func (p *GetModelRequest) InitDefault() {
}

func (p *GetModelRequest) GetModelID() (v int64) {
	return p.ModelID
}

func (p *GetModelRequest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetModelRequest(%+v)", *p)
}

type MGetModelReq struct {
	ModelIDs []int64
}

type ListModelsResponse struct {
	Models        []*Model
	HasMore       *bool   `thrift:"has_more,127,optional" json:"has_more,omitempty"`
	NextPageToken *string `thrift:"next_page_token,128,optional" json:"next_page_token,omitempty"`
	Total         *int32  `thrift:"total,129,optional" json:"total,omitempty"`
}

func NewListModelsResponse() *ListModelsResponse {
	return &ListModelsResponse{}
}

func (p *ListModelsResponse) InitDefault() {
}

func (p *ListModelsResponse) GetModels() (v []*Model) {
	return p.Models
}

func (p *ListModelsResponse) GetHasMore() (v bool) {
	return *p.HasMore
}

func (p *ListModelsResponse) GetNextPageToken() (v string) {
	return *p.NextPageToken
}

func (p *ListModelsResponse) GetTotal() (v int32) {
	return *p.Total
}

func (p *ListModelsResponse) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ListModelsResponse(%+v)", *p)
}

type GetModelResponse struct {
	Model *Model
}

func NewGetModelResponse() *GetModelResponse {
	return &GetModelResponse{}
}

func (p *GetModelResponse) GetModel() (v *Model) {
	return p.Model
}
