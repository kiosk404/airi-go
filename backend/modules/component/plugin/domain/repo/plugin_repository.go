package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/model"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/infra/dao"
)

type PluginRepository interface {
	CreateDraftPlugin(ctx context.Context, plugin *entity.PluginInfo) (pluginID int64, err error)
	CreateDraftPluginWithCode(ctx context.Context, req *CreateDraftPluginWithCodeRequest) (resp *CreateDraftPluginWithCodeResponse, err error)
	GetDraftPlugin(ctx context.Context, pluginID int64, opts ...PluginSelectedOptions) (plugin *entity.PluginInfo, exist bool, err error)
	MGetDraftPlugins(ctx context.Context, pluginIDs []int64, opts ...PluginSelectedOptions) (plugins []*entity.PluginInfo, err error)
	GetAPPAllDraftPlugins(ctx context.Context, appID int64, opts ...PluginSelectedOptions) (plugins []*entity.PluginInfo, err error)
	ListDraftPlugins(ctx context.Context, req *ListDraftPluginsRequest) (resp *ListDraftPluginsResponse, err error)
	UpdateDraftPlugin(ctx context.Context, plugin *entity.PluginInfo) (err error)
	UpdateDraftPluginWithoutURLChanged(ctx context.Context, plugin *entity.PluginInfo) (err error)
	UpdateDraftPluginWithCode(ctx context.Context, req *UpdatePluginDraftWithCode) (err error)
	DeleteDraftPlugin(ctx context.Context, pluginID int64) (err error)
	DeleteAPPAllPlugins(ctx context.Context, appID int64) (pluginIDs []int64, err error)
	UpdateDebugExample(ctx context.Context, pluginID int64, openapiDoc *model.Openapi3T) (err error)

	GetOnlinePlugin(ctx context.Context, pluginID int64, opts ...PluginSelectedOptions) (plugin *entity.PluginInfo, exist bool, err error)
	MGetOnlinePlugins(ctx context.Context, pluginIDs []int64, opts ...PluginSelectedOptions) (plugins []*entity.PluginInfo, err error)
	ListCustomOnlinePlugins(ctx context.Context, spaceID int64, pageInfo dao.PageInfo) (plugins []*entity.PluginInfo, total int64, err error)

	GetVersionPlugin(ctx context.Context, vPlugin model.VersionPlugin) (plugin *entity.PluginInfo, exist bool, err error)
	MGetVersionPlugins(ctx context.Context, vPlugins []model.VersionPlugin, opts ...PluginSelectedOptions) (plugin []*entity.PluginInfo, err error)

	PublishPlugin(ctx context.Context, draftPlugin *entity.PluginInfo) (err error)
	PublishPlugins(ctx context.Context, draftPlugins []*entity.PluginInfo) (err error)

	CopyPlugin(ctx context.Context, req *CopyPluginRequest) (plugin *entity.PluginInfo, tools []*entity.ToolInfo, err error)
	MoveAPPPluginToLibrary(ctx context.Context, draftPlugin *entity.PluginInfo, draftTools []*entity.ToolInfo) (err error)
}

type UpdatePluginDraftWithCode struct {
	PluginID   int64
	OpenapiDoc *model.Openapi3T
	Manifest   *model.PluginManifest

	UpdatedTools  []*entity.ToolInfo
	NewDraftTools []*entity.ToolInfo
}

type CreateDraftPluginWithCodeRequest struct {
	SpaceID     int64
	DeveloperID int64
	ProjectID   *int64
	Manifest    *model.PluginManifest
	OpenapiDoc  *model.Openapi3T
}

type CreateDraftPluginWithCodeResponse struct {
	Plugin *entity.PluginInfo
	Tools  []*entity.ToolInfo
}

type ListDraftPluginsRequest struct {
	SpaceID  int64
	APPID    int64
	PageInfo dao.PageInfo
}

type ListDraftPluginsResponse struct {
	Plugins []*entity.PluginInfo
	Total   int64
}

type CopyPluginRequest struct {
	Plugin *entity.PluginInfo
	Tools  []*entity.ToolInfo
}
