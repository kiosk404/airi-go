package service

import (
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/limiter"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/component/conf"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/repo"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llmfactory"
)

type ModelService struct {
	IManage
	IRuntime
}

func NewRuntime(
	llmFact llmfactory.IFactory,
	idGen idgen.IDGenerator,
	runtimeRepo repo.IRuntimeRepository,
	cfg conf.IConfigRuntime,
	limiter limiter.IRateLimiter,
) IRuntime {
	return &RuntimeImpl{
		llmFact:     llmFact,
		idGen:       idGen,
		runtimeRepo: runtimeRepo,
		runtimeCfg:  cfg,
		limiter:     limiter,
	}
}

func NewManage(cfg conf.IConfigManage) IManage {
	return &ManageImpl{
		conf: cfg,
	}
}
