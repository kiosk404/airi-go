package application

import (
	"context"
	"fmt"

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
	infra, err := appinfra.Init(ctx)
	if err != nil {
		return err
	}
	basicServices, err := initBasicServices(ctx, infra)
	if err != nil {
		return fmt.Errorf("init - initBasicServices failed, err: %v", err)
	}

	_, err = initPrimaryServices(ctx, basicServices)
	if err != nil {
		return fmt.Errorf("init - initPrimaryServices failed, err: %v", err)
	}

	return nil
}

// initBasicServices init basic services that only depends on infra.
func initBasicServices(ctx context.Context, infra *appinfra.AppDependencies) (*basicServices, error) {
	openAuthSVC := openauth.InitService(infra.DB, infra.IDGenSVC)
	userSVC := user.InitService(ctx, infra.DB, infra.TOSClient, infra.IDGenSVC)

	return &basicServices{
		infra:       infra,
		userSVC:     userSVC,
		openAuthSVC: openAuthSVC,
	}, nil
}

// initPrimaryServices init primary services that depends on basic services.
func initPrimaryServices(ctx context.Context, basicServices *basicServices) (*primaryServices, error) {

	return &primaryServices{
		basicServices: basicServices,
	}, nil
}
