package application

import (
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/imagex"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/application/singleagent"
	agentRepo "github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/repo"
	agentrun "github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/service"
	convRepo "github.com/kiosk404/airi-go/backend/modules/conversation/conversation/domain/repo"
	conversationService "github.com/kiosk404/airi-go/backend/modules/conversation/conversation/domain/service"
	messageRepo "github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/repo"
	message "github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/service"
)

type ServiceComponents struct {
	IDGen     idgen.IDGenerator
	DB        rdb.Provider
	TosClient storage.Storage
	ImageX    imagex.ImageX

	SingleAgentDomainSVC singleagent.SingleAgent
}

func InitService(s *ServiceComponents) *ConversationApplicationService {
	agentRunDomainSVC := agentrun.NewService(agentRepo.NewRunRecordRepo(s.DB, s.IDGen), s.ImageX)        // 运行记录
	conversationDomainSVC := conversationService.NewService(convRepo.NewConversationRepo(s.DB, s.IDGen)) //
	messageDomainSVC := message.NewService(messageRepo.NewMessageRepo(s.DB, s.IDGen))

	ConversationSVC.AgentRunDomainSVC = agentRunDomainSVC
	ConversationSVC.MessageDomainSVC = messageDomainSVC
	ConversationSVC.ConversationDomainSVC = conversationDomainSVC
	ConversationSVC.appContext = s

	return &ConversationApplicationService{
		appContext: s,

		AgentRunDomainSVC:     agentRunDomainSVC,
		ConversationDomainSVC: conversationDomainSVC,
		MessageDomainSVC:      messageDomainSVC,
	}
}
