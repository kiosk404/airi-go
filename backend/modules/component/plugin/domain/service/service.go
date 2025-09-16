package service

import (
	"context"

	model "github.com/kiosk404/airi-go/backend/api/crossdomain/plugin"
	"github.com/kiosk404/airi-go/backend/api/model/component/plugin_develop/common"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/entity"
)

//go:generate mockgen -destination ./mock/service.go --package mockPlugin -source service.go
type PluginService interface {
	// Draft Plugin
	CreateDraftPlugin(ctx context.Context, req *CreateDraftPluginRequest) (pluginID int64, err error)
	CreateDraftPluginWithCode(ctx context.Context, req *CreateDraftPluginWithCodeRequest) (resp *CreateDraftPluginWithCodeResponse, err error)
	GetDraftPlugin(ctx context.Context, pluginID int64) (plugin *entity.PluginInfo, err error)
	MGetDraftPlugins(ctx context.Context, pluginIDs []int64) (plugins []*entity.PluginInfo, err error)
	ListDraftPlugins(ctx context.Context, req *ListDraftPluginsRequest) (resp *ListDraftPluginsResponse, err error)
	UpdateDraftPlugin(ctx context.Context, plugin *UpdateDraftPluginRequest) (err error)
	UpdateDraftPluginWithCode(ctx context.Context, req *UpdateDraftPluginWithCodeRequest) (err error)
	DeleteDraftPlugin(ctx context.Context, pluginID int64) (err error)
	DeleteAPPAllPlugins(ctx context.Context, appID int64) (pluginIDs []int64, err error)
	GetAPPAllPlugins(ctx context.Context, appID int64) (plugins []*entity.PluginInfo, err error)
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
	AuthzType    *model.AuthzType
	Location     *model.HTTPParamLocation
	Key          *string
	ServiceToken *string
	OAuthInfo    *string
	AuthzSubType *model.AuthzSubType
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
