package service

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/infra/contract/imagex"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/repo"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/service/runtime"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/pkg"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/pkg/utils/safego"
)

type runImpl struct {
	RunRecordRepo repo.RunRecordRepo
	ImagexSVC     imagex.ImageX
}

func NewService(repo repo.RunRecordRepo, imagexSVC imagex.ImageX) Run {
	return &runImpl{
		RunRecordRepo: repo,
		ImagexSVC:     imagexSVC,
	}
}

func (c *runImpl) AgentRun(ctx context.Context, arm *entity.AgentRunMeta) (*schema.StreamReader[*entity.AgentRunResponse], error) {
	sr, sw := schema.Pipe[*entity.AgentRunResponse](20)

	defer func() {
		if pe := recover(); pe != nil {
			logs.ErrorX(pkg.ModelName, "panic recover: %v\n, [stack]:%v", pe, string(debug.Stack()))
			return
		}
	}()

	art := &runtime.AgentRuntime{
		StartTime:     time.Now(),
		RunMeta:       arm,
		SW:            sw,
		MessageEvent:  runtime.NewMessageEvent(),
		RunProcess:    runtime.NewRunProcess(c.RunRecordRepo),
		RunRecordRepo: c.RunRecordRepo,
		ImagexClient:  c.ImagexSVC,
	}
	safego.Go(ctx, func() {
		defer sw.Close()
		_ = art.Run(ctx)
	})

	return sr, nil
}

func (c *runImpl) Delete(ctx context.Context, runID []int64) error {
	return c.RunRecordRepo.Delete(ctx, runID)
}

func (c *runImpl) List(ctx context.Context, meta *entity.ListRunRecordMeta) ([]*entity.RunRecordMeta, error) {
	return c.RunRecordRepo.List(ctx, meta)
}

func (c *runImpl) Create(ctx context.Context, runRecord *entity.AgentRunMeta) (*entity.RunRecordMeta, error) {
	return c.RunRecordRepo.Create(ctx, runRecord)
}
func (c *runImpl) Cancel(ctx context.Context, req *entity.CancelRunMeta) (*entity.RunRecordMeta, error) {
	return c.RunRecordRepo.Cancel(ctx, req)
}

func (c *runImpl) GetByID(ctx context.Context, runID int64) (*entity.RunRecordMeta, error) {
	return c.RunRecordRepo.GetByID(ctx, runID)
}
