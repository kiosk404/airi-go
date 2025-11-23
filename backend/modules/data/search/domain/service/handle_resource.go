package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/kiosk404/airi-go/backend/infra/contract/eventbus"
	"github.com/kiosk404/airi-go/backend/infra/contract/search"
	"github.com/kiosk404/airi-go/backend/modules/data/search/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

const resourceIndexName = "airi_resource"

type resourceHandlerImpl struct {
	esClient search.Client
}

var defaultResourceHandler *resourceHandlerImpl

func NewResourceHandler(ctx context.Context, e search.Client) ConsumerHandler {
	handler := &resourceHandlerImpl{
		esClient: e,
	}

	defaultResourceHandler = handler
	return handler
}

func (s *resourceHandlerImpl) HandleMessage(ctx context.Context, msg *eventbus.Message) error {
	ev := &entity.ResourceDomainEvent{}

	logs.InfoX("Resource Handler receive: %s", string(msg.Body))

	err := sonic.Unmarshal(msg.Body, ev)
	if err != nil {
		return err
	}

	err = s.indexResources(ctx, ev)
	if err != nil {
		return err
	}

	return nil
}

func (s *resourceHandlerImpl) indexResources(ctx context.Context, ev *entity.ResourceDomainEvent) error {
	if ev.Meta == nil {
		ev.Meta = &entity.EventMeta{}
	}

	ev.Meta.ReceiveTimeMs = time.Now().UnixMilli()

	return s.indexResource(ctx, ev.OpType, ev.Resource)
}

func (s *resourceHandlerImpl) indexResource(ctx context.Context, opType entity.OpType, r *entity.ResourceDocument) error {
	switch opType {
	case entity.Created:
		return s.esClient.Create(ctx, resourceIndexName, conv.Int64ToStr(r.ResID), r)
	case entity.Updated:
		return s.esClient.Update(ctx, resourceIndexName, conv.Int64ToStr(r.ResID), r)
	case entity.Deleted:
		return s.esClient.Delete(ctx, resourceIndexName, conv.Int64ToStr(r.ResID))
	}

	return fmt.Errorf("unexpected op type: %v", opType)
}
