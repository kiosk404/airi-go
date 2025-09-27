package service

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/entity"
)

type Run interface {
	AgentRun(ctx context.Context, req *entity.AgentRunMeta) (*schema.StreamReader[*entity.AgentRunResponse], error)
	Delete(ctx context.Context, runID []int64) error
	Create(ctx context.Context, runRecord *entity.AgentRunMeta) (*entity.RunRecordMeta, error)
	List(ctx context.Context, ListMeta *entity.ListRunRecordMeta) ([]*entity.RunRecordMeta, error)
	GetByID(ctx context.Context, runID int64) (*entity.RunRecordMeta, error)
	Cancel(ctx context.Context, req *entity.CancelRunMeta) (*entity.RunRecordMeta, error)
}
