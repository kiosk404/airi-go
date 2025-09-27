package runtime

import (
	"context"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/api/crossdomain/agentrun"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/repo"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/pkg"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

type RunProcess struct {
	event         *Event
	SW            *schema.StreamWriter[*entity.AgentRunResponse]
	RunRecordRepo repo.RunRecordRepo
}

func NewRunProcess(runRecordRepo repo.RunRecordRepo) *RunProcess {
	return &RunProcess{
		RunRecordRepo: runRecordRepo,
	}
}

func (r *RunProcess) StepToCreate(ctx context.Context, srRecord *entity.ChunkRunItem, sw *schema.StreamWriter[*entity.AgentRunResponse]) {
	srRecord.Status = entity.RunStatusCreated
	r.event.SendRunEvent(entity.RunEventCreated, srRecord, sw)
}
func (r *RunProcess) StepToInProgress(ctx context.Context, srRecord *entity.ChunkRunItem, sw *schema.StreamWriter[*entity.AgentRunResponse]) error {
	srRecord.Status = entity.RunStatusInProgress

	updateMeta := &entity.UpdateMeta{
		Status:    entity.RunStatusInProgress,
		UpdatedAt: time.Now().UnixMilli(),
	}
	err := r.RunRecordRepo.UpdateByID(ctx, srRecord.ID, updateMeta)

	if err != nil {
		return err
	}

	r.event.SendRunEvent(entity.RunEventInProgress, srRecord, sw)
	return nil
}

func (r *RunProcess) StepToComplete(ctx context.Context, srRecord *entity.ChunkRunItem, sw *schema.StreamWriter[*entity.AgentRunResponse], usage *agentrun.Usage) {

	completedAt := time.Now().UnixMilli()

	updateMeta := &entity.UpdateMeta{
		Status:      entity.RunStatusCompleted,
		Usage:       usage,
		CompletedAt: completedAt,
		UpdatedAt:   completedAt,
	}
	err := r.RunRecordRepo.UpdateByID(ctx, srRecord.ID, updateMeta)
	if err != nil {
		logs.Error(pkg.ModelName, "RunRecordRepo.UpdateByID error: %v", err)
		r.event.SendErrEvent(entity.RunEventError, sw, &entity.RunError{
			Code: errno.ErrConversationAgentRunError,
			Msg:  err.Error(),
		})
		return
	}

	srRecord.CompletedAt = completedAt
	srRecord.Status = entity.RunStatusCompleted

	r.event.SendRunEvent(entity.RunEventCompleted, srRecord, sw)

	r.event.SendStreamDoneEvent(sw)
}
func (r *RunProcess) StepToFailed(ctx context.Context, srRecord *entity.ChunkRunItem, sw *schema.StreamWriter[*entity.AgentRunResponse]) {

	nowTime := time.Now().UnixMilli()
	updateMeta := &entity.UpdateMeta{
		Status:    entity.RunStatusFailed,
		UpdatedAt: nowTime,
		FailedAt:  nowTime,
		LastError: srRecord.Error,
	}

	err := r.RunRecordRepo.UpdateByID(ctx, srRecord.ID, updateMeta)

	if err != nil {
		r.event.SendErrEvent(entity.RunEventError, sw, &entity.RunError{
			Code: errno.ErrConversationAgentRunError,
			Msg:  err.Error(),
		})
		logs.ErrorX(pkg.ModelName, "update run record failed, err: %v", err)
		return
	}
	srRecord.Status = entity.RunStatusFailed
	srRecord.FailedAt = time.Now().UnixMilli()
	r.event.SendErrEvent(entity.RunEventError, sw, &entity.RunError{
		Code: srRecord.Error.Code,
		Msg:  srRecord.Error.Msg,
	})
}

func (r *RunProcess) StepToDone(sw *schema.StreamWriter[*entity.AgentRunResponse]) {
	r.event.SendStreamDoneEvent(sw)
}
