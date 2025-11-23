package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/model"
)

func (p *pluginServiceImpl) GetPluginNextVersion(ctx context.Context, pluginID int64) (version string, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) PublishPlugin(ctx context.Context, req *model.PublishPluginRequest) (err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) PublishAPPPlugins(ctx context.Context, req *model.PublishAPPPluginsRequest) (resp *model.PublishAPPPluginsResponse, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) CheckPluginToolsDebugStatus(ctx context.Context, pluginID int64) (err error) {
	panic("implement me")
}
