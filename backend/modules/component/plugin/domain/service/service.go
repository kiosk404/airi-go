package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/api/model/component/plugin_develop/common"
	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/consts"
	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/model"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/infra/dao"
)

//go:generate mockgen -destination ./mock/service.go --package mockPlugin -source service.go
type PluginService interface {
	// Draft Plugin
	CreateDraftPlugin(ctx context.Context, req *dao.CreateDraftPluginRequest) (pluginID int64, err error)
	CreateDraftPluginWithCode(ctx context.Context, req *dao.CreateDraftPluginWithCodeRequest) (resp *dao.CreateDraftPluginWithCodeResponse, err error)
	GetDraftPlugin(ctx context.Context, pluginID int64) (plugin *entity.PluginInfo, err error)
	MGetDraftPlugins(ctx context.Context, pluginIDs []int64) (plugins []*entity.PluginInfo, err error)
	ListDraftPlugins(ctx context.Context, req *dao.ListDraftPluginsRequest) (resp *dao.ListDraftPluginsResponse, err error)
	UpdateDraftPlugin(ctx context.Context, plugin *dao.UpdateDraftPluginRequest) (err error)
	UpdateDraftPluginWithCode(ctx context.Context, req *dao.UpdateDraftPluginWithCodeRequest) (err error)
	DeleteDraftPlugin(ctx context.Context, pluginID int64) (err error)
	DeleteAPPAllPlugins(ctx context.Context, appID int64) (pluginIDs []int64, err error)
	GetAPPAllPlugins(ctx context.Context, appID int64) (plugins []*entity.PluginInfo, err error)

	// Online Plugin
	PublishPlugin(ctx context.Context, req *model.PublishPluginRequest) (err error)
	PublishAPPPlugins(ctx context.Context, req *model.PublishAPPPluginsRequest) (resp *model.PublishAPPPluginsResponse, err error)
	GetOnlinePlugin(ctx context.Context, pluginID int64) (plugin *entity.PluginInfo, err error)
	MGetOnlinePlugins(ctx context.Context, pluginIDs []int64) (plugins []*entity.PluginInfo, err error)
	MGetPluginLatestVersion(ctx context.Context, pluginIDs []int64) (resp *model.MGetPluginLatestVersionResponse, err error)
	GetPluginNextVersion(ctx context.Context, pluginID int64) (version string, err error)
	MGetVersionPlugins(ctx context.Context, versionPlugins []model.VersionPlugin) (plugins []*entity.PluginInfo, err error)
	ListCustomOnlinePlugins(ctx context.Context, pageInfo dao.PageInfo) (plugins []*entity.PluginInfo, total int64, err error)

	// Draft Tool
	MGetDraftTools(ctx context.Context, toolIDs []int64) (tools []*entity.ToolInfo, err error)
	UpdateDraftTool(ctx context.Context, req *dao.UpdateDraftToolRequest) (err error)
	ConvertToOpenapi3Doc(ctx context.Context, req *dao.ConvertToOpenapi3DocRequest) (resp *dao.ConvertToOpenapi3DocResponse)
	CreateDraftToolsWithCode(ctx context.Context, req *dao.CreateDraftToolsWithCodeRequest) (resp *dao.CreateDraftToolsWithCodeResponse, err error)
	CheckPluginToolsDebugStatus(ctx context.Context, pluginID int64) (err error)

	// Online Tool
	GetOnlineTool(ctx context.Context, toolID int64) (tool *entity.ToolInfo, err error)
	MGetOnlineTools(ctx context.Context, toolIDs []int64) (tools []*entity.ToolInfo, err error)
	MGetVersionTools(ctx context.Context, versionTools []model.VersionTool) (tools []*entity.ToolInfo, err error)
	CopyPlugin(ctx context.Context, req *dao.CopyPluginRequest) (resp *dao.CopyPluginResponse, err error)
	MoveAPPPluginToLibrary(ctx context.Context, pluginID int64) (plugin *entity.PluginInfo, err error)

	// Agent Tool
	BindAgentTools(ctx context.Context, agentID int64, bindaools []*model.BindToolInfo) (err error)
	DuplicateDraftAgentTools(ctx context.Context, fromAgentID, toAgentID int64) (err error)
	GetDraftAgentToolByName(ctx context.Context, agentID int64, pluginID int64, toolName string) (tool *entity.ToolInfo, err error)
	MGetAgentTools(ctx context.Context, req *model.MGetAgentToolsRequest) (tools []*entity.ToolInfo, err error)
	UpdateBotDefaultParams(ctx context.Context, req *dao.UpdateBotDefaultParamsRequest) (err error)

	PublishAgentTools(ctx context.Context, agentID int64, agentVersion string) (err error)

	ExecuteTool(ctx context.Context, req *model.ExecuteToolRequest, opts ...model.ExecuteToolOpt) (resp *model.ExecuteToolResponse, err error)

	// Product
	ListPluginProducts(ctx context.Context, req *dao.ListPluginProductsRequest) (resp *dao.ListPluginProductsResponse, err error)
	GetPluginProductAllTools(ctx context.Context, pluginID int64) (tools []*entity.ToolInfo, err error)

	GetOAuthStatus(ctx context.Context, userID, pluginID int64) (resp *dao.GetOAuthStatusResponse, err error)
	GetAgentPluginsOAuthStatus(ctx context.Context, userID, agentID int64) (status []*dao.AgentPluginOAuthStatus, err error)

	OAuthCode(ctx context.Context, code string, state *dao.OAuthState) (err error)
	GetAccessToken(ctx context.Context, oa *dao.OAuthInfo) (accessToken string, err error)
	RevokeAccessToken(ctx context.Context, meta *dao.AuthorizationCodeMeta) (err error)
}

type CreateDraftPluginRequest struct {
	PluginType   common.PluginType
	IconURI      string
	DeveloperID  int64
	ProjectID    *int64
	Name         string
	Desc         string
	ServerURL    string
	CommonParams map[common.ParameterLocation][]*common.CommonParamSchema
	AuthInfo     *PluginAuthInfo
}

type UpdateDraftPluginWithCodeRequest struct {
	UserID     int64
	PluginID   int64
	OpenapiDoc *model.Openapi3T
	Manifest   *entity.PluginManifest
}

type UpdateDraftPluginRequest struct {
	PluginID     int64
	Name         *string
	Desc         *string
	URL          *string
	Icon         *common.PluginIcon
	CommonParams map[common.ParameterLocation][]*common.CommonParamSchema
	AuthInfo     *PluginAuthInfo
}

type ListDraftPluginsRequest struct {
	APPID int64
}

type CreateDraftPluginWithCodeRequest struct {
	DeveloperID int64
	ProjectID   *int64
	Manifest    *entity.PluginManifest
	OpenapiDoc  *model.Openapi3T
}

type PluginAuthInfo struct {
	AuthzType    *consts.AuthzType
	Location     *consts.HTTPParamLocation
	Key          *string
	ServiceToken *string
	OAuthInfo    *string
	AuthzSubType *consts.AuthzSubType
	AuthzPayload *string
}

type CreateDraftPluginWithCodeResponse struct {
	Plugin *entity.PluginInfo
	Tools  []*entity.ToolInfo
}

type ListDraftPluginsResponse struct {
	Plugins []*entity.PluginInfo
	Total   int64
}
