package application

import (
	"context"
	"fmt"

	"github.com/kiosk404/airi-go/backend/application/appinfra"
	"github.com/kiosk404/airi-go/backend/infra/contract/eventbus"
	implEventbus "github.com/kiosk404/airi-go/backend/infra/impl/eventbus"
	singleagentapp "github.com/kiosk404/airi-go/backend/modules/component/agent/application/singleagent"
	crossagent "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/agent"
	crossagentimpl "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/agent/impl"
	crossplugin "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin"
	crosspluginimpl "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/impl"
	pluginapp "github.com/kiosk404/airi-go/backend/modules/component/plugin/application"
	conversationapp "github.com/kiosk404/airi-go/backend/modules/conversation/conversation/application"
	crossmessage "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/message"
	crossmessageimpl "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/message/impl"
	crosssearch "github.com/kiosk404/airi-go/backend/modules/data/crossdomain/search"
	searchImpl "github.com/kiosk404/airi-go/backend/modules/data/crossdomain/search/impl"
	searchapp "github.com/kiosk404/airi-go/backend/modules/data/search/application"
	search "github.com/kiosk404/airi-go/backend/modules/data/search/domain/service"
	uploadapp "github.com/kiosk404/airi-go/backend/modules/data/upload/application"
	openauthapp "github.com/kiosk404/airi-go/backend/modules/foundation/openauth/application"
	userapp "github.com/kiosk404/airi-go/backend/modules/foundation/user/application"
	modelmgrapp "github.com/kiosk404/airi-go/backend/modules/llm/application"
	crossmodelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr"
	crossmodelmgrimpl "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/impl"
	"github.com/kiosk404/airi-go/backend/pkg/checkpoint"
)

type eventbusImpl struct {
	resourceEventBus search.ResourceEventBus
	projectEventBus  search.ProjectEventBus
}

type basicServices struct {
	infra       *appinfra.AppDependencies
	eventbus    *eventbusImpl
	userSVC     *userapp.UserApplicationService
	openAuthSVC *openauthapp.OpenAuthApplicationService
	modelMgrSVC *modelmgrapp.ModelManagerApplicationService
	uploadSVC   *uploadapp.UploadService
}

type primaryServices struct {
	basicServices *basicServices

	pluginSVC *pluginapp.PluginApplicationService
	//memorySVC    *memory.MemoryApplicationServices
}

type complexServices struct {
	primaryServices *primaryServices
	singleAgentSVC  *singleagentapp.SingleAgentApplicationService
	//appSVC          *app.APPApplicationService
	searchSVC       *searchapp.SearchApplicationService
	conversationSVC *conversationapp.ConversationApplicationService
}

func Init(ctx context.Context) (err error) {
	infra, err := appinfra.Init(ctx)
	if err != nil {
		return err
	}

	eventBus := initEventBus(infra)

	basicServices, err := initBasicServices(ctx, infra, eventBus)
	if err != nil {
		return fmt.Errorf("init - initBasicServices failed, err: %v", err)
	}

	primaryServices, err := initPrimaryServices(ctx, basicServices)
	if err != nil {
		return fmt.Errorf("init - initPrimaryServices failed, err: %v", err)
	}

	complexServices, err := initComplexServices(ctx, primaryServices)
	if err != nil {
		return fmt.Errorf("init - initVitalServices failed, err: %v", err)
	}

	crossmodelmgr.SetDefaultSVC(crossmodelmgrimpl.InitDomainService(basicServices.modelMgrSVC.DomainSVC))
	crossplugin.SetDefaultSVC(crosspluginimpl.InitDomainService(primaryServices.pluginSVC.DomainSVC, infra.TOSClient))
	crossagent.SetDefaultSVC(crossagentimpl.InitDomainService(complexServices.singleAgentSVC.DomainSVC, infra.ImageXClient))
	crossmessage.SetDefaultSVC(crossmessageimpl.InitDomainService(complexServices.conversationSVC.MessageDomainSVC))
	crosssearch.SetDefaultSVC(searchImpl.InitDomainService(complexServices.searchSVC.DomainSVC))

	return nil
}

func initEventBus(infra *appinfra.AppDependencies) *eventbusImpl {
	e := &eventbusImpl{}
	eventbus.SetDefaultSVC(implEventbus.NewConsumerService())
	return e
}

// initBasicServices init basic services that only depends on infra.
func initBasicServices(ctx context.Context, infra *appinfra.AppDependencies, e *eventbusImpl) (*basicServices, error) {
	var err error

	openAuthSVC := openauthapp.InitService(infra.DB, infra.IDGenSVC)
	userSVC := userapp.InitService(ctx, infra.DB, infra.TOSClient, infra.IDGenSVC)
	modelSVC := modelmgrapp.InitService(ctx, infra.IDGenSVC, infra.DB, infra.TOSClient, infra.ConfigFactory)
	uploadSVC := uploadapp.InitService(ctx, infra.TOSClient, infra.CacheCli, infra.DB, infra.IDGenSVC)

	return &basicServices{
		eventbus:    e,
		infra:       infra,
		userSVC:     userSVC,
		openAuthSVC: openAuthSVC,
		modelMgrSVC: modelSVC,
		uploadSVC:   uploadSVC,
	}, err
}

// initPrimaryServices init primary services that depends on basic services.
func initPrimaryServices(ctx context.Context, basicServices *basicServices) (*primaryServices, error) {
	pluginSVC, err := pluginapp.InitService(ctx, basicServices.toPluginServiceComponents())
	if err != nil {
		return nil, err
	}
	return &primaryServices{
		basicServices: basicServices,
		pluginSVC:     pluginSVC,
	}, nil
}

func initComplexServices(ctx context.Context, primaryServices *primaryServices) (*complexServices, error) {
	singleAgentSVC, err := singleagentapp.InitService(primaryServices.toSingleAgentServiceComponents())
	if err != nil {
		return nil, err
	}

	searchSVC, err := searchapp.InitService(ctx, primaryServices.toSearchComponents(singleAgentSVC))
	if err != nil {
		return nil, err
	}
	conversationSVC := conversationapp.InitService(primaryServices.toConversationComponents(singleAgentSVC))

	return &complexServices{
		primaryServices: primaryServices,
		singleAgentSVC:  singleAgentSVC,
		conversationSVC: conversationSVC,
		searchSVC:       searchSVC,
	}, nil
}

func (p *basicServices) toPluginServiceComponents() *pluginapp.ServiceComponents {
	return &pluginapp.ServiceComponents{
		IDGen:    p.infra.IDGenSVC,
		DB:       p.infra.DB,
		EventBus: p.eventbus.resourceEventBus,
		OSS:      p.infra.TOSClient,
		UserSVC:  p.userSVC.DomainSVC,
	}
}

func (p *primaryServices) toSingleAgentServiceComponents() *singleagentapp.ServiceComponents {
	return &singleagentapp.ServiceComponents{
		IDGen:       p.basicServices.infra.IDGenSVC,
		DB:          p.basicServices.infra.DB,
		Cache:       p.basicServices.infra.CacheCli,
		TosClient:   p.basicServices.infra.TOSClient,
		CPStore:     checkpoint.NewInMemoryStore(),
		ModelMgrSVC: p.basicServices.modelMgrSVC.DomainSVC,
	}
}

func (p *primaryServices) toConversationComponents(singleAgentSVC *singleagentapp.SingleAgentApplicationService) *conversationapp.ServiceComponents {
	infra := p.basicServices.infra

	return &conversationapp.ServiceComponents{
		DB:                   infra.DB,
		IDGen:                infra.IDGenSVC,
		TosClient:            infra.TOSClient,
		ImageX:               infra.ImageXClient,
		SingleAgentDomainSVC: singleAgentSVC.DomainSVC,
	}
}

func (p *primaryServices) toSearchComponents(singleAgentSVC *singleagentapp.SingleAgentApplicationService) *searchapp.ServiceComponents {
	infra := p.basicServices.infra
	return &searchapp.ServiceComponents{
		DB:                   infra.DB,
		Cache:                infra.CacheCli,
		TOS:                  infra.TOSClient,
		SingleAgentDomainSVC: singleAgentSVC.DomainSVC,
		PluginDomainSVC:      p.pluginSVC.DomainSVC,
		UserDomainSVC:        p.basicServices.userSVC.DomainSVC,
	}
}
