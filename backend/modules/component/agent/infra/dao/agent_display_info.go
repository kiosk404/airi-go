package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kiosk404/airi-go/backend/api/model/app/developer_api"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
)

func makeAgentDisplayInfoKey(userID, agentID int64) string {
	return fmt.Sprintf("agent_display_info:%d:%d", userID, agentID)
}

func (sa *SingleAgentDraftDAO) UpdateDisplayInfo(ctx context.Context, userID int64, e *entity.AgentDraftDisplayInfo) error {
	data, err := json.Marshal(e)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrAgentSetDraftBotDisplayInfo)
	}

	key := makeAgentDisplayInfoKey(userID, e.AgentID)

	_, err = sa.cacheClient.Set(ctx, key, data, 0).Result()
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrAgentSetDraftBotDisplayInfo)
	}

	return nil
}

// GetDisplayInfo 前端展示详情信息缓存
func (sa *SingleAgentDraftDAO) GetDisplayInfo(ctx context.Context, userID, agentID int64) (*entity.AgentDraftDisplayInfo, error) {
	key := makeAgentDisplayInfoKey(userID, agentID)
	data, err := sa.cacheClient.Get(ctx, key).Result()
	if err != nil {
		tabStatusDefault := developer_api.TabStatus_Default
		return &entity.AgentDraftDisplayInfo{
			AgentID: agentID,
			DisplayInfo: &developer_api.DraftBotDisplayInfoData{
				TabDisplayInfo: &developer_api.TabDisplayItems{
					PluginTabStatus:           &tabStatusDefault,
					WorkflowTabStatus:         &tabStatusDefault,
					KnowledgeTabStatus:        &tabStatusDefault,
					DatabaseTabStatus:         &tabStatusDefault,
					VariableTabStatus:         &tabStatusDefault,
					OpeningDialogTabStatus:    &tabStatusDefault,
					ScheduledTaskTabStatus:    &tabStatusDefault,
					SuggestionTabStatus:       &tabStatusDefault,
					TtsTabStatus:              &tabStatusDefault,
					FileboxTabStatus:          &tabStatusDefault,
					LongTermMemoryTabStatus:   &tabStatusDefault,
					AnswerActionTabStatus:     &tabStatusDefault,
					ImageflowTabStatus:        &tabStatusDefault,
					BackgroundImageTabStatus:  &tabStatusDefault,
					ShortcutTabStatus:         &tabStatusDefault,
					KnowledgeTableTabStatus:   &tabStatusDefault,
					KnowledgeTextTabStatus:    &tabStatusDefault,
					KnowledgePhotoTabStatus:   &tabStatusDefault,
					HookInfoTabStatus:         &tabStatusDefault,
					DefaultUserInputTabStatus: &tabStatusDefault,
				},
			},
		}, nil
	}

	e := &entity.AgentDraftDisplayInfo{}
	err = json.Unmarshal([]byte(data), e)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrAgentGetDraftBotDisplayInfoNotFound)
	}

	return e, nil
}
