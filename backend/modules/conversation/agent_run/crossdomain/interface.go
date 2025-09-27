package crossdomain

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/entity"
)

type AgentRun interface {
	Delete(ctx context.Context, runID []int64) error
	List(ctx context.Context, ListMeta *entity.ListRunRecordMeta) ([]*entity.RunRecordMeta, error)
	Create(ctx context.Context, runRecord *entity.AgentRunMeta) (*entity.RunRecordMeta, error)
}

var defaultSVC AgentRun

func DefaultSVC() AgentRun {
	return defaultSVC
}

func SetDefaultSVC(svc AgentRun) {
	defaultSVC = svc
}
