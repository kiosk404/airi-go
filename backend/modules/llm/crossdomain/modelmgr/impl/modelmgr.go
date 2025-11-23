package impl

import (
	"context"

	modelmgrapp "github.com/kiosk404/airi-go/backend/modules/llm/application"
	crossmodelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llmimpl/chatmodel"
)

var defaultSVC crossmodelmgr.Model

type impl struct {
	DomainSVC modelmgrapp.ModelServiceComponents
	factory   chatmodel.Factory
}

func InitDomainService(c modelmgrapp.ModelServiceComponents, f chatmodel.Factory) crossmodelmgr.Model {
	defaultSVC = &impl{
		DomainSVC: c,
		factory:   f,
	}

	return defaultSVC
}

func (c *impl) ListModels(ctx context.Context, req *model.ListModelsRequest) (r *model.ListModelsResponse, err error) {
	resp, err := c.DomainSVC.ListModels(ctx, convertListModelsRequest(req))
	return convertListModelsResponse(resp), err
}

func (c *impl) GetModel(ctx context.Context, req *model.GetModelRequest) (r *model.GetModelResponse, err error) {
	resp, err := c.DomainSVC.GetModel(ctx, convertGetModelRequest(req))
	return convertGetModelResponse(resp), err
}

func (c *impl) Chat(ctx context.Context, req *model.ChatRequest) (r *model.ChatResponse, err error) {
	resp, err := c.DomainSVC.Chat(ctx, convertChatRequest(req))
	return convertChatResponse(resp), err
}

func (c *impl) ChatStream(req *model.ChatRequest, stream model.ChatResponseStream) (err error) {
	return nil
}
