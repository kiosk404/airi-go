package runtime

import (
	"context"

	"github.com/kiosk404/airi-go/backend/api/crossdomain/singleagent"
	crossagent "github.com/kiosk404/airi-go/backend/modules/component/agent/crossdomain"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

func getAgentInfo(ctx context.Context, agentID int64, isDraft bool) (*singleagent.SingleAgent, error) {
	agentInfo, err := crossagent.DefaultSVC().ObtainAgentByIdentity(ctx, &singleagent.AgentIdentity{
		AgentID: agentID,
		IsDraft: isDraft,
	})
	if err != nil {
		return nil, err
	}

	return agentInfo, nil
}

func getAgentHistoryRounds(agentInfo *singleagent.SingleAgent) int32 {
	var conversationTurns int32 = entity.ConversationTurnsDefault
	if agentInfo != nil && agentInfo.ModelInfo != nil && agentInfo.ModelInfo.ShortMemoryPolicy != nil && ptr.From(agentInfo.ModelInfo.ShortMemoryPolicy.HistoryRound) > 0 {
		conversationTurns = ptr.From(agentInfo.ModelInfo.ShortMemoryPolicy.HistoryRound)
	}
	return conversationTurns
}
