package application

import (
	"context"

	"github.com/kiosk404/airi-go/backend/api/model/conversation/conversation"
)

type ConversationApplicationService struct {
	appContext *ServiceComponents
}

func (c *ConversationApplicationService) ClearHistory(ctx context.Context, req *conversation.ClearConversationHistoryRequest) (*conversation.ClearConversationHistoryResponse, error) {
	return &conversation.ClearConversationHistoryResponse{}, nil
}

func (c *ConversationApplicationService) CreateSection(ctx context.Context, conversationID int64) (int64, error) {
	return 0, nil
}

func (c *ConversationApplicationService) CreateConversation(ctx context.Context, agentID int64, connectorID int64) (*conversation.CreateConversationResponse, error) {
	return &conversation.CreateConversationResponse{}, nil
}

func (c *ConversationApplicationService) ListConversation(ctx context.Context, req *conversation.ListConversationsApiRequest) (*conversation.ListConversationsApiResponse, error) {
	return &conversation.ListConversationsApiResponse{}, nil
}

func (c *ConversationApplicationService) DeleteConversation(ctx context.Context, req *conversation.DeleteConversationApiRequest) (*conversation.DeleteConversationApiResponse, error) {
	return &conversation.DeleteConversationApiResponse{}, nil
}

func (c *ConversationApplicationService) UpdateConversation(ctx context.Context, req *conversation.UpdateConversationApiRequest) (*conversation.UpdateConversationApiResponse, error) {
	return &conversation.UpdateConversationApiResponse{}, nil
}
