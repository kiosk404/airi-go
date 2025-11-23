package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/infra/dao"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
)

func (p *pluginServiceImpl) CreateDraftPlugin(ctx context.Context, req *dao.CreateDraftPluginRequest) (pluginID int64, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) GetDraftPlugin(ctx context.Context, pluginID int64) (plugin *entity.PluginInfo, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) MGetDraftPlugins(ctx context.Context, pluginIDs []int64) (plugins []*entity.PluginInfo, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) ListDraftPlugins(ctx context.Context, req *dao.ListDraftPluginsRequest) (resp *dao.ListDraftPluginsResponse, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) CreateDraftPluginWithCode(ctx context.Context, req *dao.CreateDraftPluginWithCodeRequest) (resp *dao.CreateDraftPluginWithCodeResponse, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) UpdateDraftPluginWithCode(ctx context.Context, req *dao.UpdateDraftPluginWithCodeRequest) (err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) UpdateDraftPlugin(ctx context.Context, req *dao.UpdateDraftPluginRequest) (err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) DeleteDraftPlugin(ctx context.Context, pluginID int64) (err error) {
	return p.pluginRepo.DeleteDraftPlugin(ctx, pluginID)
}

func (p *pluginServiceImpl) MGetDraftTools(ctx context.Context, toolIDs []int64) (tools []*entity.ToolInfo, err error) {
	tools, err = p.toolRepo.MGetDraftTools(ctx, toolIDs)
	if err != nil {
		return nil, errorx.Wrapf(err, "MGetDraftTools failed, toolIDs=%v", toolIDs)
	}

	return tools, nil
}

func (p *pluginServiceImpl) UpdateDraftTool(ctx context.Context, req *dao.UpdateDraftToolRequest) (err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) ConvertToOpenapi3Doc(ctx context.Context, req *dao.ConvertToOpenapi3DocRequest) (resp *dao.ConvertToOpenapi3DocResponse) {
	panic("implement me")
}

func (p *pluginServiceImpl) CreateDraftToolsWithCode(ctx context.Context, req *dao.CreateDraftToolsWithCodeRequest) (resp *dao.CreateDraftToolsWithCodeResponse, err error) {
	panic("implement me")
}
