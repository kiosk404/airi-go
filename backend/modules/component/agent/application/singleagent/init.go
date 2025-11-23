package singleagent

import (
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
	modelmgrapp "github.com/kiosk404/airi-go/backend/modules/llm/application"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llmimpl/chatmodel"
	"github.com/kiosk404/airi-go/backend/pkg/jsoncache"
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
	ModelMgr      modelmgrapp.ModelApplicationService
}

func InitService(c *ServiceComponents) (*SingleAgentApplicationService, error) {
	modelFactory := chatmodel.NewDefaultFactory()
	agentDraft := repo.NewSingleAgentRepo(c.DB, c.IDGen, c.Cache)
	agentVersion := repo.NewSingleAgentVersionRepo(c.DB, c.IDGen)
	publishInfoRepo := jsoncache.New[entity.PublishInfo]("agent:publish:last:", c.Cache)
	cps := c.CPStore

	singleAgentDomainSVC := singleagent.NewService(modelFactory, agentDraft, agentVersion, publishInfoRepo, cps)
	SingleAgentSVC = newApplicationService(c, singleAgentDomainSVC)

	return SingleAgentSVC, nil
}
