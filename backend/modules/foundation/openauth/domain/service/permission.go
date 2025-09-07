package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/foundation/openauth/domain/entity"
)

type Permission interface {
	CheckPermission(ctx context.Context, req *entity.CheckPermissionRequest) (*entity.CheckPermissionResponse, error)
	CheckSingleAgentOperatePermission(ctx context.Context, agentID int64) (bool, error)
}
