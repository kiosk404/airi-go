package service

import (
	"context"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
)

func (s singleAgentImpl) CreateSingleAgent(ctx context.Context, version string, e *entity.SingleAgent) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s singleAgentImpl) DuplicateInMemory(ctx context.Context, req *entity.DuplicateInfo) (newAgent *entity.SingleAgent, err error) {
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
