package model

import (
	"fmt"
)

type Connection struct {
	BaseConnInfo *BaseConnectionInfo `json:"base_conn_info" query:"base_conn_info"`
	Ark          *ArkConnInfo        `json:"ark,omitempty" query:"ark"`
	Openai       *OpenAIConnInfo     `json:"openai,omitempty" query:"openai"`
	Deepseek     *DeepseekConnInfo   `json:"deepseek,omitempty" query:"deepseek"`
	Gemini       *GeminiConnInfo     `json:"gemini,omitempty" query:"gemini"`
	Qwen         *QwenConnInfo       `json:"qwen,omitempty" query:"qwen"`
	Ollama       *OllamaConnInfo     `json:"ollama,omitempty" query:"ollama"`
	Claude       *ClaudeConnInfo     `json:"claude,omitempty" query:"claude"`
}

type BaseConnectionInfo struct {
	BaseURL      string       `json:"base_url" query:"base_url"`
	APIKey       string       `json:"api_key" query:"api_key"`
	Model        string       `json:"model" query:"model"`
	ThinkingType ThinkingType `json:"thinking_type" query:"thinking_type"`
}

func NewBaseConnectionInfo() *BaseConnectionInfo {
	return &BaseConnectionInfo{}
}

func (p *BaseConnectionInfo) InitDefault() {
}

func (p *BaseConnectionInfo) GetBaseURL() (v string) {
	return p.BaseURL
}

func (p *BaseConnectionInfo) GetAPIKey() (v string) {
	return p.APIKey
}

func (p *BaseConnectionInfo) GetModel() (v string) {
	return p.Model
}

func (p *BaseConnectionInfo) GetThinkingType() (v ThinkingType) {
	return p.ThinkingType
}

type ArkConnInfo struct {
	Region  string `json:"region" query:"region"`
	APIType string `json:"api_type" query:"api_type"`
}

func NewArkConnInfo() *ArkConnInfo {
	return &ArkConnInfo{}
}

func (p *ArkConnInfo) InitDefault() {
}

func (p *ArkConnInfo) GetRegion() (v string) {
	return p.Region
}

func (p *ArkConnInfo) GetAPIType() (v string) {
	return p.APIType
}

var fieldIDToName_ArkConnInfo = map[int16]string{
	1: "region",
	3: "api_type",
}

type OpenAIConnInfo struct {
	ByAzure    bool   `json:"by_azure" query:"by_azure"`
	APIVersion string `json:"api_version" query:"api_version"`
}

func NewOpenAIConnInfo() *OpenAIConnInfo {
	return &OpenAIConnInfo{}
}

func (p *OpenAIConnInfo) InitDefault() {
}

func (p *OpenAIConnInfo) GetByAzure() (v bool) {
	return p.ByAzure
}

func (p *OpenAIConnInfo) GetAPIVersion() (v string) {
	return p.APIVersion
}

type DeepseekConnInfo struct {
}

func NewDeepseekConnInfo() *DeepseekConnInfo {
	return &DeepseekConnInfo{}
}

func (p *DeepseekConnInfo) InitDefault() {
}

type GeminiConnInfo struct {
	// "1" for BackendGeminiAPI / "2" for BackendVertexAI
	Backend  int32  `json:"backend" query:"backend"`
	Project  string `json:"project" query:"project"`
	Location string `json:"location" query:"location"`
}

func NewGeminiConnInfo() *GeminiConnInfo {
	return &GeminiConnInfo{}
}

func (p *GeminiConnInfo) InitDefault() {
}

func (p *GeminiConnInfo) GetBackend() (v int32) {
	return p.Backend
}

func (p *GeminiConnInfo) GetProject() (v string) {
	return p.Project
}

func (p *GeminiConnInfo) GetLocation() (v string) {
	return p.Location
}

type QwenConnInfo struct {
}

func NewQwenConnInfo() *QwenConnInfo {
	return &QwenConnInfo{}
}

func (p *QwenConnInfo) InitDefault() {
}

type OllamaConnInfo struct {
}

func NewOllamaConnInfo() *OllamaConnInfo {
	return &OllamaConnInfo{}
}

func (p *OllamaConnInfo) InitDefault() {
}

type ClaudeConnInfo struct {
}

func NewClaudeConnInfo() *ClaudeConnInfo {
	return &ClaudeConnInfo{}
}

func (p *ClaudeConnInfo) InitDefault() {
}

type ThinkingType int64

const (
	ThinkingType_Default ThinkingType = 0
	ThinkingType_Enable  ThinkingType = 1
	ThinkingType_Disable ThinkingType = 2
	ThinkingType_Auto    ThinkingType = 3
)

func (p ThinkingType) String() string {
	switch p {
	case ThinkingType_Default:
		return "Default"
	case ThinkingType_Enable:
		return "Enable"
	case ThinkingType_Disable:
		return "Disable"
	case ThinkingType_Auto:
		return "Auto"
	}
	return "<UNSET>"
}

func ThinkingTypeFromString(s string) (ThinkingType, error) {
	switch s {
	case "Default":
		return ThinkingType_Default, nil
	case "Enable":
		return ThinkingType_Enable, nil
	case "Disable":
		return ThinkingType_Disable, nil
	case "Auto":
		return ThinkingType_Auto, nil
	}
	return ThinkingType(0), fmt.Errorf("not a valid ThinkingType string")
}
