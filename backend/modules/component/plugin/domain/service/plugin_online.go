package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/model"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/infra/dao"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
)

func (p *pluginServiceImpl) GetOnlinePlugin(ctx context.Context, pluginID int64) (plugin *entity.PluginInfo, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) MGetOnlinePlugins(ctx context.Context, pluginIDs []int64) (plugins []*entity.PluginInfo, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) GetOnlineTool(ctx context.Context, toolID int64) (tool *entity.ToolInfo, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) MGetOnlineTools(ctx context.Context, toolIDs []int64) (tools []*entity.ToolInfo, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) MGetVersionTools(ctx context.Context, versionTools []model.VersionTool) (tools []*entity.ToolInfo, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) ListPluginProducts(ctx context.Context, req *dao.ListPluginProductsRequest) (resp *dao.ListPluginProductsResponse, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) GetPluginProductAllTools(ctx context.Context, pluginID int64) (tools []*entity.ToolInfo, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) DeleteAPPAllPlugins(ctx context.Context, appID int64) (pluginIDs []int64, err error) {
	return p.pluginRepo.DeleteAPPAllPlugins(ctx, appID)
}

func (p *pluginServiceImpl) GetAPPAllPlugins(ctx context.Context, appID int64) (plugins []*entity.PluginInfo, err error) {
	plugins, err = p.pluginRepo.GetAPPAllDraftPlugins(ctx, appID)
	if err != nil {
		return nil, errorx.Wrapf(err, "GetAPPAllDraftPlugins failed, appID=%d", appID)
	}

	return plugins, nil
}

func (p *pluginServiceImpl) MGetVersionPlugins(ctx context.Context, versionPlugins []model.VersionPlugin) (plugins []*entity.PluginInfo, err error) {
	plugins, err = p.pluginRepo.MGetVersionPlugins(ctx, versionPlugins)
	if err != nil {
		return nil, errorx.Wrapf(err, "MGetVersionPlugins failed, versionPlugins=%v", versionPlugins)
	}

	return plugins, nil
}

func (p *pluginServiceImpl) ListCustomOnlinePlugins(ctx context.Context, pageInfo dao.PageInfo) (plugins []*entity.PluginInfo, total int64, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) MGetPluginLatestVersion(ctx context.Context, pluginIDs []int64) (resp *model.MGetPluginLatestVersionResponse, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) CopyPlugin(ctx context.Context, req *dao.CopyPluginRequest) (resp *dao.CopyPluginResponse, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) MoveAPPPluginToLibrary(ctx context.Context, pluginID int64) (draftPlugin *entity.PluginInfo, err error) {
	panic("implement me")
}
