package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/entity"
)

type Message interface {
	List(ctx context.Context, req *entity.ListMeta) (*entity.ListResult, error)
	ListWithoutPair(ctx context.Context, req *entity.ListMeta) (*entity.ListResult, error)
	PreCreate(ctx context.Context, req *entity.Message) (*entity.Message, error)
	Create(ctx context.Context, req *entity.Message) (*entity.Message, error)
	GetByRunIDs(ctx context.Context, conversationID int64, runIDs []int64) ([]*entity.Message, error)
	GetByID(ctx context.Context, id int64) (*entity.Message, error)
	Edit(ctx context.Context, req *entity.Message) (*entity.Message, error)
	Delete(ctx context.Context, req *entity.DeleteMeta) error
	Broken(ctx context.Context, req *entity.BrokenMeta) error
}
