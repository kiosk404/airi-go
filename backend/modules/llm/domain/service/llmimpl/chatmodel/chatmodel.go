package chatmodel

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
)

type BaseChatModel = model.BaseChatModel

type ToolCallingChatModel = model.ToolCallingChatModel

type Factory interface {
	CreateChatModel(ctx context.Context, protocol entity.Protocol, config *Config) (ToolCallingChatModel, error)
	SupportProtocol(protocol entity.Protocol) bool
}
