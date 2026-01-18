package runtime

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/infra/contract/imagex"
	crossagent "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/agent"
	singleagent "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/agent/model"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/pkg"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/pkg/errno"
	agentrun "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/agentrun/model"
	crossmessage "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/message"
	message "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/message/model"
	msgEntity "github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/pkg/utils/safego"
	"github.com/mohae/deepcopy"
)

// AgentStreamExecute 执行单 Agent 的流式对话。
//
// 该方法是 Agent 执行的核心入口，采用生产者-消费者模式：
//   - pull goroutine (生产者): 从 Agent 执行引擎接收事件流
//   - push goroutine (消费者): 处理事件并推送给客户端
//
// 执行流程：
//  1. 构建 AgentRuntime 配置（包含版本、用户、历史消息等）
//  2. 调用跨域 Agent 服务执行流式推理
//  3. 启动两个并发 goroutine 处理事件流
//  4. 等待所有 goroutine 完成
//
// 架构图：
//
//	┌─────────────────┐     ┌─────────────┐     ┌─────────────────┐
//	│  Agent Engine   │────▶│  mainChan   │────▶│  Client (SSE)   │
//	│  (StreamExecute)│     │  (buffer:100)│     │  (WebSocket)    │
//	└─────────────────┘     └─────────────┘     └─────────────────┘
//	       │                       │                    │
//	       │    pull() goroutine   │   push() goroutine │
//	       └───────────────────────┴────────────────────┘
//
// 参数：
//   - ctx: 上下文，用于控制超时和取消
//   - imagex: 图片处理服务，用于处理消息中的图片
//
// 返回值：
//   - err: 执行过程中的错误
func (art *AgentRuntime) AgentStreamExecute(ctx context.Context, imagex imagex.ImageX) (err error) {
	mainChan := make(chan *entity.AgentRespEvent, 100)

	ar := &crossagent.AgentRuntime{
		AgentVersion:     art.GetRunMeta().Version,
		AgentID:          art.GetRunMeta().AgentID,
		IsDraft:          art.GetRunMeta().IsDraft,
		UserID:           art.GetRunMeta().UserID,
		PreRetrieveTools: art.GetRunMeta().PreRetrieveTools, // 预检工具
		// 将用户输入消息转为 schema.Message 类型 格式
		Input: transMessageToSchemaMessage(ctx, []*msgEntity.Message{art.GetInput()}, imagex)[0],
		// 将历史消息转为 schema.Message 类型 格式
		HistoryMsg: transMessageToSchemaMessage(ctx, historyPairs(art.GetHistory()), imagex),
		// 解析恢复信息（用于断点续传场景）
		ResumeInfo: parseResumeInfo(ctx, art.GetHistory()),
	}

	streamer, err := crossagent.DefaultSVC().StreamExecute(ctx, ar)
	if err != nil {
		return errors.New(errorx.ErrorWithoutStack(err))
	}

	var wg sync.WaitGroup
	wg.Add(2)
	safego.Go(ctx, func() {
		defer wg.Done()
		// 从 LLM 接收事件流
		art.pull(ctx, mainChan, streamer)
	})

	safego.Go(ctx, func() {
		defer wg.Done()
		// 推送消息给客户端
		art.push(ctx, mainChan)
	})

	wg.Wait()

	return err
}

// push 是消费者方法，负责从 mainChan 接收事件并处理。
//
// 该方法是事件处理的核心，根据不同的事件类型执行相应的处理逻辑：
//   - MessageTypeFunctionCall: 处理工具/函数调用事件
//   - MessageTypeToolResponse: 处理工具执行结果
//   - MessageTypeKnowledge: 处理知识库检索结果
//   - MessageTypeToolMidAnswer: 处理工具执行过程中的中间回答
//   - MessageTypeToolAsAnswer: 处理工具直接作为最终回答的情况
//   - MessageTypeAnswer: 处理模型的正常回答（包含推理内容和最终内容）
//   - MessageTypeFlowUp: 处理后续问题建议
//   - MessageTypeInterrupt: 处理用户打断（待实现）
//
// 状态管理：
//   - reasoningContent: 累积模型的推理过程内容（思维链）
//   - firstAnswerMsg: 记录第一条回答消息，用于关联推理内容
//   - isSendFinishAnswer: 标记是否已发送完成信号
//
// 参数：
//   - ctx: 上下文
//   - mainChan: 事件接收通道
func (art *AgentRuntime) push(ctx context.Context, mainChan chan *entity.AgentRespEvent) {

	mh := &MesssageEventHanlder{
		sw:           art.SW, // StreamWriter, 用于向客户端推送消息
		messageEvent: art.MessageEvent,
	}

	var err error
	defer func() {
		if err != nil {
			logs.ErrorX(pkg.ModelName, "run.push error: %v", err)
			mh.handlerErr(ctx, err)
		}
	}()

	// ====================== 状态管理 ======================
	// reasoningContent 用于累积模型的推理过程内容（思维链/CoT）

	// 思考模型（如 Deepseek R1）会在回答前输出推理内容
	reasoningContent := bytes.NewBuffer([]byte{})

	var firstAnswerMsg *msgEntity.Message // 用于将累计的推理内容关联到正确的消息上
	var reasoningMsg *msgEntity.Message   // 记录推理内容对应的消息实体

	isSendFinishAnswer := false                         // 标记是否已发送完成信号，防止重复发送
	var preToolResponseMsg *msgEntity.Message           // 用于在工具执行完成后关联结果
	toolResponseMsgContent := bytes.NewBuffer([]byte{}) // 累计工具的执行结果
	for {
		chunk, ok := <-mainChan
		if !ok || chunk == nil {
			return
		}

		// 错误处理
		if chunk.Err != nil {
			// EOF 表示流式推理结束
			if errors.Is(chunk.Err, io.EOF) {
				if !isSendFinishAnswer {
					isSendFinishAnswer = true

					// 如果有累计的推理内容，则保存到 firstAnswerMsg 对应的消息中
					if firstAnswerMsg != nil && len(reasoningContent.String()) > 0 {
						art.saveReasoningContent(ctx, firstAnswerMsg, reasoningContent.String())
						reasoningContent.Reset()
					}

					// 发送最终完成信号
					finishErr := mh.handlerFinalAnswerFinish(ctx, art)
					if finishErr != nil {
						err = finishErr
						return
					}
				}
				return
			}
			// 其他错误，记录错误信息
			mh.handlerErr(ctx, chunk.Err)
			return
		}

		switch chunk.EventType {
		// 当模型决定调用工具时触发
		case message.MessageTypeFunctionCall:
			if chunk.FuncCall != nil && chunk.FuncCall.ResponseMeta != nil {
				if usage := handlerUsage(chunk.FuncCall.ResponseMeta); usage != nil {
					art.SetUsage(&agentrun.Usage{
						LlmPromptTokens:     usage.InputTokens,
						LlmCompletionTokens: usage.OutputTokens,
						LlmTotalTokens:      usage.TotalCount,
					})
				}
			}
			err = mh.handlerFunctionCall(ctx, chunk, art)
			if err != nil {
				return
			}

			if preToolResponseMsg == nil {
				var cErr error
				preToolResponseMsg, cErr = preCreateAnswer(ctx, art)
				if cErr != nil {
					err = cErr
					return
				}
			}
		// 工具执行完成后响应
		case message.MessageTypeToolResponse:
			err = mh.handlerTooResponse(ctx, chunk, art, preToolResponseMsg, toolResponseMsgContent.String())
			if err != nil {
				return
			}
			preToolResponseMsg = nil // reset
		// 当 Agent 决定调用知识库时触发
		case message.MessageTypeKnowledge:
			err = mh.handlerKnowledge(ctx, chunk, art)
			if err != nil {
				return
			}
		// 执行过程中产生的中间输出（如工作流节点输出）
		case message.MessageTypeToolMidAnswer:
			fullMidAnswerContent := bytes.NewBuffer([]byte{})
			var usage *msgEntity.UsageExt
			toolMidAnswerMsg, cErr := preCreateAnswer(ctx, art)

			if cErr != nil {
				err = cErr
				return
			}

			var preMsgIsFinish = false
			for {
				streamMsg, receErr := chunk.ToolMidAnswer.Recv()
				if receErr != nil {
					if errors.Is(receErr, io.EOF) {
						break
					}
					err = receErr
					return
				}
				if preMsgIsFinish {
					toolMidAnswerMsg, cErr = preCreateAnswer(ctx, art)
					if cErr != nil {
						err = cErr
						return
					}
					preMsgIsFinish = false
				}
				if streamMsg == nil {
					continue
				}
				if firstAnswerMsg == nil && len(streamMsg.Content) > 0 {
					if reasoningMsg != nil {
						toolMidAnswerMsg = deepcopy.Copy(reasoningMsg).(*msgEntity.Message)
					}
					firstAnswerMsg = deepcopy.Copy(toolMidAnswerMsg).(*msgEntity.Message)
				}

				if streamMsg.Extra != nil {
					if val, ok := streamMsg.Extra["workflow_node_name"]; ok && val != nil {
						toolMidAnswerMsg.Ext["message_title"] = val.(string)
					}
				}

				sendMidAnswerMsg := buildSendMsg(ctx, toolMidAnswerMsg, false, art)
				sendMidAnswerMsg.Content = streamMsg.Content
				toolResponseMsgContent.WriteString(streamMsg.Content)
				fullMidAnswerContent.WriteString(streamMsg.Content)

				art.MessageEvent.SendMsgEvent(entity.RunEventMessageDelta, sendMidAnswerMsg, art.SW)

				if streamMsg != nil && streamMsg.ResponseMeta != nil {
					usage = handlerUsage(streamMsg.ResponseMeta)
				}

				if streamMsg.Extra["is_finish"] == true {
					preMsgIsFinish = true
					sendMidAnswerMsg := buildSendMsg(ctx, toolMidAnswerMsg, false, art)
					sendMidAnswerMsg.Content = fullMidAnswerContent.String()
					fullMidAnswerContent.Reset()
					hfErr := mh.handlerAnswer(ctx, sendMidAnswerMsg, usage, art, toolMidAnswerMsg)
					if hfErr != nil {
						err = hfErr
						return
					}
				}
			}
		// 某些工具的输出直接作为最终回答返回给用户
		case message.MessageTypeToolAsAnswer:
			var usage *msgEntity.UsageExt
			fullContent := bytes.NewBuffer([]byte{})
			toolAsAnswerMsg, cErr := preCreateAnswer(ctx, art)
			if cErr != nil {
				err = cErr
				return
			}
			if firstAnswerMsg == nil {
				firstAnswerMsg = toolAsAnswerMsg
			}

			for {
				streamMsg, receErr := chunk.ToolAsAnswer.Recv()
				if receErr != nil {
					if errors.Is(receErr, io.EOF) {

						answer := buildSendMsg(ctx, toolAsAnswerMsg, false, art)
						answer.Content = fullContent.String()
						hfErr := mh.handlerAnswer(ctx, answer, usage, art, toolAsAnswerMsg)
						if hfErr != nil {
							err = hfErr
							return
						}
						break
					}
					err = receErr
					return
				}

				if streamMsg != nil && streamMsg.ResponseMeta != nil {
					usage = handlerUsage(streamMsg.ResponseMeta)
				}
				sendMsg := buildSendMsg(ctx, toolAsAnswerMsg, false, art)
				fullContent.WriteString(streamMsg.Content)
				sendMsg.Content = streamMsg.Content
				art.MessageEvent.SendMsgEvent(entity.RunEventMessageDelta, sendMsg, art.SW)
			}
		// LLM 直接生成的回答，可能包好推理内容和最终内容
		case message.MessageTypeAnswer:
			fullContent := bytes.NewBuffer([]byte{})
			var usage *msgEntity.UsageExt
			var isToolCalls = false
			var modelAnswerMsg *msgEntity.Message
			for {
				streamMsg, receErr := chunk.ModelAnswer.Recv()
				if receErr != nil {
					if errors.Is(receErr, io.EOF) {

						if isToolCalls {
							break
						}
						if modelAnswerMsg == nil {
							break
						}
						answer := buildSendMsg(ctx, modelAnswerMsg, false, art)
						answer.Content = fullContent.String()
						hfErr := mh.handlerAnswer(ctx, answer, usage, art, modelAnswerMsg)
						if hfErr != nil {
							err = hfErr
							return
						}
						break
					}
					err = receErr
					return
				}

				if streamMsg != nil && len(streamMsg.ToolCalls) > 0 {
					isToolCalls = true
				}

				if streamMsg != nil && streamMsg.ResponseMeta != nil {
					usage = handlerUsage(streamMsg.ResponseMeta)
				}

				if streamMsg != nil && len(streamMsg.ReasoningContent) == 0 && len(streamMsg.Content) == 0 {
					continue
				}

				if len(streamMsg.ReasoningContent) > 0 {
					if reasoningMsg == nil {
						reasoningMsg, err = preCreateAnswer(ctx, art)
						if err != nil {
							return
						}
					}

					sendReasoningMsg := buildSendMsg(ctx, reasoningMsg, false, art)
					reasoningContent.WriteString(streamMsg.ReasoningContent)
					sendReasoningMsg.ReasoningContent = ptr.Of(streamMsg.ReasoningContent)
					art.MessageEvent.SendMsgEvent(entity.RunEventMessageDelta, sendReasoningMsg, art.SW)
				}
				if len(streamMsg.Content) > 0 {
					if modelAnswerMsg == nil {
						modelAnswerMsg, err = preCreateAnswer(ctx, art)
						if err != nil {
							return
						}
						if firstAnswerMsg == nil {
							if reasoningMsg != nil {
								modelAnswerMsg.ID = reasoningMsg.ID
							}
							firstAnswerMsg = modelAnswerMsg
						}
					}

					sendAnswerMsg := buildSendMsg(ctx, modelAnswerMsg, false, art)
					fullContent.WriteString(streamMsg.Content)
					sendAnswerMsg.Content = streamMsg.Content
					art.MessageEvent.SendMsgEvent(entity.RunEventMessageDelta, sendAnswerMsg, art.SW)
				}
			}
		// Agent 完成回答后，生成后续问题建议
		case message.MessageTypeFlowUp:
			if isSendFinishAnswer {

				if firstAnswerMsg != nil && len(reasoningContent.String()) > 0 {
					art.saveReasoningContent(ctx, firstAnswerMsg, reasoningContent.String())
				}

				isSendFinishAnswer = true
				finishErr := mh.handlerFinalAnswerFinish(ctx, art)
				if finishErr != nil {
					err = finishErr
					return
				}
			}

			err = mh.handlerSuggest(ctx, chunk, art)
			if err != nil {
				return
			}
		// 用户打断 Agent 输出
		case message.MessageTypeInterrupt:
			// ToDo: 支持打断
			if err != nil {
				return
			}
		}
	}
}

// pull 是生产者方法，负责从 Agent 执行引擎拉取事件并写入 mainChan。
//
// 该方法持续从 StreamReader 读取 Agent 事件，将其转换为内部事件格式后
// 写入 mainChan 供 push 方法消费。当流结束或发生错误时，会关闭通道
// 并退出循环。
//
// 事件转换流程：
//
//	crossagent.AgentEvent  →  transformEventMap()  →  entity.AgentRespEvent
//	      (外部格式)                                         (内部格式)
//
// 参数：
//   - ctx: 上下文（当前未使用，预留扩展）
//   - mainChan: 事件写入通道
//   - events: Agent 事件流读取器
func (art *AgentRuntime) pull(_ context.Context, mainChan chan *entity.AgentRespEvent, events *schema.StreamReader[*crossagent.AgentEvent]) {
	defer func() {
		close(mainChan)
	}()

	for {
		rm, re := events.Recv()
		if re != nil {
			errChunk := &entity.AgentRespEvent{
				Err: re,
			}
			mainChan <- errChunk
			return
		}

		// 外部事件转换为内部事件
		eventType, tErr := transformEventMap(rm.EventType)

		if tErr != nil {
			// 转换失败，将错误写入通道并退出
			errChunk := &entity.AgentRespEvent{
				Err: tErr,
			}
			mainChan <- errChunk
			return
		}

		// 构建内部事件响应对象
		respChunk := &entity.AgentRespEvent{
			EventType:    eventType,
			ModelAnswer:  rm.ChatModelAnswer,
			ToolsMessage: rm.ToolsMessage,
			FuncCall:     rm.FuncCall,
			Knowledge:    rm.Knowledge,
			Suggest:      rm.Suggest,
			Interrupt:    rm.Interrupt,

			ToolMidAnswer: rm.ToolMidAnswer,
			ToolAsAnswer:  rm.ToolAsChatModelAnswer,
		}

		mainChan <- respChunk
	}
}

func transformEventMap(eventType singleagent.EventType) (message.MessageType, error) {
	var eType message.MessageType
	switch eventType {
	case singleagent.EventTypeOfFuncCall:
		return message.MessageTypeFunctionCall, nil
	case singleagent.EventTypeOfKnowledge:
		return message.MessageTypeKnowledge, nil
	case singleagent.EventTypeOfToolsMessage:
		return message.MessageTypeToolResponse, nil
	case singleagent.EventTypeOfChatModelAnswer:
		return message.MessageTypeAnswer, nil
	case singleagent.EventTypeOfToolsAsChatModelStream:
		return message.MessageTypeToolAsAnswer, nil
	case singleagent.EventTypeOfToolMidAnswer:
		return message.MessageTypeToolMidAnswer, nil
	case singleagent.EventTypeOfSuggest:
		return message.MessageTypeFlowUp, nil
	case singleagent.EventTypeOfInterrupt:
		return message.MessageTypeInterrupt, nil
	}
	return eType, errorx.New(errno.ErrReplyUnknowEventType)
}

func (art *AgentRuntime) saveReasoningContent(ctx context.Context, firstAnswerMsg *msgEntity.Message, reasoningContent string) {
	_, err := crossmessage.DefaultSVC().Edit(ctx, &message.Message{
		ID:               firstAnswerMsg.ID,
		ReasoningContent: reasoningContent,
	})
	if err != nil {
		logs.Info(pkg.ModelName, "save reasoning content failed, err: %v", err)
	}
}
