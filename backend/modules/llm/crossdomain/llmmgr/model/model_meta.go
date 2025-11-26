package model

import (
	"fmt"
)

type BasicConfiguration struct {
	PluginConfiguration *PluginConfiguration `json:"plugin_configuration"`
	CodeRunnerType      CodeRunnerType       `json:"code_runner_type"`
	SandboxConfig       *SandboxConfig       `json:"sandbox_config,omitempty"`
	ServerHost          string               `json:"server_host"`
}

type PluginConfiguration struct {
	SaasPluginEnabled bool   `json:"saas_plugin_enabled"`
	APIToken          string `json:"api_token"`
	SaasAPIBaseURL    string `json:"saas_api_base_url"`
}

type CodeRunnerType int64

const (
	CoderunnerTypeLocal   CodeRunnerType = 0
	CoderunnerTypeSandbox CodeRunnerType = 1
)

func (p CodeRunnerType) String() string {
	switch p {
	case CoderunnerTypeLocal:
		return "Local"
	case CoderunnerTypeSandbox:
		return "Sandbox"
	}
	return "<UNSET>"
}

func CodeRunnerTypeFromString(s string) (CodeRunnerType, error) {
	switch s {
	case "Local":
		return CoderunnerTypeLocal, nil
	case "Sandbox":
		return CoderunnerTypeSandbox, nil
	}
	return CodeRunnerType(0), fmt.Errorf("not a valid CodeRunnerType string")
}

type SandboxConfig struct {
	AllowEnv       string  `json:"allow_env"`
	AllowRead      string  `json:"allow_read"`
	AllowWrite     string  `json:"allow_write"`
	AllowRun       string  `json:"allow_run"`
	AllowNet       string  `json:"allow_net"`
	AllowFfi       string  `json:"allow_ffi"`
	NodeModulesDir string  `json:"node_modules_dir"`
	TimeoutSeconds float64 `json:"timeout_seconds"`
	MemoryLimitMb  int64   `json:"memory_limit_mb"`
}

func NewSandboxConfig() *SandboxConfig {
	return &SandboxConfig{}
}

func (p *SandboxConfig) InitDefault() {
}

func (p *SandboxConfig) GetAllowEnv() (v string) {
	return p.AllowEnv
}

func (p *SandboxConfig) GetAllowRead() (v string) {
	return p.AllowRead
}

func (p *SandboxConfig) GetAllowWrite() (v string) {
	return p.AllowWrite
}

func (p *SandboxConfig) GetAllowRun() (v string) {
	return p.AllowRun
}

func (p *SandboxConfig) GetAllowNet() (v string) {
	return p.AllowNet
}

func (p *SandboxConfig) GetAllowFfi() (v string) {
	return p.AllowFfi
}

func (p *SandboxConfig) GetNodeModulesDir() (v string) {
	return p.NodeModulesDir
}

func (p *SandboxConfig) GetTimeoutSeconds() (v float64) {
	return p.TimeoutSeconds
}

func (p *SandboxConfig) GetMemoryLimitMb() (v int64) {
	return p.MemoryLimitMb
}
func (p *SandboxConfig) SetAllowEnv(val string) {
	p.AllowEnv = val
}
func (p *SandboxConfig) SetAllowRead(val string) {
	p.AllowRead = val
}
func (p *SandboxConfig) SetAllowWrite(val string) {
	p.AllowWrite = val
}
func (p *SandboxConfig) SetAllowRun(val string) {
	p.AllowRun = val
}
func (p *SandboxConfig) SetAllowNet(val string) {
	p.AllowNet = val
}
func (p *SandboxConfig) SetAllowFfi(val string) {
	p.AllowFfi = val
}
func (p *SandboxConfig) SetNodeModulesDir(val string) {
	p.NodeModulesDir = val
}
func (p *SandboxConfig) SetTimeoutSeconds(val float64) {
	p.TimeoutSeconds = val
}
func (p *SandboxConfig) SetMemoryLimitMb(val int64) {
	p.MemoryLimitMb = val
}

func (p *SandboxConfig) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SandboxConfig(%+v)", *p)
}

type ModelStatus int64

const (
	ModelStatus_StatusDefault ModelStatus = 0
	ModelStatus_StatusInUse   ModelStatus = 1
	ModelStatus_StatusDeleted ModelStatus = 2
)

func (p ModelStatus) String() string {
	switch p {
	case ModelStatus_StatusDefault:
		return "StatusDefault"
	case ModelStatus_StatusInUse:
		return "StatusInUse"
	case ModelStatus_StatusDeleted:
		return "StatusDeleted"
	}
	return "<UNSET>"
}

func ModelStatusFromString(s string) (ModelStatus, error) {
	switch s {
	case "StatusDefault":
		return ModelStatus_StatusDefault, nil
	case "StatusInUse":
		return ModelStatus_StatusInUse, nil
	case "StatusDeleted":
		return ModelStatus_StatusDeleted, nil
	}
	return ModelStatus(0), fmt.Errorf("not a valid ModelStatus string")
}

type EmbeddingType int64

const (
	EmbeddingType_Ark    EmbeddingType = 0
	EmbeddingType_OpenAI EmbeddingType = 1
	EmbeddingType_Ollama EmbeddingType = 2
	EmbeddingType_Gemini EmbeddingType = 3
	EmbeddingType_HTTP   EmbeddingType = 4
)

func (p EmbeddingType) String() string {
	switch p {
	case EmbeddingType_Ark:
		return "Ark"
	case EmbeddingType_OpenAI:
		return "OpenAI"
	case EmbeddingType_Ollama:
		return "Ollama"
	case EmbeddingType_Gemini:
		return "Gemini"
	case EmbeddingType_HTTP:
		return "HTTP"
	}
	return "<UNSET>"
}

func EmbeddingTypeFromString(s string) (EmbeddingType, error) {
	switch s {
	case "Ark":
		return EmbeddingType_Ark, nil
	case "OpenAI":
		return EmbeddingType_OpenAI, nil
	case "Ollama":
		return EmbeddingType_Ollama, nil
	case "Gemini":
		return EmbeddingType_Gemini, nil
	case "HTTP":
		return EmbeddingType_HTTP, nil
	}
	return EmbeddingType(0), fmt.Errorf("not a valid EmbeddingType string")
}

type ModelConfig struct {
}

type Config struct {
	//base      *base.BaseConfig
	//knowledge *knowledge.KnowledgeConfig
	model *ModelConfig
}

type ModelProvider struct {
	Name        *I18nText  `thrift:"name,1" form:"name" json:"name" query:"name"`
	IconURI     string     `thrift:"icon_uri,2" form:"icon_uri" json:"icon_uri" query:"icon_uri"`
	IconURL     string     `thrift:"icon_url,3" form:"icon_url" json:"icon_url" query:"icon_url"`
	Description *I18nText  `thrift:"description,4" form:"description" json:"description" query:"description"`
	ModelClass  ModelClass `thrift:"model_class,5" form:"model_class" json:"model_class" query:"model_class"`
}

type ModelClass int64

const (
	ModelClass_GPT      ModelClass = 1
	ModelClass_QWen     ModelClass = 2
	ModelClass_Gemini   ModelClass = 3
	ModelClass_DeepSeek ModelClass = 4
	ModelClass_Ollama   ModelClass = 5
	ModelClass_Other    ModelClass = 999
)

func (p ModelClass) String() string {
	switch p {
	case ModelClass_GPT:
		return "GPT"
	case ModelClass_QWen:
		return "QWen"
	case ModelClass_Gemini:
		return "Gemini"
	case ModelClass_DeepSeek:
		return "DeepSeek"
	case ModelClass_Ollama:
		return "Ollama"
	case ModelClass_Other:
		return "Other"
	}
	return "<UNSET>"
}

func ModelClassFromString(s string) (ModelClass, error) {
	switch s {
	case "GPT":
		return ModelClass_GPT, nil
	case "QWen":
		return ModelClass_QWen, nil
	case "Gemini":
		return ModelClass_Gemini, nil
	case "DeepSeek":
		return ModelClass_DeepSeek, nil
	case "Ollama":
		return ModelClass_Ollama, nil
	case "Other":
		return ModelClass_Other, nil
	}
	return ModelClass(0), fmt.Errorf("not a valid ModelClass string")
}
