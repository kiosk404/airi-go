package application

import (
	"context"

	"github.com/kiosk404/airi-go/backend/api/model/llm/manage"
	"github.com/kiosk404/airi-go/backend/api/model/llm/runtime"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/limiter"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	modelservice "github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llmfactory"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llmimpl/chatmodel"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/config"
	llmmodelrepo "github.com/kiosk404/airi-go/backend/modules/llm/infra/repo"
	llmmodeldao "github.com/kiosk404/airi-go/backend/modules/llm/infra/repo/dao"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/rpc"
	"github.com/kiosk404/airi-go/backend/modules/llm/pkg"
	"github.com/kiosk404/airi-go/backend/pkg/conf"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

type ModelApplicationService struct {
	DomainSVC    ModelServiceComponents
	ChatModelSVC chatmodel.Factory
}

type ModelServiceComponents struct {
	runtime.LLMRuntimeService
	manage.LLMManageService
}

var ModelSVC *ModelApplicationService

func InitService(ctx context.Context, conf conf.IConfigLoaderFactory, idGen idgen.IDGenerator,
	rdb rdb.Provider, limiter limiter.IRateLimiterFactory) *ModelApplicationService {

	modelManagerConf, err := config.NewManage(ctx, conf)
	if err != nil {
		logs.ErrorX(pkg.ModelName, "failed to load model mgr config")
		return nil
	}
	modelRuntimeConf, err := config.NewRuntime(ctx, conf)
	runTimeRepo := llmmodelrepo.NewRuntimeRepo(rdb, llmmodeldao.NewModelRequestRecordDao(rdb))

	manager := modelservice.NewManage(modelManagerConf)
	rateLimiter := limiter.NewRateLimiter()
	llmRuntime := modelservice.NewRuntime(llmfactory.ModelF, idGen, runTimeRepo, modelRuntimeConf, rateLimiter)

	ModelSVC = &ModelApplicationService{
		DomainSVC: ModelServiceComponents{
			LLMRuntimeService: NewRuntimeApplication(manager, llmRuntime, limiter),
			LLMManageService:  NewManageApplication(manager, rpc.NewAuthRPCProvider()),
		},
		ChatModelSVC: chatmodel.NewDefaultFactory(),
	}

	return ModelSVC
}
