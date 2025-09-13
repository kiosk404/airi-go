package service

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/repo"
	llm "github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llmimpl/chatmodel"
	"github.com/kiosk404/airi-go/backend/pkg/jsoncache"
)

type singleAgentImpl struct {
	ModelMgr     llm.IManage
	ModelFactory chatmodel.Factory

	AgentDraftRepo   repo.SingleAgentDraftRepo
	AgentVersionRepo repo.SingleAgentVersionRepo
	PublishInfoRepo  *jsoncache.JsonCache[entity.PublishInfo]

	CPStore compose.CheckPointStore
}

func NewService() SingleAgent {
	s := &singleAgentImpl{}

	return s
}

func (s singleAgentImpl) CreateSingleAgentDraft(ctx context.Context, creatorID int64, draft *entity.SingleAgent) (agentID int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) CreateSingleAgentDraftWithID(ctx context.Context, creatorID, agentID int64, draft *entity.SingleAgent) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) MGetSingleAgentDraft(ctx context.Context, agentIDs []int64) (agents []*entity.SingleAgent, err error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) GetSingleAgentDraft(ctx context.Context, agentID int64) (agentInfo *entity.SingleAgent, err error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) UpdateSingleAgentDraft(ctx context.Context, agentInfo *entity.SingleAgent) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) DeleteAgentDraft(ctx context.Context, agentID int64) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) UpdateAgentDraftDisplayInfo(ctx context.Context, userID int64, e *entity.AgentDraftDisplayInfo) error {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) GetAgentDraftDisplayInfo(ctx context.Context, userID, agentID int64) (*entity.AgentDraftDisplayInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) CreateSingleAgent(ctx context.Context, version string, e *entity.SingleAgent) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) DuplicateInMemory(ctx context.Context, req *entity.DuplicateInfo) (newAgent *entity.SingleAgent, err error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) StreamExecute(ctx context.Context, req *entity.ExecuteRequest) (events *schema.StreamReader[*entity.AgentEvent], err error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) GetSingleAgent(ctx context.Context, agentID int64, version string) (botInfo *entity.SingleAgent, err error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) ListAgentPublishHistory(ctx context.Context, agentID int64, pageIndex, pageSize int32) ([]*entity.SingleAgentPublish, error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) ObtainAgentByIdentity(ctx context.Context, identity *entity.AgentIdentity) (*entity.SingleAgent, error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) GetPublishedTime(ctx context.Context, agentID int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) GetPublishedInfo(ctx context.Context, agentID int64) (*entity.PublishInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) SavePublishRecord(ctx context.Context, p *entity.SingleAgentPublish, e *entity.SingleAgent) error {
	//TODO implement me
	panic("implement me")
}
