package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/kiosk404/airi-go/backend/infra/contract/eventbus"
	"github.com/kiosk404/airi-go/backend/infra/contract/search"
	"github.com/kiosk404/airi-go/backend/modules/data/search/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/data/search/pkg"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

const projectIndexName = "project_draft"

type projectHandlerImpl struct {
	esClient search.Client
}

type ConsumerHandler = eventbus.ConsumerHandler

var defaultProjectHandle *projectHandlerImpl

func NewProjectHandler(ctx context.Context, e search.Client) ConsumerHandler {
	handler := &projectHandlerImpl{
		esClient: e,
	}

	defaultProjectHandle = handler
	return handler
}

func (s *projectHandlerImpl) HandleMessage(ctx context.Context, msg *eventbus.Message) error {
	ev := &entity.ProjectDomainEvent{}

	logs.InfoX(pkg.ModelName, "Project Handler receive: %s", string(msg.Body))
	err := sonic.Unmarshal(msg.Body, ev)
	if err != nil {
		return err
	}

	err = s.indexProject(ctx, ev)
	if err != nil {
		return err
	}

	return nil
}

func (s *projectHandlerImpl) indexProject(ctx context.Context, ev *entity.ProjectDomainEvent) error {
	if ev.Project == nil {
		return fmt.Errorf("project is nil")
	}

	if ev.Meta == nil {
		ev.Meta = &entity.EventMeta{}
	}

	ev.Meta.ReceiveTimeMs = time.Now().UnixMilli()

	switch ev.OpType {
	case entity.Created:
		return s.esClient.Create(ctx, projectIndexName, conv.Int64ToStr(ev.Project.ID), ev.Project)
	case entity.Updated:
		return s.esClient.Update(ctx, projectIndexName, conv.Int64ToStr(ev.Project.ID), ev.Project)
	case entity.Deleted:
		return s.esClient.Delete(ctx, projectIndexName, conv.Int64ToStr(ev.Project.ID))
	}

	return fmt.Errorf("unexpected op type: %v", ev.OpType)
}
