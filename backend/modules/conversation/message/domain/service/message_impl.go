package service

import (
	"context"
	"sort"

	message "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/message/model"
	"github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/repo"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

type messageImpl struct {
	MessageRepo repo.MessageRepo
}

func NewService(repo repo.MessageRepo) Message {
	return &messageImpl{
		MessageRepo: repo,
	}
}

func (m *messageImpl) PreCreate(ctx context.Context, msg *entity.Message) (*entity.Message, error) {
	// create message
	return m.MessageRepo.PreCreate(ctx, msg)
}

func (m *messageImpl) Create(ctx context.Context, msg *entity.Message) (*entity.Message, error) {
	// create message
	return m.MessageRepo.Create(ctx, msg)
}

func (m *messageImpl) BatchCreate(ctx context.Context, req []*entity.Message) ([]*entity.Message, error) {
	return m.MessageRepo.BatchCreate(ctx, req)
}

func (m *messageImpl) List(ctx context.Context, req *entity.ListMeta) (*entity.ListResult, error) {
	resp := &entity.ListResult{}
	req.MessageType = []*message.MessageType{ptr.Of(message.MessageTypeQuestion)}
	// get message with query
	messageList, hasMore, err := m.MessageRepo.List(ctx, req)
	if err != nil {
		return resp, err
	}

	resp.Direction = req.Direction
	resp.HasMore = hasMore

	if len(messageList) > 0 {
		sort.Slice(messageList, func(i, j int) bool {
			return messageList[i].CreatedAt > messageList[j].CreatedAt
		})
		resp.PrevCursor = messageList[len(messageList)-1].ID
		resp.NextCursor = messageList[0].ID

		var runIDs []int64
		for _, m := range messageList {
			runIDs = append(runIDs, m.RunID)
		}
		orderBy := "DESC"
		if req.OrderBy != nil {
			orderBy = *req.OrderBy
		}
		allMessageList, err := m.MessageRepo.GetByRunIDs(ctx, runIDs, orderBy)
		if err != nil {
			return resp, err
		}
		resp.Messages = allMessageList
	}
	return resp, nil
}

func (m *messageImpl) ListWithoutPair(ctx context.Context, req *entity.ListMeta) (*entity.ListResult, error) {
	resp := &entity.ListResult{}
	messageList, hasMore, err := m.MessageRepo.List(ctx, req)
	if err != nil {
		return resp, err
	}
	resp.Direction = req.Direction
	resp.HasMore = hasMore
	resp.Messages = messageList
	if len(messageList) > 0 {
		resp.PrevCursor = messageList[0].ID
		resp.NextCursor = messageList[len(messageList)-1].ID
	}

	return resp, nil
}

func (m *messageImpl) GetByRunIDs(ctx context.Context, conversationID int64, runIDs []int64) ([]*entity.Message, error) {
	return m.MessageRepo.GetByRunIDs(ctx, runIDs, "ASC")
}

func (m *messageImpl) Edit(ctx context.Context, req *entity.Message) (*entity.Message, error) {
	_, err := m.MessageRepo.Edit(ctx, req.ID, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (m *messageImpl) Delete(ctx context.Context, req *entity.DeleteMeta) error {
	return m.MessageRepo.Delete(ctx, req)
}

func (m *messageImpl) GetByID(ctx context.Context, id int64) (*entity.Message, error) {
	return m.MessageRepo.GetByID(ctx, id)
}

func (m *messageImpl) Broken(ctx context.Context, req *entity.BrokenMeta) error {

	_, err := m.MessageRepo.Edit(ctx, req.ID, &message.Message{
		Status:   message.MessageStatusBroken,
		Position: ptr.From(req.Position),
	})
	if err != nil {
		return err
	}
	return nil
}
