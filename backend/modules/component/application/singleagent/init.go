package singleagent

import (
	"github.com/kiosk404/airi-go/backend/infra/contract/cache"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	singleagent "github.com/kiosk404/airi-go/backend/modules/component/agent/domain/service"
	llm "github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
)

type (
	SingleAgent = singleagent.SingleAgent
)

var SingleAgentSVC *SingleAgentApplicationService

type ServiceComponents struct {
	IDGen     idgen.IDGenerator
	DB        rdb.Provider
	Cache     cache.Cmdable
	TosClient storage.Storage

	ModelMgr llm.IManage
}
