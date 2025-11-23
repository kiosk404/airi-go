package model

import (
	"github.com/kiosk404/airi-go/backend/api/model/component/plugin_develop/common"
	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/consts"
)

func NewDefaultPluginManifest() *PluginManifest {
	return &PluginManifest{
		SchemaVersion: "v1",
		API: APIDesc{
			Type: consts.PluginTypeOfCloud,
		},
		Auth: &AuthV2{
			Type: consts.AuthzTypeOfNone,
		},
		CommonParams: map[consts.HTTPParamLocation][]*common.CommonParamSchema{
			consts.ParamInBody: {},
			consts.ParamInHeader: {
				{
					Name:  "User-Agent",
					Value: "Coze/1.0",
				},
			},
			consts.ParamInQuery: {},
		},
	}
}
