package model

import (
	"fmt"
)

type ModelClass int64

const (
	ModelClass_GPT      ModelClass = 1
	ModelClass_QWen     ModelClass = 2
	ModelClass_Gemini   ModelClass = 3
	ModelClass_DeepSeek ModelClass = 4
	ModelClass_Ollama   ModelClass = 5
	ModelClass_Claude   ModelClass = 6
	ModelClass_Other    ModelClass = 999
)

func (p ModelClass) String() string {
	switch p {
	case ModelClass_GPT:
		return "gpt"
	case ModelClass_QWen:
		return "qwen"
	case ModelClass_Gemini:
		return "gemini"
	case ModelClass_DeepSeek:
		return "deepseek"
	case ModelClass_Ollama:
		return "ollama"
	case ModelClass_Claude:
		return "claude"
	case ModelClass_Other:
		return "other"
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
