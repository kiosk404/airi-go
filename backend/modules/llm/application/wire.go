//go:build wireinject
// +build wireinject

package application

import (
	"context"

	"github.com/google/wire"
	"github.com/kiosk404/airi-go/backend/api/model/llm/manage"
	"github.com/kiosk404/airi-go/backend/api/model/llm/runtime"
	"github.com/kiosk404/airi-go/backend/infra/contract/cache"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/limiter"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llmfactory"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/config"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/repo"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/repo/dao"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/rpc"
	"github.com/kiosk404/airi-go/backend/pkg/conf"
)

var (
	llmDomainSet = wire.NewSet(
		llmfactory.NewFactory,
		config.NewManage,
		config.NewRuntime,
		service.NewRuntime,
		service.NewManage,
		repo.NewRuntimeRepo,
		dao.NewModelRequestRecordDao,
		rpc.NewAuthRPCProvider,
	)
	runtimeSet = wire.NewSet(
		NewRuntimeApplication,
		llmDomainSet,
	)
	manageSet = wire.NewSet(
		NewManageApplication,
		llmDomainSet,
	)
)

func InitRuntimeApplication(
	ctx context.Context,
	idGen idgen.IDGenerator,
	configFactory conf.IConfigLoaderFactory,
	db rdb.Provider,
	redis cache.Cmdable,
	factory limiter.IRateLimiterFactory) (runtime.LLMRuntimeService, error) {
	wire.Build(runtimeSet)
	return nil, nil
}

func InitManageApplication(
	ctx context.Context,
	configFactory conf.IConfigLoaderFactory) (manage.LLMManageService, error) {
	wire.Build(manageSet)
	return nil, nil
}
