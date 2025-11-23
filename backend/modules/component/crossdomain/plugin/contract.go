package plugin

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/model"
)

//go:generate  mockgen -destination pluginmock/plugin_mock.go --package pluginmock -source plugin.go
type PluginService interface {
	BindAgentTools(ctx context.Context, agentID int64, bindTools []*model.BindToolInfo) (err error)
	MGetAgentTools(ctx context.Context, req *model.MGetAgentToolsRequest) (tools []*model.ToolInfo, err error)
	ExecuteTool(ctx context.Context, req *model.ExecuteToolRequest, opts ...model.ExecuteToolOpt) (resp *model.ExecuteToolResponse, err error)
	PublishAPPPlugins(ctx context.Context, req *model.PublishAPPPluginsRequest) (resp *model.PublishAPPPluginsResponse, err error)
	GetAPPAllPlugins(ctx context.Context, appID int64) (plugins []*model.PluginInfo, err error)
	MGetDraftPlugins(ctx context.Context, pluginIDs []int64) (plugins []*model.PluginInfo, err error)
	MGetOnlinePlugins(ctx context.Context, pluginIDs []int64) (plugins []*model.PluginInfo, err error)
	MGetVersionPlugins(ctx context.Context, versionPlugins []model.VersionPlugin) (plugins []*model.PluginInfo, err error)
	MGetDraftTools(ctx context.Context, pluginIDs []int64) (tools []*model.ToolInfo, err error)
	MGetOnlineTools(ctx context.Context, pluginIDs []int64) (tools []*model.ToolInfo, err error)
	MGetVersionTools(ctx context.Context, versionTools []model.VersionTool) (tools []*model.ToolInfo, err error)
}

type InvokableTool interface {
	Info(ctx context.Context) (*schema.ToolInfo, error)
	PluginInvoke(ctx context.Context, argumentsInJSON string) (string, error)
}

var defaultSVC PluginService

func DefaultSVC() PluginService {
	return defaultSVC
}

func SetDefaultSVC(svc PluginService) {
	defaultSVC = svc
}
