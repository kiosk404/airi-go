package model

import (
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/kiosk404/airi-go/backend/api/model/llm/domain/manage"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/pkg/errors"
)

type ParamConfig manage.ParamConfig
type Ability manage.Ability

//type ProtocolConfig manage.ProtocolConfig

type ProtocolConfig struct {
	manage.ProtocolConfig
}

type Protocol manage.Protocol
type ScenarioConfig manage.ScenarioConfig
type ParamSchema manage.ParamSchema

type AbilityMultiModal manage.AbilityMultiModal

type Model struct {
	ID              int64  `json:"id" yaml:"id" mapstructure:"id"`       // id
	Name            string `json:"name" yaml:"name" mapstructure:"name"` // 模型展示名称
	Desc            string `json:"desc" yaml:"desc" mapstructure:"desc"` // 模型描述
	IconURI         string `json:"icon_uri,omitempty" yaml:"icon_uri" mapstructure:"icon_uri"`
	IconURL         string `json:"icon_url,omitempty" yaml:"icon_url" mapstructure:"icon_url"`
	EnableBase64URL bool   `json:"enable_base64_url" yaml:"enable_base64_url" mapstructure:"enable_base64_url"`

	Ability         *Ability                     `json:"ability" yaml:"ability" mapstructure:"ability"`                            // 模型能力
	Protocol        *Protocol                    `json:"protocol" yaml:"protocol" mapstructure:"protocol"`                         // 该模型的协议类型，如qwen/deepseek/openai等
	ProtocolConfig  *ProtocolConfig              `json:"protocol_config" yaml:"protocol_config" mapstructure:"protocol_config"`    // 该模型的协议配置
	ScenarioConfigs map[Scenario]*ScenarioConfig `json:"scenario_configs" yaml:"scenario_configs" mapstructure:"scenario_configs"` // 该模型的场景配置
	ParamConfig     ParamConfig                  `json:"param_config" yaml:"param_config" mapstructure:"param_config"`             // 该模型的参数配置
	CommonParam     CommonParam                  `json:"common_param" yaml:"common_param" mapstructure:"common_param"`             // 该模型的通用参数配置
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) InitDefault() {
}

func (m *Model) GetModelID() (v int64) {
	return m.ID
}

func (m *Model) GetIconURI() (v string) {
	return m.IconURI
}

func (m *Model) GetIconURL() (v string) {
	return m.IconURL
}

func (m *Model) Valid() error {
	if m == nil {
		return errors.Errorf("model is nil")
	}
	if m.ID == 0 {
		return errors.Errorf("model id is zero")
	}
	if m.Name == "" {
		return errors.Errorf("model name is empty")
	}
	if err := ValidAbility(m.Ability); err != nil {
		return err
	}
	if err := ValidProtocolConfig(m.Protocol); err != nil {
		return err
	}
	return nil
}

func ValidAbility(ability *Ability) error {
	if ability == nil {
		return nil
	}
	if ptr.From(ability.MultiModal) {
		if ability.AbilityMultiModal == nil {
			return errors.Errorf("multi modal is true but ability multi modal is nil")
		}
		if ptr.From(ability.AbilityMultiModal.Image) {
			if !ability.AbilityMultiModal.GetImage() {
				return errors.Errorf("multi modal Image is true but ability multi modal ability image is nil")
			}
		}
	}
	return nil
}

func ValidProtocolConfig(protocol *Protocol) error {
	if protocol == nil {
		return errors.Errorf("protocol is empty")
	}
	return nil
}

func (m *Model) GetModelName() string {
	if m == nil {
		return ""
	}
	return ptr.From(m.ProtocolConfig.Model)
}

func (m *Model) GetProtocol() Protocol {
	if m == nil {
		return ""
	}
	return ptr.From(m.Protocol)
}

func (m *Model) GetProtocolConfig() ProtocolConfig {
	if m == nil {
		return ProtocolConfig{}
	}
	return ptr.From(m.ProtocolConfig)
}

func (m *Model) SupportMultiModalInput() bool {
	if m == nil || m.Ability == nil {
		return false
	}
	return ptr.From(m.Ability.MultiModal)
}

func (m *Model) SupportImageURL() (bool, int64) {
	if m == nil || m.Ability == nil || m.Ability.AbilityMultiModal == nil || m.Ability.AbilityMultiModal.Image == nil {
		return false, 0
	}
	return true, 10
}

func (m *Model) SupportImageBinary() (bool, int64, int64) {
	if m == nil || m.Ability == nil || m.Ability.AbilityMultiModal == nil || m.Ability.AbilityMultiModal.Image == nil {
		return false, 0, 0
	}
	return true, 10, 1024
}

func (m *Model) SupportFunctionCall() bool {
	if m == nil || m.Ability == nil {
		return false
	}
	return ptr.From(m.Ability.FunctionCall)
}

func (m *Model) Available(scenario *Scenario) bool {
	// 默认都是available
	if scenario == nil || m.ScenarioConfigs == nil {
		return true
	}
	scenarioConfig, ok := m.ScenarioConfigs[*scenario]
	if !ok || scenarioConfig == nil {
		return true
	}
	return !ptr.From(scenarioConfig.Unavailable)
}

func (m *Model) GetScenarioConfig(scenario *Scenario) *ScenarioConfig {
	if m.ScenarioConfigs == nil {
		return nil
	}
	if scenario == nil {
		return m.ScenarioConfigs[ScenarioDefault]
	}
	cfg, ok := m.ScenarioConfigs[*scenario]
	if ok && cfg != nil {
		return cfg
	}
	return m.ScenarioConfigs[ScenarioDefault]
}

func (m *Model) FindParameter(name ParameterName) (*ParamSchema, bool) {
	if len(m.ParamConfig.ParamSchemas) == 0 {
		return nil, false
	}
	for _, param := range m.ParamConfig.ParamSchemas {
		if ptr.From(param.Name) == name.String() {
			return ptr.PtrConvert[manage.ParamSchema, ParamSchema](param), true
		}
	}
	return nil, false
}

type CommonParam struct {
	MaxTokens        *int     `json:"max_tokens,omitempty" yaml:"max_tokens" mapstructure:"max_tokens"`
	Temperature      *float32 `json:"temperature,omitempty" yaml:"temperature" mapstructure:"temperature"`
	TopP             *float32 `json:"top_p,omitempty" yaml:"top_p" mapstructure:"top_p"`
	TopK             *int     `json:"top_k,omitempty" yaml:"top_k" mapstructure:"top_k"`
	Stop             []string `json:"stop,omitempty" yaml:"stop" mapstructure:"stop"`
	FrequencyPenalty *float32 `json:"frequency_penalty,omitempty" yaml:"frequency_penalty" mapstructure:"frequency_penalty"`
	PresencePenalty  *float32 `json:"presence_penalty,omitempty" yaml:"presence_penalty" mapstructure:"presence_penalty"`
	EnableThinking   *bool    `json:"enable_thinking,omitempty" yaml:"enable_thinking,omitempty"`
}

type DefaultValue map[DefaultType]string

func (p *ParamSchema) GetFloat(tp DefaultType) (float64, error) {
	if ptr.FromPtrConvert[manage.ParamType, ParamType](p.Type) != ParamTypeFloat {
		return 0, fmt.Errorf("unexpected paramerter type, name=%v, expect=%v, given=%v", p.Name, ParamTypeFloat, p.Type)
	}
	mtp := manage.DefaultType(tp)
	if tp == DefaultTypeDefault && p.DefaultValue[mtp] == "" {
		tp = DefaultTypeDefault
	}
	val, ok := p.DefaultValue[mtp]
	if !ok {
		return 0, fmt.Errorf("unexpected default type, name=%v, type=%v", p.Name, tp)
	}

	return strconv.ParseFloat(val, 64)
}

func (p *ParamSchema) GetInt(tp DefaultType) (int64, error) {
	if ptr.FromPtrConvert[manage.ParamType, ParamType](p.Type) != ParamTypeInt {
		return 0, fmt.Errorf("unexpected paramerter type, name=%v, expect=%v, given=%v", p.Name, ParamTypeInt, p.Type)
	}
	mtp := manage.DefaultType(tp)
	if tp != DefaultTypeDefault && p.DefaultValue[mtp] == "" {
		tp = DefaultTypeDefault
	}
	val, ok := p.DefaultValue[mtp]
	if !ok {
		return 0, fmt.Errorf("unexpected default type, name=%v, type=%v", p.Name, tp)
	}
	return strconv.ParseInt(val, 10, 64)
}

func (p *ParamSchema) GetBool(tp DefaultType) (bool, error) {
	if ptr.FromPtrConvert[manage.ParamType, ParamType](p.Type) != ParamTypeBoolean {
		return false, fmt.Errorf("unexpected paramerter type, name=%v, expect=%v, given=%v",
			p.Name, ParamTypeBoolean, p.Type)
	}
	mtp := manage.DefaultType(tp)

	if tp != DefaultTypeDefault && p.DefaultValue[mtp] == "" {
		tp = DefaultTypeDefault
	}
	val, ok := p.DefaultValue[mtp]
	if !ok {
		return false, fmt.Errorf("unexpected default type, name=%v, type=%v", p.Name, tp)
	}
	return strconv.ParseBool(val)
}

func (p *ParamSchema) GetString(tp DefaultType) (string, error) {
	if tp != DefaultTypeDefault && p.DefaultValue[manage.DefaultType(tp)] == "" {
		tp = DefaultTypeDefault
	}
	mtp := manage.DefaultType(tp)

	val, ok := p.DefaultValue[mtp]
	if !ok {
		return "", fmt.Errorf("unexpected default type, name=%v, type=%v", p.Name, tp)
	}
	return val, nil
}

type ParamOption struct {
	Value string `json:"value" yaml:"value" mapstructure:"value"`
	Label string `json:"label" yaml:"label" mapstructure:"label"`
}

// type Protocol string
const (
	ProtocolOpenAI   Protocol = "openai"
	ProtocolDeepseek Protocol = "deepseek"
	ProtocolClaude   Protocol = "claude"
	ProtocolOllama   Protocol = "ollama"
	ProtocolGemini   Protocol = "gemini"
	ProtocolQwen     Protocol = "qwen"
)

func (p *ParamConfig) GetCommonParamDefaultVal() CommonParam {
	rawDf := p.GetDefaultVal([]ParameterName{MaxTokens, Temperature, TopP, TopK, FrequencyPenalty, PresencePenalty, Stop})
	cp := CommonParam{}
	if rawDf == nil {
		return cp
	}
	if rawDf[MaxTokens] != "" {
		maxTokens, _ := strconv.ParseInt(rawDf[MaxTokens], 10, 32)
		cp.MaxTokens = ptr.Of(int(maxTokens))
	}
	if rawDf[Temperature] != "" {
		temperature, _ := strconv.ParseFloat(rawDf[Temperature], 32)
		cp.Temperature = ptr.Of(float32(temperature))
	}
	if rawDf[TopP] != "" {
		topP, _ := strconv.ParseFloat(rawDf[TopP], 32)
		cp.TopP = ptr.Of(float32(topP))
	}
	if rawDf[TopK] != "" {
		topK, _ := strconv.ParseInt(rawDf[TopK], 10, 32)
		cp.TopK = ptr.Of(int(topK))
	}
	if rawDf[Stop] != "" {
		var stop []string
		_ = sonic.UnmarshalString(rawDf[Stop], &stop)
		cp.Stop = stop
	}
	if rawDf[FrequencyPenalty] != "" {
		frequencyPenalty, _ := strconv.ParseFloat(rawDf[FrequencyPenalty], 32)
		cp.FrequencyPenalty = ptr.Of(float32(frequencyPenalty))
	}
	if rawDf[PresencePenalty] != "" {
		presencePenalty, _ := strconv.ParseFloat(rawDf[PresencePenalty], 32)
		cp.PresencePenalty = ptr.Of(float32(presencePenalty))
	}
	return cp
}

func (p *ParamConfig) GetDefaultVal(params []ParameterName) map[ParameterName]string {
	if p == nil || len(p.ParamSchemas) == 0 {
		return nil
	}
	res := make(map[ParameterName]string)
	for _, param := range params {
		for _, ps := range p.ParamSchemas {
			if param == ptr.FromPtrConvert[string, ParameterName](ps.Name) {
				res[param] = ps.DefaultValue[manage.DefaultTypeValue]
			}
		}
	}
	return res
}

func (a *Ability) GetMultiModal() bool {
	if a.MultiModal == nil {
		return false
	}
	return ptr.From(a.MultiModal)
}

type ListModelsRequest struct {
	Scenario  *Scenario
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
