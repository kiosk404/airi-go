package dao

import (
	"context"
	"time"

	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
)

func (sa *SingleAgentVersionDAO) List(ctx context.Context, agentID int64, pageIndex, pageSize int32) ([]*entity.SingleAgentPublish, error) {
	sap := sa.dbQuery.SingleAgentPublish
	offset := (pageIndex - 1) * pageSize

	query := sap.WithContext(ctx).
		Where(sap.AgentID.Eq(agentID)).
		Order(sap.PublishTime.Desc())

	result, _, err := query.FindByPage(int(offset), int(pageSize))
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrAgentGetCode)
	}

	dos := make([]*entity.SingleAgentPublish, 0, len(result))
	for _, po := range result {
		dos = append(dos, sa.singleAgentPublishPo2Do(po))
	}

	return dos, nil
}

func (sa *SingleAgentVersionDAO) SavePublishRecord(ctx context.Context, p *entity.SingleAgentPublish, e *entity.SingleAgent) (err error) {
	publishID := p.PublishID
	version := p.Version

	id, err := sa.IDGen.GenID(ctx)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrAgentIDGenFailCode, errorx.KV("msg", "PublishDraftAgent"))
	}

	now := time.Now()

	po := &model.SingleAgentPublish{
		ID:          id,
		AgentID:     e.AgentID,
		PublishID:   publishID,
		Version:     version,
		PublishInfo: nil,
		PublishTime: now.UnixMilli(),
		Status:      0,
		Extra:       nil,
	}

	if p.PublishInfo != nil {
		po.PublishInfo = p.PublishInfo
	}

	err = sa.dbQuery.SingleAgentPublish.WithContext(ctx).Create(po)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrAgentPublishSingleAgentCode)
	}

	return nil
}

func (sa *SingleAgentVersionDAO) Create(ctx context.Context, version string, e *entity.SingleAgent) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (sa *SingleAgentVersionDAO) singleAgentPublishPo2Do(po *model.SingleAgentPublish) *entity.SingleAgentPublish {
	if po == nil {
		return nil
	}
	return &entity.SingleAgentPublish{
		ID:          po.ID,
		AgentID:     po.AgentID,
		PublishID:   po.PublishID,
		Version:     po.Version,
		PublishInfo: po.PublishInfo,
		PublishTime: po.PublishTime,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
		Status:      po.Status,
		Extra:       po.Extra,
	}
}
