package impl

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	crossplugin "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin"
	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/model"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/entity"
	plugin "github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/service"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
)

var defaultSVC crossplugin.PluginService

type impl struct {
	DomainSVC plugin.PluginService
	tos       storage.Storage
}

func InitDomainService(c plugin.PluginService, tos storage.Storage) crossplugin.PluginService {
	defaultSVC = &impl{
		DomainSVC: c,
		tos:       tos,
	}

	return defaultSVC
}

func (s *impl) BindAgentTools(ctx context.Context, agentID int64, bindTools []*model.BindToolInfo) (err error) {
	return s.DomainSVC.BindAgentTools(ctx, agentID, bindTools)
}

func (s *impl) MGetAgentTools(ctx context.Context, req *model.MGetAgentToolsRequest) (tools []*model.ToolInfo, err error) {
	return []*model.ToolInfo{}, nil
	//return s.DomainSVC.MGetAgentTools(ctx, req)
}

func (s *impl) ExecuteTool(ctx context.Context, req *model.ExecuteToolRequest, opts ...model.ExecuteToolOpt) (resp *model.ExecuteToolResponse, err error) {
	return s.DomainSVC.ExecuteTool(ctx, req, opts...)
}

func (s *impl) PublishAPPPlugins(ctx context.Context, req *model.PublishAPPPluginsRequest) (resp *model.PublishAPPPluginsResponse, err error) {
	return s.DomainSVC.PublishAPPPlugins(ctx, req)
}

func (s *impl) GetAPPAllPlugins(ctx context.Context, appID int64) (plugins []*model.PluginInfo, err error) {
	_plugins, err := s.DomainSVC.GetAPPAllPlugins(ctx, appID)
	if err != nil {
		return nil, err
	}

	plugins = slices.Transform(_plugins, func(e *entity.PluginInfo) *model.PluginInfo {
		return e.PluginInfo
	})

	return plugins, nil
}

func (s *impl) MGetDraftPlugins(ctx context.Context, pluginIDs []int64) (plugins []*model.PluginInfo, err error) {
	ePlugins, err := s.DomainSVC.MGetDraftPlugins(ctx, pluginIDs)
	if err != nil {
		return nil, err
	}

	plugins = slices.Transform(ePlugins, func(e *entity.PluginInfo) *model.PluginInfo {
		return e.PluginInfo
	})

	return plugins, nil
}

func (s *impl) MGetOnlinePlugins(ctx context.Context, pluginIDs []int64) (plugins []*model.PluginInfo, err error) {
	ePlugins, err := s.DomainSVC.MGetOnlinePlugins(ctx, pluginIDs)
	if err != nil {
		return nil, err
	}

	plugins = slices.Transform(ePlugins, func(e *entity.PluginInfo) *model.PluginInfo {
		return e.PluginInfo
	})

	return plugins, nil
}

func (s *impl) MGetVersionPlugins(ctx context.Context, versionPlugins []model.VersionPlugin) (plugins []*model.PluginInfo, err error) {
	ePlugins, err := s.DomainSVC.MGetVersionPlugins(ctx, versionPlugins)
	if err != nil {
		return nil, err
	}

	plugins = slices.Transform(ePlugins, func(e *entity.PluginInfo) *model.PluginInfo {
		return e.PluginInfo
	})

	return plugins, nil
}

func (s *impl) MGetDraftTools(ctx context.Context, pluginIDs []int64) (tools []*model.ToolInfo, err error) {
	return s.DomainSVC.MGetDraftTools(ctx, pluginIDs)
}

func (s *impl) MGetOnlineTools(ctx context.Context, pluginIDs []int64) (tools []*model.ToolInfo, err error) {
	return s.DomainSVC.MGetOnlineTools(ctx, pluginIDs)
}

func (s *impl) MGetVersionTools(ctx context.Context, versionTools []model.VersionTool) (tools []*model.ToolInfo, err error) {
	return s.DomainSVC.MGetVersionTools(ctx, versionTools)
}
