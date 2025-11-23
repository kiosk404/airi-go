package model

import (
	"github.com/kiosk404/airi-go/backend/api/model/llm/domain/manage"
)

type ParameterName string

const (
	Temperature      ParameterName = "temperature"
	TopP             ParameterName = "top_p"
	TopK             ParameterName = "top_k"
	MaxTokens        ParameterName = "max_tokens"
	RespFormat       ParameterName = "response_format"
	FrequencyPenalty ParameterName = "frequency_penalty"
	PresencePenalty  ParameterName = "presence_penalty"
	EnableThinking   ParameterName = "enable_thinking"
	Stop             ParameterName = "stop"
)

func (p ParameterName) String() string {
	return string(p)
}

type ParamType manage.ParamType

const (
	ParamTypeFloat   ParamType = "float"
	ParamTypeInt     ParamType = "int"
	ParamTypeBoolean ParamType = "boolean"
	ParamTypeString  ParamType = "string"
)

type DefaultType manage.DefaultType

const (
	DefaultTypeDefault  DefaultType = "default_val"
	DefaultTypeCreative DefaultType = "creative"
	DefaultTypeBalance  DefaultType = "balance"
	DefaultTypePrecise  DefaultType = "precise"
)

type Scenario string

const (
	ScenarioDefault   Scenario = "default"
	ScenarioEvaluator Scenario = "evaluator"
)

func ScenarioValue(scenario *Scenario) Scenario {
	if scenario == nil {
		return ScenarioDefault
	}
	return *scenario
}

type Modal string

const (
	ModalText  Modal = "text"
	ModalImage Modal = "image"
	ModalFile  Modal = "file"
	ModalAudio Modal = "audio"
	ModalVideo Modal = "video"
)

type ModelStatus int64

const (
	StatusDefault ModelStatus = 0  // Default state when not configured, equivalent to StatusInUse
	StatusInUse   ModelStatus = 1  // In the application, it can be used to create new
	StatusPending ModelStatus = 5  // To be offline, it can be used and cannot be created.
	StatusDeleted ModelStatus = 10 // It is offline, unusable, and cannot be created.
)
