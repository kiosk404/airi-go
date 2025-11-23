package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/pkg/errno"
	agentrun "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/agentrun/model"
	crossmessage "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/message"
	message "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/message/model"
	msgEntity "github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/json"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

type Event struct {
}

func NewMessageEvent() *Event {
	return &Event{}
}

func (e *Event) buildMessageEvent(runEvent entity.RunEvent, chunkMsgItem *entity.ChunkMessageItem) *entity.AgentRunResponse {
	return &entity.AgentRunResponse{
		Event:            runEvent,
		ChunkMessageItem: chunkMsgItem,
	}
}

func (e *Event) buildRunEvent(runEvent entity.RunEvent, chunkRunItem *entity.ChunkRunItem) *entity.AgentRunResponse {
	return &entity.AgentRunResponse{
		Event:        runEvent,
		ChunkRunItem: chunkRunItem,
	}
}

func (e *Event) buildErrEvent(runEvent entity.RunEvent, err *entity.RunError) *entity.AgentRunResponse {
	return &entity.AgentRunResponse{
		Event: runEvent,
		Error: err,
	}
}

func (e *Event) buildStreamDoneEvent() *entity.AgentRunResponse {

	return &entity.AgentRunResponse{
		Event: entity.RunEventStreamDone,
	}
}

func (e *Event) SendRunEvent(runEvent entity.RunEvent, runItem *entity.ChunkRunItem, sw *schema.StreamWriter[*entity.AgentRunResponse]) {
	resp := e.buildRunEvent(runEvent, runItem)
	sw.Send(resp, nil)
}

func (e *Event) SendMsgEvent(runEvent entity.RunEvent, messageItem *entity.ChunkMessageItem, sw *schema.StreamWriter[*entity.AgentRunResponse]) {
	resp := e.buildMessageEvent(runEvent, messageItem)
	sw.Send(resp, nil)
}

func (e *Event) SendErrEvent(runEvent entity.RunEvent, sw *schema.StreamWriter[*entity.AgentRunResponse], err *entity.RunError) {
	resp := e.buildErrEvent(runEvent, err)
	sw.Send(resp, nil)
}

func (e *Event) SendStreamDoneEvent(sw *schema.StreamWriter[*entity.AgentRunResponse]) {
	resp := e.buildStreamDoneEvent()
	sw.Send(resp, nil)
}

type MesssageEventHanlder struct {
	messageEvent *Event
	sw           *schema.StreamWriter[*entity.AgentRunResponse]
}

func (mh *MesssageEventHanlder) handlerErr(_ context.Context, err error) {

	var errMsg string
	var statusErr errorx.StatusError
	if errors.As(err, &statusErr) {
		errMsg = statusErr.Msg()
	} else {
		if strings.ToLower(os.Getenv(consts.RunMode)) != "debug" {
			errMsg = "Internal Server Error"
		} else {
			errMsg = errorx.ErrorWithoutStack(err)
		}
	}

	mh.messageEvent.SendErrEvent(entity.RunEventError, mh.sw, &entity.RunError{
		Code: errno.ErrAgentRun,
		Msg:  errMsg,
	})
}

func (mh *MesssageEventHanlder) handlerAckMessage(_ context.Context, input *msgEntity.Message) error {
	sendMsg := &entity.ChunkMessageItem{
		ID:             input.ID,
		ConversationID: input.ConversationID,
		SectionID:      input.SectionID,
		AgentID:        input.AgentID,
		Role:           entity.RoleType(input.Role),
		MessageType:    message.MessageTypeAck,
		ReplyID:        input.ID,
		Content:        input.Content,
		ContentType:    message.ContentTypeText,
		IsFinish:       true,
	}

	mh.messageEvent.SendMsgEvent(entity.RunEventAck, sendMsg, mh.sw)

	return nil
}

func (mh *MesssageEventHanlder) handlerFunctionCall(ctx context.Context, chunk *entity.AgentRespEvent, rtDependence *AgentRuntime) error {
	cm := buildAgentMessage2Create(ctx, chunk, message.MessageTypeFunctionCall, rtDependence)

	cmData, err := crossmessage.DefaultSVC().Create(ctx, cm)
	if err != nil {
		return err
	}

	sendMsg := buildSendMsg(ctx, cmData, true, rtDependence)

	mh.messageEvent.SendMsgEvent(entity.RunEventMessageCompleted, sendMsg, mh.sw)
	return nil
}

func (mh *MesssageEventHanlder) handlerTooResponse(ctx context.Context, chunk *entity.AgentRespEvent, rtDependence *AgentRuntime, preToolResponseMsg *msgEntity.Message, toolResponseMsgContent string) error {

	cm := buildAgentMessage2Create(ctx, chunk, message.MessageTypeToolResponse, rtDependence)

	var cmData *message.Message
	var err error

	if preToolResponseMsg != nil {
		cm.ID = preToolResponseMsg.ID
		cm.CreatedAt = preToolResponseMsg.CreatedAt
		cm.UpdatedAt = preToolResponseMsg.UpdatedAt
		if len(toolResponseMsgContent) > 0 {
			cm.Content = toolResponseMsgContent + "\n" + cm.Content
		}
	}

	cmData, err = crossmessage.DefaultSVC().Create(ctx, cm)
	if err != nil {
		return err
	}

	sendMsg := buildSendMsg(ctx, cmData, true, rtDependence)

	mh.messageEvent.SendMsgEvent(entity.RunEventMessageCompleted, sendMsg, mh.sw)

	return nil
}

func (mh *MesssageEventHanlder) handlerSuggest(ctx context.Context, chunk *entity.AgentRespEvent, rtDependence *AgentRuntime) error {
	cm := buildAgentMessage2Create(ctx, chunk, message.MessageTypeFlowUp, rtDependence)

	cmData, err := crossmessage.DefaultSVC().Create(ctx, cm)
	if err != nil {
		return err
	}

	sendMsg := buildSendMsg(ctx, cmData, true, rtDependence)

	mh.messageEvent.SendMsgEvent(entity.RunEventMessageCompleted, sendMsg, mh.sw)

	return nil
}

func (mh *MesssageEventHanlder) handlerKnowledge(ctx context.Context, chunk *entity.AgentRespEvent, rtDependence *AgentRuntime) error {
	cm := buildAgentMessage2Create(ctx, chunk, message.MessageTypeKnowledge, rtDependence)
	cmData, err := crossmessage.DefaultSVC().Create(ctx, cm)
	if err != nil {
		return err
	}

	sendMsg := buildSendMsg(ctx, cmData, true, rtDependence)

	mh.messageEvent.SendMsgEvent(entity.RunEventMessageCompleted, sendMsg, mh.sw)
	return nil
}

func (mh *MesssageEventHanlder) handlerAnswer(ctx context.Context, msg *entity.ChunkMessageItem, usage *msgEntity.UsageExt, rtDependence *AgentRuntime, preAnswerMsg *msgEntity.Message) error {

	if len(msg.Content) == 0 && len(ptr.From(msg.ReasoningContent)) == 0 {
		return nil
	}

	msg.IsFinish = true

	if msg.Ext == nil {
		msg.Ext = map[string]string{}
	}
	if usage != nil {
		msg.Ext[string(msgEntity.MessageExtKeyToken)] = strconv.FormatInt(usage.TotalCount, 10)
		msg.Ext[string(msgEntity.MessageExtKeyInputTokens)] = strconv.FormatInt(usage.InputTokens, 10)
		msg.Ext[string(msgEntity.MessageExtKeyOutputTokens)] = strconv.FormatInt(usage.OutputTokens, 10)

		rtDependence.Usage = &agentrun.Usage{
			LlmPromptTokens:     usage.InputTokens,
			LlmCompletionTokens: usage.OutputTokens,
			LlmTotalTokens:      usage.TotalCount,
		}
	}

	if _, ok := msg.Ext[string(msgEntity.MessageExtKeyTimeCost)]; !ok {
		msg.Ext[string(msgEntity.MessageExtKeyTimeCost)] = fmt.Sprintf("%.1f", float64(time.Since(rtDependence.GetStartTime()).Milliseconds())/1000.00)
	}

	buildModelContent := &schema.Message{
		Role:    schema.Assistant,
		Content: msg.Content,
	}

	mc, err := json.Marshal(buildModelContent)
	if err != nil {
		return err
	}
	preAnswerMsg.Content = msg.Content
	preAnswerMsg.ReasoningContent = ptr.From(msg.ReasoningContent)
	preAnswerMsg.Ext = msg.Ext
	preAnswerMsg.ContentType = msg.ContentType
	preAnswerMsg.ModelContent = string(mc)
	preAnswerMsg.CreatedAt = 0
	preAnswerMsg.UpdatedAt = 0

	_, err = crossmessage.DefaultSVC().Create(ctx, preAnswerMsg)
	if err != nil {
		return err
	}
	mh.messageEvent.SendMsgEvent(entity.RunEventMessageCompleted, msg, mh.sw)

	return nil
}

func (mh *MesssageEventHanlder) handlerFinalAnswerFinish(ctx context.Context, rtDependence *AgentRuntime) error {
	cm := buildAgentMessage2Create(ctx, nil, message.MessageTypeVerbose, rtDependence)
	cmData, err := crossmessage.DefaultSVC().Create(ctx, cm)
	if err != nil {
		return err
	}

	sendMsg := buildSendMsg(ctx, cmData, true, rtDependence)

	mh.messageEvent.SendMsgEvent(entity.RunEventMessageCompleted, sendMsg, mh.sw)
	return nil
}

func (mh *MesssageEventHanlder) handlerInterruptVerbose(ctx context.Context, chunk *entity.AgentRespEvent, rtDependence *AgentRuntime) error {
	cm := buildAgentMessage2Create(ctx, chunk, message.MessageTypeInterrupt, rtDependence)
	cmData, err := crossmessage.DefaultSVC().Create(ctx, cm)
	if err != nil {
		return err
	}

	sendMsg := buildSendMsg(ctx, cmData, true, rtDependence)

	mh.messageEvent.SendMsgEvent(entity.RunEventMessageCompleted, sendMsg, mh.sw)
	return nil
}

func (mh *MesssageEventHanlder) handlerWfUsage(ctx context.Context, msg *entity.ChunkMessageItem, usage *msgEntity.UsageExt) error {

	if msg.Ext == nil {
		msg.Ext = map[string]string{}
	}
	if usage != nil {
		msg.Ext[string(msgEntity.MessageExtKeyToken)] = strconv.FormatInt(usage.TotalCount, 10)
		msg.Ext[string(msgEntity.MessageExtKeyInputTokens)] = strconv.FormatInt(usage.InputTokens, 10)
		msg.Ext[string(msgEntity.MessageExtKeyOutputTokens)] = strconv.FormatInt(usage.OutputTokens, 10)
	}

	_, err := crossmessage.DefaultSVC().Edit(ctx, &msgEntity.Message{
		ID:  msg.ID,
		Ext: msg.Ext,
	})
	if err != nil {
		return err
	}

	mh.messageEvent.SendMsgEvent(entity.RunEventMessageCompleted, msg, mh.sw)
	return nil
}

func (mh *MesssageEventHanlder) HandlerInput(ctx context.Context, rtDependence *AgentRuntime) (*msgEntity.Message, error) {
	msgMeta := buildAgentMessage2Create(ctx, nil, message.MessageTypeQuestion, rtDependence)

	cm, err := crossmessage.DefaultSVC().Create(ctx, msgMeta)
	if err != nil {
		return nil, err
	}

	ackErr := mh.handlerAckMessage(ctx, cm)
	if ackErr != nil {
		return msgMeta, ackErr
	}
	return cm, nil
}
