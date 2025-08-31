package application

import (
	"context"

	"github.com/kiosk404/airi-go/backend/application/base/appinfra"
	"github.com/kiosk404/airi-go/backend/application/openauth"
	"github.com/kiosk404/airi-go/backend/application/user"
)

type basicServices struct {
	infra       *appinfra.AppDependencies
	userSVC     *user.UserApplicationService
	openAuthSVC *openauth.OpenAuthApplicationService
}

type primaryServices struct {
	basicServices *basicServices

	//pluginSVC    *plugin.PluginApplicationService
	//memorySVC    *memory.MemoryApplicationServices
}

func Init(ctx context.Context) (err error) {

	return nil
}
