package modelmgr

import (
	"context"

	model "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
)

type Model interface {
	LLMManageService
	LLMRuntimeService
}

type LLMManageService interface {
	ListModels(ctx context.Context, req *model.ListModelsRequest) (r *model.ListModelsResponse, err error)
	GetModel(ctx context.Context, req *model.GetModelRequest) (r *model.GetModelResponse, err error)
}

type LLMRuntimeService interface {
	Chat(ctx context.Context, req *model.ChatRequest) (r *model.ChatResponse, err error)
	ChatStream(req *model.ChatRequest, stream model.ChatResponseStream) (err error)
}

type ModelAppMeta = model.Model

var defaultSVC Model

func DefaultSVC() Model {
	return defaultSVC
}

func SetDefaultSVC(svc Model) {
	defaultSVC = svc
}
