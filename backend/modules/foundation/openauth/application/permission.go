package application

import (
	openapiauth "github.com/kiosk404/airi-go/backend/modules/foundation/openauth/domain/service"
)

type PermissionCheckApp struct {
	permissionSVC openapiauth.Permission
}

func NewPermissionCheckAPP() *PermissionCheckApp {
	return &PermissionCheckApp{}
}

func (p *PermissionCheckApp) CheckPermission() {

}
