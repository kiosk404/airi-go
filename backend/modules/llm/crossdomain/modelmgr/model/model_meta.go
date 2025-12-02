package model

import (
	"fmt"
)

type ModelMeta struct {
	DisplayInfo     *DisplayInfo      `json:"display_info,omitempty"`
	Capability      *ModelAbility     `json:"capability,omitempty"`
	Connection      *Connection       `json:"connection,omitempty"`
	Parameters      []*ModelParameter `json:"parameters,omitempty"`
	EnableBase64URL bool              `json:"enable_base64_url,omitempty"`
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

type ModelExtra struct {
	EnableBase64URL bool `json:"enable_base64_url"`
}
