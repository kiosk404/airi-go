package application

import (
	"context"
	"fmt"
	"os"

	"github.com/kiosk404/airi-go/backend/infra/contract/cache"
	"github.com/kiosk404/airi-go/backend/infra/contract/eventbus"
	searchC "github.com/kiosk404/airi-go/backend/infra/contract/search"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/application/singleagent"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/service"
	prompt "github.com/kiosk404/airi-go/backend/modules/component/prompt/domain/service"
	search "github.com/kiosk404/airi-go/backend/modules/data/search/domain/service"
	user "github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/service"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/types/consts"
	"gorm.io/gorm"
)

type ServiceComponents struct {
	DB                   *gorm.DB
	Cache                cache.Cmdable
	TOS                  storage.Storage
	ESClient             searchC.Client
	ProjectEventBus      ProjectEventBus
	ResourceEventBus     ResourceEventBus
	SingleAgentDomainSVC singleagent.SingleAgent
	//APPDomainSVC         app.AppService
	//KnowledgeDomainSVC   knowledge.Knowledge
	//WorkflowDomainSVC    workflow.Service
	PluginDomainSVC service.PluginService
	UserDomainSVC   user.User
	PromptDomainSVC prompt.Prompt
}

func InitService(ctx context.Context, s *ServiceComponents) (*SearchApplicationService, error) {
	searchDomainSVC := search.NewDomainService(ctx, s.ESClient)

	SearchSVC.DomainSVC = searchDomainSVC
	SearchSVC.ServiceComponents = s

	// setup consumer
	searchConsumer := search.NewProjectHandler(ctx, s.ESClient)

	logs.Info("start search domain consumer...")
	nameServer := os.Getenv(consts.MQServer)

	err := eventbus.GetDefaultSVC().RegisterConsumer(nameServer, consts.RMQTopicApp, consts.RMQConsumeGroupApp, searchConsumer)
	if err != nil {
		return nil, fmt.Errorf("register search consumer failed, err=%w", err)
	}

	searchResourceConsumer := search.NewResourceHandler(ctx, s.ESClient)

	err = eventbus.GetDefaultSVC().RegisterConsumer(nameServer, consts.RMQTopicResource, consts.RMQConsumeGroupResource, searchResourceConsumer)
	if err != nil {
		return nil, fmt.Errorf("register search consumer failed, err=%w", err)
	}

	return SearchSVC, nil
}

type (
	ResourceEventBus = search.ResourceEventBus
	ProjectEventBus  = search.ProjectEventBus
)

func NewResourceEventBus(p eventbus.Producer) search.ResourceEventBus {
	return search.NewResourceEventBus(p)
}

func NewProjectEventBus(p eventbus.Producer) search.ProjectEventBus {
	return search.NewProjectEventBus(p)
}
