package plugin

import (
	"github.com/kiosk404/airi-go/backend/api/model/component/plugin_develop/common"
)

type ToolInfo struct {
	ID        int64
	PluginID  int64
	CreatedAt int64
	UpdatedAt int64
	Version   *string

	ActivatedStatus *ActivatedStatus
	DebugStatus     *common.APIDebugStatus

	Method    *string
	SubURL    *string
	Operation *Openapi3Operation
}
