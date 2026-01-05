package service

import (
	"context"
	"time"

	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg/consts"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

func (s singleAgentImpl) CreateSingleAgent(ctx context.Context, version string, e *entity.SingleAgent) (int64, error) {
	return s.AgentVersionRepo.Create(ctx, version, e)
}

func (s singleAgentImpl) GetPublishedTime(ctx context.Context, agentID int64) (int64, error) {
	pubInfo, err := s.GetPublishedInfo(ctx, agentID)
	if err != nil {
		return 0, err
	}

	return pubInfo.LastPublishTimeMS, nil
}

func (s singleAgentImpl) GetPublishedInfo(ctx context.Context, agentID int64) (*entity.PublishInfo, error) {
	// 获取最新的发布记录
	publishHistory, err := s.ListAgentPublishHistory(ctx, agentID, 1, 1)
	if err != nil {
		return nil, err
	}
	if len(publishHistory) == 0 {
		return &entity.PublishInfo{
			AgentID:           agentID,
			LastPublishTimeMS: 0,
		}, nil
	}
	return &entity.PublishInfo{
		AgentID:           agentID,
		LastPublishTimeMS: publishHistory[0].PublishTime,
	}, nil
}

func (s singleAgentImpl) SavePublishRecord(ctx context.Context, p *entity.SingleAgentPublish, e *entity.SingleAgent) error {
	err := s.AgentVersionRepo.SavePublishRecord(ctx, p, e)
	if err != nil {
		return err
	}

	err = s.UpdatePublishInfo(ctx, e.AgentID)
	if err != nil {
		logs.WarnX(pkg.ModelName, "update publish info failed: %v, agentID: %d", err, e.AgentID)
	}

	return nil
}

func (s *singleAgentImpl) UpdatePublishInfo(ctx context.Context, agentID int64) error {
	now := time.Now().UnixMilli()
	pubInfo, err := s.GetPublishedInfo(ctx, agentID)
	if err != nil {
		return err
	}

	if pubInfo.LastPublishTimeMS > now {
		return nil
	}

	pubInfo.LastPublishTimeMS = now
	pubInfo.AgentID = agentID

	err = s.PublishInfoRepo.Save(ctx, consts.PublishInfoKeyPrefix, conv.Int64ToStr(agentID), pubInfo)

	return err
}
