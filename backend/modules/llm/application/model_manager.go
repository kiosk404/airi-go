package application

import (
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
)

type ModelApplicationService struct {
	service.IManage
	service.IRuntime
}

func NewModelApplicationService(manageSrv service.IManage,
	runtimeSrv service.IRuntime) *ModelApplicationService {
	return &ModelApplicationService{
		IRuntime: runtimeSrv,
		IManage:  manageSrv,
	}
}
