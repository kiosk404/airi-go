package application

import (
	"context"

	"github.com/kiosk404/airi-go/backend/api/model/resource/common"
)

var defaultAction = []*common.ResourceAction{
	{
		Key:    common.ActionKey_Edit,
		Enable: true,
	},
	{
		Key:    common.ActionKey_Delete,
		Enable: true,
	},
	{
		Key:    common.ActionKey_Copy,
		Enable: true,
	},
}

type ResourcePacker interface {
	GetDataInfo(ctx context.Context) (*dataInfo, error)
	GetActions(ctx context.Context) []*common.ResourceAction
	GetProjectDefaultActions(ctx context.Context) []*common.ProjectResourceAction
}

type dataInfo struct {
	iconURI *string
	iconURL string
	desc    *string
	status  *int32
}
