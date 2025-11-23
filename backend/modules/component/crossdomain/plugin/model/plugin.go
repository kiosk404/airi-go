package model

import (
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	api "github.com/kiosk404/airi-go/backend/api/model/component/plugin_develop/common"
)

type BindToolInfo struct {
	ToolID   int64
	PluginID int64
	Source   *bot_common.PluginFrom
}

type VersionPlugin struct {
	PluginID int64
	Version  string
}

type MGetPluginLatestVersionResponse struct {
	Versions map[int64]string // pluginID vs version
}

type PluginInfo struct {
	ID           int64
	PluginType   api.PluginType
	DeveloperID  int64
	APPID        *int64
	RefProductID *int64 // for product plugin
	IconURI      *string
	ServerURL    *string
	Version      *string
	VersionDesc  *string

	CreatedAt int64
	UpdatedAt int64

	Manifest   *PluginManifest
	OpenapiDoc *Openapi3T
}

func (p PluginInfo) SetName(name string) {
	if p.Manifest == nil || p.OpenapiDoc == nil {
		return
	}
	p.Manifest.NameForModel = name
	p.Manifest.NameForHuman = name
	p.OpenapiDoc.Info.Title = name
}

func (p PluginInfo) GetName() string {
	if p.Manifest == nil {
		return ""
	}
	return p.Manifest.NameForHuman
}

func (p PluginInfo) GetDesc() string {
	if p.Manifest == nil {
		return ""
	}
	return p.Manifest.DescriptionForHuman
}

func (p PluginInfo) GetAuthInfo() *AuthV2 {
	if p.Manifest == nil {
		return nil
	}
	return p.Manifest.Auth
}

func (p PluginInfo) IsOfficial() bool {
	return p.RefProductID != nil
}

func (p PluginInfo) GetIconURI() string {
	if p.IconURI == nil {
		return ""
	}
	return *p.IconURI
}

func (p PluginInfo) Published() bool {
	return p.Version != nil
}

type PublishPluginRequest struct {
	PluginID    int64
	Version     string
	VersionDesc string
}

type PublishAPPPluginsRequest struct {
	APPID   int64
	Version string
}

type PublishAPPPluginsResponse struct {
	FailedPlugins   []*PluginInfo
	AllDraftPlugins []*PluginInfo
}

type CheckCanPublishPluginsRequest struct {
	PluginIDs []int64
	Version   string
}

type CheckCanPublishPluginsResponse struct {
	InvalidPlugins []*PluginInfo
}
