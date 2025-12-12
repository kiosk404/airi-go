package application

import (
	"context"
	"errors"
	"io"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/api/model/conversation/message"
	"github.com/kiosk404/airi-go/backend/api/model/conversation/run"
	"github.com/kiosk404/airi-go/backend/application/ctxutil"
	singleagentEntity "github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/pkg"
	convEntity "github.com/kiosk404/airi-go/backend/modules/conversation/conversation/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/pkg/errno"
	crossDomainMessage "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/message/model"
	msgEntity "github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	sseImpl "github.com/kiosk404/airi-go/backend/pkg/http/sse"
	"github.com/kiosk404/airi-go/backend/pkg/json"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

func (c *ConversationApplicationService) Run(ctx context.Context, sseSender *sseImpl.SSenderImpl, ar *run.AgentRunRequest) error {
	agentInfo, caErr := c.checkAgent(ctx, ar)
	if caErr != nil {
		logs.ErrorX(pkg.ModelName, "checkAgent err:%v", caErr)
		return caErr
	}
	logs.DebugX(pkg.ModelName, "agent run req id:%v, req name:%v", agentInfo.AgentID, agentInfo.Name)

	userID := ctxutil.MustGetUIDFromCtx(ctx)

	// 验证对话是否存在以及是否拥有权限访问该对话
	conversationData, ccErr := c.checkConversation(ctx, ar, userID)
	if ccErr != nil {
		logs.ErrorX(pkg.ModelName, "checkConversation err:%v", ccErr)
		return ccErr
	}

	// 处理消息的重生成逻辑
	if ar.RegenMessageID != nil && ptr.From(ar.RegenMessageID) > 0 {
		// 获取重生成消息的元数据
		msgMeta, err := c.MessageDomainSVC.GetByID(ctx, ptr.From(ar.RegenMessageID))
		if err != nil {
			return err
		}
		// 验证消息存在性并检查用户权限
		if msgMeta != nil {
			if msgMeta.UserID != conv.Int64ToStr(userID) {
				return errorx.New(errno.ErrConversationPermissionCode, errorx.KV("msg", "message not match"))
			}

			err = c.AgentRunDomainSVC.Delete(ctx, []int64{msgMeta.RunID})
			if err != nil {
				return err
			}

			delErr := c.MessageDomainSVC.Delete(ctx, &msgEntity.DeleteMeta{
				RunIDs: []int64{msgMeta.RunID},
			})
			if delErr != nil {
				return delErr
			}
		}
	}

	// 构建AgentRunMeta请求
	arr, err := c.buildAgentRunRequest(ctx, ar, userID, conversationData)
	if err != nil {
		logs.ErrorX(pkg.ModelName, "buildAgentRunRequest err:%v", err)
		return err
	}
	// 启动智能体运行
	streamer, err := c.AgentRunDomainSVC.AgentRun(ctx, arr)
	if err != nil {
		return err
	}
	// 处理流式响应
	c.pullStream(ctx, sseSender, streamer, ar)
	return nil
}

func (c *ConversationApplicationService) checkAgent(ctx context.Context, ar *run.AgentRunRequest) (*singleagentEntity.SingleAgent, error) {
	agentInfo, err := c.appContext.SingleAgentDomainSVC.GetSingleAgent(ctx, ar.BotID, "")
	if err != nil {
		return nil, err
	}

	if agentInfo == nil {
		return nil, errorx.New(errno.ErrAgentNotExists)
	}
	return agentInfo, nil
}

func (c *ConversationApplicationService) buildAgentRunRequest(ctx context.Context, ar *run.AgentRunRequest, userID int64, conversationData *convEntity.Conversation) (*entity.AgentRunMeta, error) {
	var contentType crossDomainMessage.ContentType
	contentType = crossDomainMessage.ContentTypeText

	if ptr.From(ar.ContentType) != string(crossDomainMessage.ContentTypeText) {
		contentType = crossDomainMessage.ContentTypeMix
	}

	preRetrieveTools, err := c.buildTools(ctx, ar.ToolList)
	if err != nil {
		return nil, err
	}

	arm := &entity.AgentRunMeta{
		ConversationID:   conversationData.ID,
		AgentID:          ar.BotID,
		Content:          c.buildMultiContent(ctx, ar),
		DisplayContent:   c.buildDisplayContent(ctx, ar),
		UserID:           conv.Int64ToStr(userID),
		SectionID:        conversationData.SectionID,
		PreRetrieveTools: preRetrieveTools,
		IsDraft:          ptr.From(ar.DraftMode),
		ContentType:      contentType,
		Ext:              ar.Extra,
	}
	return arm, nil
}

func (c *ConversationApplicationService) GenID(ctx context.Context) (int64, error) {
	id, err := c.appContext.IDGen.GenID(ctx)
	return id, err
}

func (c *ConversationApplicationService) buildDisplayContent(ctx context.Context, ar *run.AgentRunRequest) string {
	if *ar.ContentType == run.ContentTypeText {
		return ""
	}
	return ar.Query
}

func (c *ConversationApplicationService) checkConversation(ctx context.Context, ar *run.AgentRunRequest, userID int64) (*convEntity.Conversation, error) {
	var conversationData *convEntity.Conversation
	if ar.ConversationID > 0 {
		realCurrCon, err := c.ConversationDomainSVC.GetCurrentConversation(ctx, &convEntity.GetCurrent{
			UserID:  userID,
			AgentID: ar.BotID,
			Scene:   ptr.From(ar.Scene),
		})
		logs.InfoX(pkg.ModelName, "conversation data: %v", conv.DebugJsonToStr(realCurrCon))
		if err != nil {
			return nil, err
		}
		if realCurrCon != nil {
			conversationData = realCurrCon
		}
	}
	if ar.ConversationID == 0 || conversationData == nil {

		conData, err := c.ConversationDomainSVC.Create(ctx, &convEntity.CreateMeta{
			AgentID: ar.BotID,
			UserID:  userID,
			Scene:   ptr.From(ar.Scene),
		})
		if err != nil {
			return nil, err
		}
		logs.InfoX(pkg.ModelName, "conversatioin create data:%v", conv.DebugJsonToStr(conData))
		conversationData = conData

		ar.ConversationID = conversationData.ID
	}

	if conversationData.CreatorID != userID {
		return nil, errorx.New(errno.ErrConversationPermissionCode, errorx.KV("msg", "conversation not match"))
	}

	return conversationData, nil
}

func (c *ConversationApplicationService) buildTools(ctx context.Context, tools []*run.Tool) ([]*entity.Tool, error) {
	// todo 实现 快捷指令.
	return nil, nil
}

func (c *ConversationApplicationService) buildMultiContent(ctx context.Context, ar *run.AgentRunRequest) []*crossDomainMessage.InputMetaData {
	var multiContents []*crossDomainMessage.InputMetaData

	switch *ar.ContentType {
	case run.ContentTypeText:
		multiContents = append(multiContents, &crossDomainMessage.InputMetaData{
			Type: crossDomainMessage.InputTypeText,
			Text: ar.Query,
		})
	case run.ContentTypeImage, run.ContentTypeFile, run.ContentTypeAudio, run.ContentTypeVideo:
		var mc *run.MixContentModel

		err := json.Unmarshal([]byte(ar.Query), &mc)
		if err != nil {
			multiContents = append(multiContents, &crossDomainMessage.InputMetaData{
				Type: crossDomainMessage.InputTypeText,
				Text: ar.Query,
			})
			return multiContents
		}
		mcContent, newItemList := c.parseMultiContent(ctx, mc.ItemList)

		multiContents = append(multiContents, mcContent...)

		mc.ItemList = newItemList
		mcByte, err := json.Marshal(mc)
		if err == nil {
			ar.Query = string(mcByte)
		}
	default:
		logs.ErrorX(pkg.ModelName, "unknown content type:%v", *ar.ContentType)
	}
	return multiContents
}

func (c *ConversationApplicationService) parseMultiContent(ctx context.Context, mc []*run.Item) (multiContents []*crossDomainMessage.InputMetaData, mcNew []*run.Item) {
	for index, item := range mc {
		switch item.Type {
		case run.ContentTypeText:
			multiContents = append(multiContents, &crossDomainMessage.InputMetaData{
				Type: crossDomainMessage.InputTypeText,
				Text: item.Text,
			})
		case run.ContentTypeImage:

			resourceUrl, err := c.getUrlByUri(ctx, item.Image.Key)
			if err != nil {
				logs.ErrorX(pkg.ModelName, "failed to unescape resource url, err is %v", err)
				continue
			}

			if resourceUrl == "" {
				logs.ErrorX(pkg.ModelName, "failed to unescape resource url, uri is %v", item.Image.Key)
				continue
			}

			mc[index].Image.ImageThumb.URL = resourceUrl
			mc[index].Image.ImageOri.URL = resourceUrl

			multiContents = append(multiContents, &crossDomainMessage.InputMetaData{
				Type: crossDomainMessage.InputTypeImage,
				FileData: []*crossDomainMessage.FileData{
					{
						Url: resourceUrl,
						URI: item.Image.Key,
					},
				},
			})
		case run.ContentTypeFile, run.ContentTypeAudio, run.ContentTypeVideo:

			resourceUrl, err := c.getUrlByUri(ctx, item.File.FileKey)
			if err != nil {
				continue
			}

			mc[index].File.FileURL = resourceUrl

			multiContents = append(multiContents, &crossDomainMessage.InputMetaData{
				Type: c.getType(item.File.FileType),
				FileData: []*crossDomainMessage.FileData{
					{
						Url: resourceUrl,
						URI: item.File.FileKey,
					},
				},
			})
		}
	}

	return multiContents, mc
}

func (c *ConversationApplicationService) getUrlByUri(ctx context.Context, uri string) (string, error) {
	url, err := c.appContext.ImageX.GetResourceURL(ctx, uri)
	if err != nil {
		return "", err
	}
	return url.URL, nil
}

func (c *ConversationApplicationService) getType(fileType string) crossDomainMessage.InputType {
	switch fileType {
	case string(crossDomainMessage.InputTypeAudio):
		return crossDomainMessage.InputTypeAudio
	case string(crossDomainMessage.InputTypeVideo):
		return crossDomainMessage.InputTypeVideo
	default:
		return crossDomainMessage.InputTypeFile
	}
}

func (c *ConversationApplicationService) pullStream(ctx context.Context, sseSender *sseImpl.SSenderImpl, arStream *schema.StreamReader[*entity.AgentRunResponse], req *run.AgentRunRequest) {
	var ackMessageInfo *entity.ChunkMessageItem
	for {
		chunk, recvErr := arStream.Recv()
		if recvErr != nil {
			if errors.Is(recvErr, io.EOF) {
				return
			}
			sseSender.Send(ctx, buildErrorEvent(errno.ErrConversationAgentRunError, recvErr.Error()))
			return
		}
		switch chunk.Event {
		case entity.RunEventCreated, entity.RunEventInProgress, entity.RunEventCompleted:
		case entity.RunEventError:
			id, err := c.GenID(ctx)
			if err != nil {
				sseSender.Send(ctx, buildErrorEvent(errno.ErrConversationAgentRunError, err.Error()))
			} else {
				sseSender.Send(ctx, buildMessageChunkEvent(run.RunEventMessage, buildErrMsg(ackMessageInfo, chunk.Error, id)))
			}
		case entity.RunEventStreamDone:
			sseSender.Send(ctx, buildDoneEvent(run.RunEventDone))
		case entity.RunEventAck:
			ackMessageInfo = chunk.ChunkMessageItem
			sseSender.Send(ctx, buildMessageChunkEvent(run.RunEventMessage, buildARSM2Message(chunk, req)))
		case entity.RunEventMessageDelta, entity.RunEventMessageCompleted:
			sseSender.Send(ctx, buildMessageChunkEvent(run.RunEventMessage, buildARSM2Message(chunk, req)))
		default:
			logs.ErrorX(pkg.ModelName, "unknown handler event:%v", chunk.Event)
		}
	}
}

func buildARSM2Message(chunk *entity.AgentRunResponse, req *run.AgentRunRequest) []byte {
	chunkMessageItem := chunk.ChunkMessageItem

	chunkMessage := &run.RunStreamResponse{
		ConversationID: conv.Int64ToStr(chunkMessageItem.ConversationID),
		IsFinish:       ptr.Of(chunk.ChunkMessageItem.IsFinish),
		Message: &message.ChatMessage{
			Role:        string(chunkMessageItem.Role),
			ContentType: string(chunkMessageItem.ContentType),
			MessageID:   conv.Int64ToStr(chunkMessageItem.ID),
			SectionID:   conv.Int64ToStr(chunkMessageItem.SectionID),
			ContentTime: chunkMessageItem.CreatedAt,
			ExtraInfo:   buildExt(chunkMessageItem.Ext),
			ReplyID:     conv.Int64ToStr(chunkMessageItem.ReplyID),

			Status:           "",
			Type:             string(chunkMessageItem.MessageType),
			Content:          chunkMessageItem.Content,
			ReasoningContent: chunkMessageItem.ReasoningContent,
			RequiredAction:   chunkMessageItem.RequiredAction,
		},
		Index: int32(chunkMessageItem.Index),
		SeqID: int32(chunkMessageItem.SeqID),
	}
	if chunkMessageItem.MessageType == crossDomainMessage.MessageTypeAck {
		chunkMessage.Message.Content = req.GetQuery()
		chunkMessage.Message.ContentType = req.GetContentType()
		chunkMessage.Message.ExtraInfo = &message.ExtraInfo{
			LocalMessageID: req.GetLocalMessageID(),
		}
	} else {
		chunkMessage.Message.ExtraInfo = buildExt(chunkMessageItem.Ext)
		chunkMessage.Message.SenderID = ptr.Of(conv.Int64ToStr(chunkMessageItem.AgentID))
		chunkMessage.Message.Content = chunkMessageItem.Content

		if chunkMessageItem.MessageType == crossDomainMessage.MessageTypeKnowledge {
			chunkMessage.Message.Type = string(crossDomainMessage.MessageTypeVerbose)
		}
	}

	if chunk.ChunkMessageItem.IsFinish && chunkMessageItem.MessageType == crossDomainMessage.MessageTypeAnswer {
		chunkMessage.Message.Content = ""
		chunkMessage.Message.ReasoningContent = ptr.Of("")
	}

	mCM, _ := json.Marshal(chunkMessage)
	return mCM
}

func buildExt(extra map[string]string) *message.ExtraInfo {
	if extra == nil {
		return nil
	}

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

func buildErrMsg(ackChunk *entity.ChunkMessageItem, err *entity.RunError, id int64) []byte {
	chunkMessage := &run.RunStreamResponse{
		IsFinish:       ptr.Of(true),
		ConversationID: conv.Int64ToStr(ackChunk.ConversationID),
		Message: &message.ChatMessage{
			Role:        string(schema.Assistant),
			ContentType: string(crossDomainMessage.ContentTypeText),
			Type:        string(crossDomainMessage.MessageTypeAnswer),
			MessageID:   conv.Int64ToStr(id),
			SectionID:   conv.Int64ToStr(ackChunk.SectionID),
			ReplyID:     conv.Int64ToStr(ackChunk.ReplyID),
			Content:     "Something error:" + err.Msg,
			ExtraInfo:   &message.ExtraInfo{},
		},
	}

	mCM, _ := json.Marshal(chunkMessage)
	return mCM
}
