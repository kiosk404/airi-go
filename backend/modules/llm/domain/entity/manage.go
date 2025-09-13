package entity

import (
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/pkg/errors"
	"google.golang.org/genai"
)

type Model struct {
	ID   int64  `json:"id" yaml:"id" mapstructure:"id"`       // id
	Name string `json:"name" yaml:"name" mapstructure:"name"` // 模型展示名称
	Desc string `json:"desc" yaml:"desc" mapstructure:"desc"` // 模型描述

	Ability         *Ability                     `json:"ability" yaml:"ability" mapstructure:"ability"`                            // 模型能力
	Protocol        Protocol                     `json:"protocol" yaml:"protocol" mapstructure:"protocol"`                         // 该模型的协议类型，如qwen/deepseek/openai等
	ProtocolConfig  ProtocolConfig               `json:"protocol_config" yaml:"protocol_config" mapstructure:"protocol_config"`    // 该模型的协议配置
	ScenarioConfigs map[Scenario]*ScenarioConfig `json:"scenario_configs" yaml:"scenario_configs" mapstructure:"scenario_configs"` // 该模型的场景配置
	ParamConfig     ParamConfig                  `json:"param_config" yaml:"param_config" mapstructure:"param_config"`             // 该模型的参数配置
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
	if err := m.Ability.ValidAbility(); err != nil {
		return err
	}
	if err := m.ProtocolConfig.ValidProtocolConfig(m.Protocol); err != nil {
		return err
	}
	return nil
}

func (a *Ability) ValidAbility() error {
	if a == nil {
		return nil
	}
	if a.MultiModal {
		if a.AbilityMultiModal == nil {
			return errors.Errorf("multi modal is true but ability multi modal is nil")
		}
		if a.AbilityMultiModal.Image {
			if a.AbilityMultiModal.AbilityImage == nil {
				return errors.Errorf("multi modal Image is true but ability multi modal ability image is nil")
			}
		}
	}
	return nil
}

func (p *ProtocolConfig) ValidProtocolConfig(protocol Protocol) error {
	if p == nil {
		return errors.Errorf("protocol config is nil")
	}
	if protocol == "" {
		return errors.Errorf("protocol is empty")
	}
	return nil
}

func (m *Model) GetModel() string {
	if m == nil {
		return ""
	}
	return m.ProtocolConfig.Model
}

func (m *Model) SupportMultiModalInput() bool {
	if m == nil || m.Ability == nil {
		return false
	}
	return m.Ability.MultiModal
}

func (m *Model) SupportImageURL() (bool, int64) {
	if m == nil || m.Ability == nil || m.Ability.AbilityMultiModal == nil || m.Ability.AbilityMultiModal.AbilityImage == nil {
		return false, 0
	}
	return m.Ability.AbilityMultiModal.AbilityImage.URLEnabled, m.Ability.AbilityMultiModal.AbilityImage.MaxImageCount
}

func (m *Model) SupportImageBinary() (bool, int64, int64) {
	if m == nil || m.Ability == nil || m.Ability.AbilityMultiModal == nil || m.Ability.AbilityMultiModal.AbilityImage == nil {
		return false, 0, 0
	}
	return m.Ability.AbilityMultiModal.AbilityImage.BinaryEnabled,
		m.Ability.AbilityMultiModal.AbilityImage.MaxImageCount, m.Ability.AbilityMultiModal.AbilityImage.MaxImageSize
}

func (m *Model) SupportFunctionCall() bool {
	if m == nil || m.Ability == nil {
		return false
	}
	return m.Ability.FunctionCall
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
	return !scenarioConfig.Unavailable
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

type Ability struct {
	MaxContextTokens  *int64             `json:"max_context_tokens" yaml:"max_context_tokens" mapstructure:"max_context_tokens"`
	MaxInputTokens    *int64             `json:"max_input_tokens" yaml:"max_input_tokens" mapstructure:"max_input_tokens"`
	MaxOutputTokens   *int64             `json:"max_output_tokens" yaml:"max_output_tokens" mapstructure:"max_output_tokens"`
	FunctionCall      bool               `json:"function_call" yaml:"function_call" mapstructure:"function_call"`
	JsonMode          bool               `json:"json_mode" yaml:"json_mode" mapstructure:"json_mode"`
	MultiModal        bool               `json:"multi_modal" yaml:"multi_modal" mapstructure:"multi_modal"`
	AbilityMultiModal *AbilityMultiModal `json:"ability_multi_modal" yaml:"ability_multi_modal" mapstructure:"ability_multi_modal"`
}

type AbilityMultiModal struct {
	Image        bool          `json:"image" yaml:"image" mapstructure:"image"`
	AbilityImage *AbilityImage `json:"ability_image" yaml:"ability_image" mapstructure:"ability_image"`
}

type AbilityImage struct {
	URLEnabled    bool  `json:"url_enabled" yaml:"url_enabled" mapstructure:"url_enabled"`
	BinaryEnabled bool  `json:"binary_enabled" yaml:"binary_enabled" mapstructure:"binary_enabled"`
	MaxImageSize  int64 `json:"max_image_size" yaml:"max_image_size" mapstructure:"max_image_size"`
	MaxImageCount int64 `json:"max_image_count" yaml:"max_image_count" mapstructure:"max_image_count"`
}

type ProtocolConfig struct {
	BaseURL                string                  `json:"base_url" yaml:"base_url" mapstructure:"base_url"`
	APIKey                 string                  `json:"api_key" yaml:"api_key" mapstructure:"api_key"`
	Model                  string                  `json:"model" yaml:"model" mapstructure:"model"`
	TimeoutMs              *int64                  `json:"timeout_ms" yaml:"timeout_ms" mapstructure:"timeout_ms"`
	ProtocolConfigOpenAI   *ProtocolConfigOpenAI   `json:"protocol_config_open_ai" yaml:"protocol_config_open_ai" mapstructure:"protocol_config_open_ai"`
	ProtocolConfigClaude   *ProtocolConfigClaude   `json:"protocol_config_claude" yaml:"protocol_config_claude" mapstructure:"protocol_config_claude"`
	ProtocolConfigDeepSeek *ProtocolConfigDeepSeek `json:"protocol_config_deep_seek" yaml:"protocol_config_deep_seek" mapstructure:"protocol_config_deep_seek"`
	ProtocolConfigGemini   *ProtocolConfigGemini   `json:"protocol_config_gemini" yaml:"protocol_config_gemini" mapstructure:"protocol_config_gemini"`
	ProtocolConfigOllama   *ProtocolConfigOllama   `json:"protocol_config_ollama" yaml:"protocol_config_ollama" mapstructure:"protocol_config_ollama"`
	ProtocolConfigQwen     *ProtocolConfigQwen     `json:"protocol_config_qwen" yaml:"protocol_config_qwen" mapstructure:"protocol_config_qwen"`
}

type ProtocolConfigOpenAI struct {
	ByAzure                  bool   `json:"by_azure" yaml:"by_azure" mapstructure:"by_azure"`
	ApiVersion               string `json:"api_version" yaml:"api_version" mapstructure:"api_version"`
	ResponseFormatType       string `json:"response_format_type" yaml:"response_format_type" mapstructure:"response_format_type"`
	ResponseFormatJsonSchema string `json:"response_format_json_schema" yaml:"response_format_json_schema" mapstructure:"response_format_json_schema"`
}

type ProtocolConfigClaude struct {
	ByBedrock       bool   `json:"by_bedrock" yaml:"by_bedrock" mapstructure:"by_bedrock"`
	AccessKey       string `json:"access_key" yaml:"access_key" mapstructure:"access_key"`
	SecretAccessKey string `json:"secret_access_key" yaml:"secret_access_key" mapstructure:"secret_access_key"`
	SessionToken    string `json:"session_token" yaml:"session_token" mapstructure:"session_token"`
	Region          string `json:"region" yaml:"region" mapstructure:"region"`
	BudgetTokens    *int   `json:"budget_tokens,omitempty" yaml:"budget_tokens" mapstructure:"budget_tokens"`
}

type ProtocolConfigDeepSeek struct {
	ResponseFormatType string `json:"response_format_type" yaml:"response_format_type" mapstructure:"response_format_type"`
}

type ProtocolConfigGemini struct {
	Backend    genai.Backend       `json:"backend,omitempty" yaml:"backend"`
	Project    string              `json:"project,omitempty" yaml:"project"`
	Location   string              `json:"location,omitempty" yaml:"location"`
	APIVersion string              `json:"api_version,omitempty" yaml:"api_version"`
	Headers    map[string][]string `json:"headers,omitempty" yaml:"headers"`

	IncludeThoughts *bool  `json:"include_thoughts,omitempty" yaml:"include_thoughts"` // default true
	ThinkingBudget  *int32 `json:"thinking_budget,omitempty" yaml:"thinking_budget"`   // default nil

	ResponseSchema      *string                             `json:"response_schema" yaml:"response_schema" mapstructure:"response_schema"`
	EnableCodeExecution bool                                `json:"enable_code_execution" yaml:"enable_code_execution" mapstructure:"enable_code_execution"`
	SafetySettings      []ProtocolConfigGeminiSafetySetting `json:"safety_settings" yaml:"safety_settings" mapstructure:"safety_settings"`
}

type ProtocolConfigGeminiSafetySetting struct {
	// Required. The category for this setting.
	Category int32 `json:"category" yaml:"category" mapstructure:"category"`
	// Required. Controls the probability threshold at which harm is blocked.
	Threshold int32 `json:"threshold" yaml:"threshold" mapstructure:"threshold"`
}

type ProtocolConfigOllama struct {
	Format      *string `json:"format" yaml:"format" mapstructure:"format"`
	KeepAliveMs *int64  `json:"keep_alive_ms" yaml:"keep_alive_ms" mapstructure:"keep_alive_ms"`
}

type ProtocolConfigQwen struct {
	ResponseFormatType       *string `json:"response_format_type" yaml:"response_format_type" mapstructure:"response_format_type"`
	ResponseFormatJsonSchema *string `json:"response_format_json_schema" yaml:"response_format_json_schema" mapstructure:"response_format_json_schema"`
}

type ScenarioConfig struct {
	Scenario    Scenario `json:"scenario" yaml:"scenario" mapstructure:"scenario"`
	Quota       *Quota   `json:"quota" yaml:"quota" mapstructure:"quota"`
	Unavailable bool     `json:"unavailable" yaml:"unavailable" mapstructure:"unavailable"`
}

type Quota struct {
	Qpm int64 `json:"qpm" yaml:"qpm" mapstructure:"qpm"`
	Tpm int64 `json:"tpm" yaml:"tpm" mapstructure:"tpm"`
}

type ParamConfig struct {
	ParamSchemas []*ParamSchema `json:"param_schemas" yaml:"param_schemas" mapstructure:"param_schemas"`
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

func (p *ParamConfig) GetCommonParamDefaultVal() CommonParam {
	rawDf := p.GetDefaultVal([]string{"max_tokens", "temperature", "top_p", "top_k", "frequency_penalty", "presence_penalty", "enable_thinking", "stop"})
	cp := CommonParam{}
	if rawDf == nil {
		return cp
	}
	if rawDf["max_tokens"] != "" {
		maxTokens, _ := strconv.ParseInt(rawDf["max_tokens"], 10, 32)
		cp.MaxTokens = ptr.Of(int(maxTokens))
	}
	if rawDf["temperature"] != "" {
		temperature, _ := strconv.ParseFloat(rawDf["temperature"], 32)
		cp.Temperature = ptr.Of(float32(temperature))
	}
	if rawDf["top_p"] != "" {
		topP, _ := strconv.ParseFloat(rawDf["top_p"], 32)
		cp.TopP = ptr.Of(float32(topP))
	}
	if rawDf["top_k"] != "" {
		topK, _ := strconv.ParseInt(rawDf["top_k"], 10, 32)
		cp.TopK = ptr.Of(int(topK))
	}
	if rawDf["stop"] != "" {
		var stop []string
		_ = sonic.UnmarshalString(rawDf["stop"], &stop)
		cp.Stop = stop
	}
	if rawDf["frequency_penalty"] != "" {
		frequencyPenalty, _ := strconv.ParseFloat(rawDf["frequency_penalty"], 32)
		cp.FrequencyPenalty = ptr.Of(float32(frequencyPenalty))
	}
	if rawDf["presence_penalty"] != "" {
		presencePenalty, _ := strconv.ParseFloat(rawDf["presence_penalty"], 32)
		cp.PresencePenalty = ptr.Of(float32(presencePenalty))
	}
	if rawDf["enable_thinking"] != "" {
		enableThinking, _ := strconv.ParseBool(rawDf["enable_thinking"])
		cp.EnableThinking = ptr.Of(enableThinking)
	}
	return cp
}

func (p *ParamConfig) GetDefaultVal(params []string) map[string]string {
	if p == nil || len(p.ParamSchemas) == 0 {
		return nil
	}
	res := make(map[string]string)
	for _, param := range params {
		for _, ps := range p.ParamSchemas {
			if param == ps.Name {
				res[param] = ps.DefaultValue
			}
		}
	}
	return res
}

type ParamSchema struct {
	Name         string         `json:"name" yaml:"name" mapstructure:"name"`
	Label        string         `json:"label" yaml:"label" mapstructure:"label"`
	Desc         string         `json:"desc" yaml:"desc" mapstructure:"desc"`
	Type         ParamType      `json:"type" yaml:"type" mapstructure:"type"`
	Min          string         `json:"min" yaml:"min" mapstructure:"min"`
	Max          string         `json:"max" yaml:"max" mapstructure:"max"`
	DefaultValue string         `json:"default_value" yaml:"default_value" mapstructure:"default_value"`
	Options      []*ParamOption `json:"options" yaml:"options" mapstructure:"options"`
}

type ParamOption struct {
	Value string `json:"value" yaml:"value" mapstructure:"value"`
	Label string `json:"label" yaml:"label" mapstructure:"label"`
}

type ParamType string

const (
	ParamTypeFloat   ParamType = "float"
	ParamTypeInt     ParamType = "int"
	ParamTypeBoolean ParamType = "boolean"
	ParamTypeString  ParamType = "string"
)

type Protocol string

const (
	ProtocolOpenAI   Protocol = "openai"
	ProtocolDeepseek Protocol = "deepseek"
	ProtocolClaude   Protocol = "claude"
	ProtocolOllama   Protocol = "ollama"
	ProtocolGemini   Protocol = "gemini"
	ProtocolQwen     Protocol = "qwen"
)

type ListModelReq struct {
	Scenario  *Scenario
	PageToken int64
	PageSize  int64
}

type GetModelReq struct {
	ModelID int64
}
