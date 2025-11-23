package application

import (
	"context"
	"strconv"

	"github.com/kiosk404/airi-go/backend/api/model/llm/manage"
	"github.com/kiosk404/airi-go/backend/modules/llm/application/convertor"
	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/component/rpc"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

type manageApp struct {
	manageSrv service.IManage
	auth      rpc.IAuthProvider
}

func NewManageApplication(
	manageSrv service.IManage,
	auth rpc.IAuthProvider,
) manage.LLMManageService {
	return &manageApp{
		manageSrv: manageSrv,
		auth:      auth,
	}
}

func (m *manageApp) ListModels(ctx context.Context, req *manage.ListModelsRequest) (r *manage.ListModelsResponse, err error) {
	r = manage.NewListModelsResponse()
	if err := m.auth.CheckPermission(ctx, "listModels"); err != nil {
		return r, err
	}
	var scenario *model.Scenario
	if req.Scenario != nil {
		scenario = ptr.Of(model.Scenario(req.GetScenario()))
	} else {
		scenario = ptr.Of(model.ScenarioDefault)
	}
	pageToken, _ := strconv.ParseInt(req.GetPageToken(), 10, 64)
	models, total, hasMore, nextPageToken, err := m.manageSrv.ListModels(ctx, entity.ListModelsRequest{
		Scenario:  scenario,
		PageSize:  int64(req.GetPageSize()),
		PageToken: pageToken,
	})
	if err != nil {
		return r, err
	}
	r.SetModels(convertor.ModelsDO2DTO(models, true))
	r.SetTotal(ptr.Of(int32(total)))
	r.SetNextPageToken(ptr.Of(strconv.FormatInt(nextPageToken, 10)))
	r.SetHasMore(ptr.Of(hasMore))
	return r, nil
}

func (m *manageApp) GetModel(ctx context.Context, req *manage.GetModelRequest) (r *manage.GetModelResponse, err error) {
	r = manage.NewGetModelResponse()
	if err := m.auth.CheckPermission(ctx, "getModel"); err != nil {
		return r, err
	}
	model, err := m.manageSrv.GetModelByID(ctx, req.GetModelID())
	if err != nil {
		return r, err
	}
	r.SetModel(convertor.ModelDO2DTO(model, true))
	return r, nil
}
