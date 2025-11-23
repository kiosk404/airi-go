package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/model"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/infra/dao"
)

func (p *pluginServiceImpl) BindAgentTools(ctx context.Context, agentID int64, bindTools []*model.BindToolInfo) (err error) {
	return p.toolRepo.BindDraftAgentTools(ctx, agentID, bindTools)
}

func (p *pluginServiceImpl) DuplicateDraftAgentTools(ctx context.Context, fromAgentID, toAgentID int64) (err error) {
	return p.toolRepo.DuplicateDraftAgentTools(ctx, fromAgentID, toAgentID)
}

func (p *pluginServiceImpl) GetDraftAgentToolByName(ctx context.Context, agentID int64, pluginID int64, toolName string) (tool *entity.ToolInfo, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) MGetAgentTools(ctx context.Context, req *model.MGetAgentToolsRequest) (tools []*entity.ToolInfo, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) PublishAgentTools(ctx context.Context, agentID int64, agentVersion string) (err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) UpdateBotDefaultParams(ctx context.Context, req *dao.UpdateBotDefaultParamsRequest) (err error) {
	panic("implement me")
}
