package application

import (
	"context"
	"strconv"

	"github.com/kiosk404/airi-go/backend/api/model/conversation/common"
	"github.com/kiosk404/airi-go/backend/api/model/conversation/message"
	"github.com/kiosk404/airi-go/backend/api/model/conversation/run"
	"github.com/kiosk404/airi-go/backend/application/ctxutil"
	singleAgentEntity "github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	convEntity "github.com/kiosk404/airi-go/backend/modules/conversation/conversation/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/pkg/errno"
	model "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/message/model"
	"github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/json"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
)

func (c *ConversationApplicationService) GetMessageList(ctx context.Context, mr *message.GetMessageListRequest) (*message.GetMessageListResponse, error) {
	// Get Conversation ID by agent id & userID & scene
	userID := ctxutil.GetUIDFromCtx(ctx)

	agentID, err := strconv.ParseInt(mr.BotID, 10, 64)
	if err != nil {
		return nil, err
	}

	// 获取或者新建一个会话
	currentConversation, isNewCreate, err := c.getCurrentConversation(ctx, *userID, agentID, *mr.Scene)
	if err != nil {
		return nil, err
	}

	if isNewCreate {
		return &message.GetMessageListResponse{
			MessageList:    []*message.ChatMessage{},
			Cursor:         mr.Cursor,
			NextCursor:     "0",
			NextHasMore:    false,
			ConversationID: strconv.FormatInt(currentConversation.ID, 10),
			LastSectionID:  ptr.Of(strconv.FormatInt(currentConversation.SectionID, 10)),
		}, nil
	}

	cursor, err := strconv.ParseInt(mr.Cursor, 10, 64)
	if err != nil {
		return nil, err
	}

	// 查找历史的会话内容
	mListMessages, err := c.MessageDomainSVC.List(ctx, &entity.ListMeta{
		ConversationID: currentConversation.ID,
		AgentID:        agentID,
		Limit:          int(mr.Count),
		Cursor:         cursor,
		Direction:      loadDirectionToScrollDirection(mr.LoadDirection),
	})
	if err != nil {
		return nil, err
	}

	// get agent id
	var agentIDs []int64
	for _, mOne := range mListMessages.Messages {
		agentIDs = append(agentIDs, mOne.AgentID)
	}

	// 获取该消息体重 Agent 的相关信息
	agentInfo, err := c.buildAgentInfo(ctx, agentIDs)
	if err != nil {
		return nil, err
	}
	// 获取历史消息体
	resp := c.buildMessageListResponse(ctx, mListMessages, currentConversation)

	resp.ParticipantInfoMap = map[string]*message.MsgParticipantInfo{}
	for _, aOne := range agentInfo {
		resp.ParticipantInfoMap[aOne.ID] = aOne
	}
	return resp, err
}

func (c *ConversationApplicationService) buildAgentInfo(ctx context.Context, agentIDs []int64) ([]*message.MsgParticipantInfo, error) {
	var result []*message.MsgParticipantInfo
	if len(agentIDs) > 0 {
		agentInfos, err := c.appContext.SingleAgentDomainSVC.MGetSingleAgentDraft(ctx, agentIDs)
		if err != nil {
			return nil, err
		}

		result = slices.Transform(agentInfos, func(a *singleAgentEntity.SingleAgent) *message.MsgParticipantInfo {
			return &message.MsgParticipantInfo{
				ID:        strconv.FormatInt(a.AgentID, 10),
				Name:      a.Name,
				UserID:    strconv.FormatInt(a.CreatorID, 10),
				Desc:      a.Desc,
				AvatarURL: a.IconURI,
			}
		})
	}

	return result, nil
}

func (c *ConversationApplicationService) buildMessageListResponse(ctx context.Context, mListMessages *entity.ListResult, currentConversation *convEntity.Conversation) *message.GetMessageListResponse {
	var messages []*message.ChatMessage
	runToQuestionIDMap := make(map[int64]int64)

	for _, mMessage := range mListMessages.Messages {
		if mMessage.MessageType == model.MessageTypeQuestion {
			runToQuestionIDMap[mMessage.RunID] = mMessage.ID
		}
	}

	for _, mMessage := range mListMessages.Messages {
		messages = append(messages, c.buildDomainMsg2VOMessage(ctx, mMessage, runToQuestionIDMap))
	}

	resp := &message.GetMessageListResponse{
		MessageList:             messages,
		Cursor:                  strconv.FormatInt(mListMessages.PrevCursor, 10),
		NextCursor:              strconv.FormatInt(mListMessages.NextCursor, 10),
		ConversationID:          strconv.FormatInt(currentConversation.ID, 10),
		LastSectionID:           ptr.Of(strconv.FormatInt(currentConversation.SectionID, 10)),
		ConnectorConversationID: strconv.FormatInt(currentConversation.ID, 10),
	}

	if mListMessages.Direction == entity.ScrollPageDirectionPrev {
		resp.Hasmore = mListMessages.HasMore
	} else {
		resp.NextHasMore = mListMessages.HasMore
	}

	return resp
}

func (c *ConversationApplicationService) buildDomainMsg2VOMessage(ctx context.Context, dm *entity.Message, runToQuestionIDMap map[int64]int64) *message.ChatMessage {
	cm := &message.ChatMessage{
		MessageID:        strconv.FormatInt(dm.ID, 10),
		Role:             string(dm.Role),
		Type:             string(dm.MessageType),
		Content:          dm.Content,
		ContentType:      string(dm.ContentType),
		ReplyID:          "0",
		SectionID:        strconv.FormatInt(dm.SectionID, 10),
		ExtraInfo:        buildDExt2ApiExt(dm.Ext),
		ContentTime:      dm.CreatedAt,
		Status:           "available",
		Source:           0,
		ReasoningContent: ptr.Of(dm.ReasoningContent),
	}

	if dm.Status == model.MessageStatusBroken {
		cm.BrokenPos = ptr.Of(dm.Position)
	}

	// 富文本、图片、音频等链接的URL
	if dm.ContentType == model.ContentTypeMix && dm.DisplayContent != "" {
		cm.Content = c.buildParseMessageURI(ctx, dm.DisplayContent)
	}

	if dm.MessageType != model.MessageTypeQuestion {
		cm.ReplyID = strconv.FormatInt(runToQuestionIDMap[dm.RunID], 10)
		cm.SenderID = ptr.Of(strconv.FormatInt(dm.AgentID, 10))
	}
	return cm
}

func (c *ConversationApplicationService) buildParseMessageURI(ctx context.Context, msgContent string) string {

	if msgContent == "" {
		return msgContent
	}

	var mc *run.MixContentModel
	err := json.Unmarshal([]byte(msgContent), &mc)
	if err != nil {
		return msgContent
	}
	for k, item := range mc.ItemList {
		switch item.Type {
		case run.ContentTypeImage:

			url, pErr := c.appContext.ImageX.GetResourceURL(ctx, item.Image.Key)
			if pErr == nil {
				mc.ItemList[k].Image.ImageThumb.URL = url.URL
				mc.ItemList[k].Image.ImageOri.URL = url.URL
			}

		case run.ContentTypeFile, run.ContentTypeAudio, run.ContentTypeVideo:
			url, pErr := c.appContext.ImageX.GetResourceURL(ctx, item.File.FileKey)
			if pErr == nil {
				mc.ItemList[k].File.FileURL = url.URL

			}

		default:

		}
	}
	jsonMsg, err := json.Marshal(mc)
	if err != nil {
		return msgContent
	}

	return string(jsonMsg)
}

func buildDExt2ApiExt(extra map[string]string) *message.ExtraInfo {
	return &message.ExtraInfo{
		InputTokens:         extra["input_tokens"],
		OutputTokens:        extra["output_tokens"],
		Token:               extra["token"],
		PluginStatus:        extra["plugin_status"],
		TimeCost:            extra["time_cost"],
		WorkflowTokens:      extra["workflow_tokens"],
		BotState:            extra["bot_state"],
		PluginRequest:       extra["plugin_request"],
		ToolName:            extra["tool_name"],
		Plugin:              extra["plugin"],
		MockHitInfo:         extra["mock_hit_info"],
		MessageTitle:        extra["message_title"],
		StreamPluginRunning: extra["stream_plugin_running"],
		ExecuteDisplayName:  extra["execute_display_name"],
		TaskType:            extra["task_type"],
		ReferFormat:         extra["refer_format"],
	}
}

func (c *ConversationApplicationService) getCurrentConversation(ctx context.Context, userID int64, agentID int64, scene common.Scene) (*convEntity.Conversation, bool, error) {
	var currentConversation *convEntity.Conversation
	var isNewCreate bool

	currentConversation, err := c.ConversationDomainSVC.GetCurrentConversation(ctx, &convEntity.GetCurrent{
		UserID:  userID,
		Scene:   scene,
		AgentID: agentID,
	})
	if err != nil {
		return nil, isNewCreate, err
	}

	// 新会话
	if currentConversation == nil {
		// 创建一个新的会话
		ccNew, err := c.ConversationDomainSVC.Create(ctx, &convEntity.CreateMeta{
			AgentID: agentID,
			UserID:  userID,
			Scene:   scene,
		})
		if err != nil {
			return nil, isNewCreate, err
		}
		if ccNew == nil {
			return nil, isNewCreate, errorx.New(errno.ErrConversationNotFound)
		}
		isNewCreate = true
		currentConversation = ccNew
	}

	return currentConversation, isNewCreate, nil
}

func loadDirectionToScrollDirection(direction *message.LoadDirection) entity.ScrollPageDirection {
	if direction != nil && *direction == message.LoadDirection_Next {
		return entity.ScrollPageDirectionNext
	}
	return entity.ScrollPageDirectionPrev
}
