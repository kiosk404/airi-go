package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/foundation/openauth/domain/entity"
)

type permissionImpl struct{}

func NewPermissionService() Permission {
	return &permissionImpl{}
}

func (p *permissionImpl) CheckPermission(ctx context.Context, req *entity.CheckPermissionRequest) (*entity.CheckPermissionResponse, error) {
	return &entity.CheckPermissionResponse{Decision: 0}, nil
}

func (p *permissionImpl) CheckSingleAgentOperatePermission(ctx context.Context, botID int64) (bool, error) {
	return true, nil
}
