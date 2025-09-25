package application

import (
	"context"

	"github.com/kiosk404/airi-go/backend/api/model/conversation/run"
	"github.com/kiosk404/airi-go/backend/pkg/http/sse"
)

func (c *ConversationApplicationService) Run(ctx context.Context, sseSender *sse.SSenderImpl, ar *run.AgentRunRequest) error {
	return nil
}
