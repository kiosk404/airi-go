package service

import (
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/component/conf"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/repo"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llmfactory"
)

func NewRuntime(
	llmFact llmfactory.IFactory,
	idGen idgen.IDGenerator,
	runtimeRepo repo.IRuntimeRepo,
	cfg conf.IConfigRuntime,
) IRuntime {
	return &RuntimeImpl{
		llmFact:     llmFact,
		idGen:       idGen,
		runtimeRepo: runtimeRepo,
		runtimeCfg:  cfg,
	}
}

func NewManage(cfg conf.IConfigManage) IManage {
	return &ManageImpl{
		conf: cfg,
	}
}
