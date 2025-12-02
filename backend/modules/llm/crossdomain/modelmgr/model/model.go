package model

import (
	"fmt"
)

type Model struct {
	ID              int64            `json:"id" query:"id"`
	Provider        *ModelProvider   `json:"provider" query:"provider"`
	DisplayInfo     *DisplayInfo     `json:"display_info" query:"display_info"`
	Capability      *ModelAbility    `json:"capability" query:"capability"`
	Connection      *Connection      `json:"connection" query:"connection"`
	Type            ModelType        `json:"type" query:"type"`
	Parameters      []ModelParameter `json:"parameters" query:"parameters"`
	Status          ModelStatus      `json:"status" query:"status"`
	EnableBase64URL bool             `json:"enable_base64_url" query:"enable_base64_url"`
	DeleteAtMs      int64            `json:"delete_at_ms" query:"delete_at_ms"`
}

type ModelType int64

const (
	ModelType_LLM           ModelType = 0
	ModelType_TextEmbedding ModelType = 1
	ModelType_Rerank        ModelType = 2
)

func (p ModelType) String() string {
	switch p {
	case ModelType_LLM:
		return "LLM"
	case ModelType_TextEmbedding:
		return "TextEmbedding"
	case ModelType_Rerank:
		return "Rerank"
	}
	return "<UNSET>"
}

func (p ModelType) Int32() int32 {
	return int32(p)
}

func ModelTypeFromString(s string) (ModelType, error) {
	switch s {
	case "LLM":
		return ModelType_LLM, nil
	case "TextEmbedding":
		return ModelType_TextEmbedding, nil
	case "Rerank":
		return ModelType_Rerank, nil
	}
	return ModelType(0), fmt.Errorf("not a valid ModelType string")
}

type CreateModelRequest struct {
	ModelClass    ModelClass
	ModelShowName string
	Conn          Connection
	Extra         ModelExtra
}

type UpdateModelRequest struct {
	ID    int64
	Conn  Connection
	Extra ModelExtra
}
