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
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
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
			ConversationID: conv.Int64ToStr(currentConversation.ID),
			LastSectionID:  ptr.Of(conv.Int64ToStr(currentConversation.SectionID)),
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
				ID:        conv.Int64ToStr(a.AgentID),
				Name:      a.Name,
				UserID:    conv.Int64ToStr(a.CreatorID),
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
		Cursor:                  conv.Int64ToStr(mListMessages.PrevCursor),
		NextCursor:              conv.Int64ToStr(mListMessages.NextCursor),
		ConversationID:          conv.Int64ToStr(currentConversation.ID),
		LastSectionID:           ptr.Of(conv.Int64ToStr(currentConversation.SectionID)),
		ConnectorConversationID: conv.Int64ToStr(currentConversation.ID),
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
		MessageID:        conv.Int64ToStr(dm.ID),
		Role:             string(dm.Role),
		Type:             string(dm.MessageType),
		Content:          dm.Content,
		ContentType:      string(dm.ContentType),
		ReplyID:          "0",
		SectionID:        conv.Int64ToStr(dm.SectionID),
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
		cm.ReplyID = conv.Int64ToStr(runToQuestionIDMap[dm.RunID])
		cm.SenderID = ptr.Of(conv.Int64ToStr(dm.AgentID))
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

func (c *ConversationApplicationService) DeleteMessage(ctx context.Context, mr *message.DeleteMessageRequest) (*message.DeleteMessageResponse, error) {
	resp := new(message.DeleteMessageResponse)
	messageInfo, err := c.MessageDomainSVC.GetByID(ctx, mr.MessageID)
	if err != nil {
		return resp, err
	}
	if messageInfo == nil {
		return resp, errorx.New(errno.ErrConversationMessageNotFound)
	}

	userID := ctxutil.GetUIDFromCtx(ctx)
	if messageInfo.UserID != conv.Int64ToStr(*userID) {
		return resp, errorx.New(errno.ErrConversationPermissionCode, errorx.KV("msg", "permission denied"))
	}

	err = c.AgentRunDomainSVC.Delete(ctx, []int64{messageInfo.RunID})
	if err != nil {
		return resp, err
	}

	err = c.MessageDomainSVC.Delete(ctx, &entity.DeleteMeta{
		RunIDs: []int64{messageInfo.RunID},
	})
	if err != nil {
		return resp, nil
	}

	return resp, nil
}

func (c *ConversationApplicationService) BreakMessage(ctx context.Context, mr *message.BreakMessageRequest) (*message.BreakMessageResponse, error) {
	resp := new(message.BreakMessageResponse)
	messageInfo, err := c.MessageDomainSVC.GetByID(ctx, mr.GetAnswerMessageID())
	if err != nil {
		return resp, err
	}
	if messageInfo == nil {
		return resp, errorx.New(errno.ErrConversationMessageNotFound)
	}

	userID := ctxutil.GetUIDFromCtx(ctx)
	if messageInfo.UserID != conv.Int64ToStr(*userID) {
		return resp, errorx.New(errno.ErrConversationPermissionCode, errorx.KV("msg", "permission denied"))
	}

	if messageInfo.ConversationID != mr.ConversationID {
		return resp, errorx.New(errno.ErrConversationPermissionCode, errorx.KV("msg", "conversation not match"))
	}

	err = c.MessageDomainSVC.Broken(ctx, &entity.BrokenMeta{
		ID:       *mr.AnswerMessageID,
		Position: mr.BrokenPos,
	})
	if err != nil {
		return resp, err
	}

	return resp, nil
}
