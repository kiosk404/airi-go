package application

import (
	"context"
	"fmt"
	openauth "github.com/kiosk404/airi-go/backend/modules/foundation/application"

	"github.com/kiosk404/airi-go/backend/application/base/appinfra"
	model "github.com/kiosk404/airi-go/backend/modules/llm/application"
	modelservice "github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llmfactory"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/config"
	modelrepo "github.com/kiosk404/airi-go/backend/modules/llm/infra/repo"
	modeldao "github.com/kiosk404/airi-go/backend/modules/llm/infra/repo/dao"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

const Application = "application"

type basicServices struct {
	infra       *appinfra.AppDependencies
	userSVC     *openauth.UserApplicationService
	openAuthSVC *openauth.OpenAuthApplicationService
	modelMgrSVC *model.ModelApplicationService
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
	var err error

	openAuthSVC := openauth.InitService(infra.DB, infra.IDGenSVC)
	userSVC := openauth.InitUserService(ctx, infra.DB, infra.TOSClient, infra.IDGenSVC)

	modelManagerConf, err := config.NewManage(ctx, infra.ConfigFactory)
	if err != nil {
		logs.ErrorX(Application, "failed to load model mgr config")
		return nil, err
	}
	modelRuntimeConf, err := config.NewRuntime(ctx, infra.ConfigFactory)
	runTimeRepo := modelrepo.NewRuntimeRepo(infra.DB, modeldao.NewModelRequestRecordDao(infra.DB))

	modelSVC := model.NewModelApplicationService(
		modelservice.NewManage(modelManagerConf),
		modelservice.NewRuntime(llmfactory.ModelF, infra.IDGenSVC, runTimeRepo, modelRuntimeConf))

	return &basicServices{
		infra:       infra,
		userSVC:     userSVC,
		openAuthSVC: openAuthSVC,
		modelMgrSVC: modelSVC,
	}, nil
}

// initPrimaryServices init primary services that depends on basic services.
func initPrimaryServices(ctx context.Context, basicServices *basicServices) (*primaryServices, error) {

	return &primaryServices{
		basicServices: basicServices,
	}, nil
}
