package singleagent

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/kiosk404/airi-go/backend/infra/contract/cache"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/imagex"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/repo"
	singleagent "github.com/kiosk404/airi-go/backend/modules/component/agent/domain/service"
	search "github.com/kiosk404/airi-go/backend/modules/data/search/domain/service"
	user "github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/service"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
	"github.com/kiosk404/airi-go/backend/pkg/kvstore"
)

type (
	SingleAgent = singleagent.SingleAgent
)

var SingleAgentSVC *SingleAgentApplicationService

type ServiceComponents struct {
	IDGen         idgen.IDGenerator
	DB            rdb.Provider
	Cache         cache.Cmdable
	TosClient     storage.Storage
	ImageX        imagex.ImageX
	UserDomainSVC user.User
	CPStore       compose.CheckPointStore
	EventBus      search.ProjectEventBus
	ModelMgrSVC   modelmgr.ModelManager
}

func InitService(c *ServiceComponents) (*SingleAgentApplicationService, error) {
	agentDraft := repo.NewSingleAgentRepo(c.DB, c.IDGen, c.Cache)
	agentVersion := repo.NewSingleAgentVersionRepo(c.DB, c.IDGen)
	publishInfoRepo := kvstore.New[entity.PublishInfo](c.DB.NewSession(context.Background()).DB())
	cps := c.CPStore

	singleAgentDomainSVC := singleagent.NewService(agentDraft, agentVersion, publishInfoRepo, cps)
	SingleAgentSVC = newApplicationService(c, singleAgentDomainSVC)

	return SingleAgentSVC, nil
}
