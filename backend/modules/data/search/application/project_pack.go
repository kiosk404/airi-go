package application

import (
	"context"

	"github.com/kiosk404/airi-go/backend/api/model/app/intelligence"
	"github.com/kiosk404/airi-go/backend/api/model/app/intelligence/common"
)

type projectInfo struct {
	iconURI string
	desc    string
}

type ProjectPacker interface {
	GetProjectInfo(ctx context.Context) (*projectInfo, error)
	GetPermissionInfo() *intelligence.IntelligencePermissionInfo
	GetPublishedInfo(ctx context.Context) *intelligence.IntelligencePublishInfo
	GetUserInfo(ctx context.Context, userID int64) *common.User
}
