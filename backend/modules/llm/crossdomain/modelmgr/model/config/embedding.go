package config

import (
	"fmt"
)

type EmbeddingConnection struct {
	BaseConnInfo  *BaseConnectionInfo `json:"base_conn_info" query:"base_conn_info"`
	EmbeddingInfo *EmbeddingInfo      `json:"embedding_info" query:"embedding_info"`
	Ark           *ArkConnInfo        `json:"ark,omitempty" query:"ark"`
	Openai        *OpenAIConnInfo     `json:"openai,omitempty" query:"openai"`
	Ollama        *OllamaConnInfo     `json:"ollama,omitempty" query:"ollama"`
	Gemini        *GeminiConnInfo     `json:"gemini,omitempty" query:"gemini"`
	HTTP          *HttpConnection     `json:"http,omitempty" query:"http"`
}

type BaseConnectionInfo struct {
	BaseURL      string       `json:"base_url" query:"base_url"`
	APIKey       string       `json:"api_key" query:"api_key"`
	Model        string       `json:"model" query:"model"`
	ThinkingType ThinkingType `json:"thinking_type" query:"thinking_type"`
}

type EmbeddingInfo struct {
	Dims int32 `json:"dims" query:"dims"`
}

type ArkConnInfo struct {
	Region  string `json:"region" query:"region"`
	APIType string `json:"api_type" query:"api_type"`
}

type OpenAIConnInfo struct {
	ByAzure    bool   `json:"by_azure" query:"by_azure"`
	APIVersion string `json:"api_version" query:"api_version"`
}

type OllamaConnInfo struct {
}

type GeminiConnInfo struct {
	// "1" for BackendGeminiAPI / "2" for BackendVertexAI
	Backend  int32  `json:"backend" query:"backend"`
	Project  string `json:"project" query:"project"`
	Location string `json:"location" query:"location"`
}

type HttpConnection struct {
	Address string `json:"address" query:"address"`
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
