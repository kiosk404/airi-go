package agentflow

import (
	"context"
	"errors"
	"io"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg"
	singleagent "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/agent/model"
	model2 "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/model"
	"github.com/kiosk404/airi-go/backend/pkg/json"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

func newReplyCallback(_ context.Context, executeID string, returnDirectlyTools mapset.Set[string]) (
	clb callbacks.Handler, sr *schema.StreamReader[*entity.AgentEvent], sw *schema.StreamWriter[*entity.AgentEvent],
) {
	sr, sw = schema.Pipe[*entity.AgentEvent](10)

	rcc := &replyChunkCallback{
		sw:                  sw,
		executeID:           executeID,
		returnDirectlyTools: returnDirectlyTools,
	}

	clb = callbacks.NewHandlerBuilder().
		OnStartFn(rcc.OnStart).
		OnEndFn(rcc.OnEnd).
		OnEndWithStreamOutputFn(rcc.OnEndWithStreamOutput).
		OnErrorFn(rcc.OnError).
		Build()

	return clb, sr, sw
}

type replyChunkCallback struct {
	sw                  *schema.StreamWriter[*entity.AgentEvent]
	executeID           string
	returnDirectlyTools mapset.Set[string]
}

func (r *replyChunkCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	logs.InfoX(pkg.ModelName, "info-OnError, info=%v, err=%v", conv.DebugJsonToStr(info), err)

	switch info.Component {
	case compose.ComponentOfGraph:
		// 暂停情况
		if interruptInfo, ok := compose.ExtractInterruptInfo(err); ok {
			if info.Name != "" {
				return ctx
			}

			// 提取中断的信息
			interruptData := convInterruptInfo(ctx, interruptInfo)
			interruptData.InterruptID = r.executeID

			// 发送中断事件给客户端，提示需要用户回应（如输入什么信息, 等待授权什么的）
			toolMessageEvent := &entity.AgentEvent{
				EventType: singleagent.EventTypeOfToolsMessage,
				ToolsMessage: []*schema.Message{
					{
						Role:       schema.Tool,
						Content:    "directly streaming reply",
						ToolCallID: interruptData.ToolCallID,
					},
				},
			}
			r.sw.Send(toolMessageEvent, nil)

			// 发送中断事件给客户端
			interruptEvent := &entity.AgentEvent{
				EventType: singleagent.EventTypeOfInterrupt,
				Interrupt: interruptData,
			}
			r.sw.Send(interruptEvent, nil)

		} else {
			// 普通错误
			logs.ErrorX(pkg.ModelName, "[AgentRunError] | node execute failed, component=%v, name=%v, err=%v",
				info.Component, info.Name, err)
			r.sw.Send(nil, err)
		}
	}
	return ctx
}

func (r *replyChunkCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	logs.InfoX(pkg.ModelName, "info-OnStart, info=%v, input=%v", conv.DebugJsonToStr(info), conv.DebugJsonToStr(input))

	switch info.Component {
	case compose.ComponentOfToolsNode:
		if info.Name != keyOfReActAgentToolsNode {
			return ctx
		}
		ae := &entity.AgentEvent{
			EventType: singleagent.EventTypeOfFuncCall,
			FuncCall:  convToolsNodeCallbackInput(input),
		}
		r.sw.Send(ae, nil)
	}

	return ctx
}

func (r *replyChunkCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	logs.InfoX(pkg.ModelName, "info-OnEnd, info=%v, output=%v", conv.DebugJsonToStr(info), conv.DebugJsonToStr(output))
	switch info.Name {
	case keyOfKnowledgeRetriever:
		knowledgeEvent := &entity.AgentEvent{
			EventType: singleagent.EventTypeOfKnowledge,
			Knowledge: retriever.ConvCallbackOutput(output).Docs,
		}

		if knowledgeEvent.Knowledge != nil {
			r.sw.Send(knowledgeEvent, nil)
		}
	case keyOfToolsPreRetriever:
		result := convToolsPreRetrieverCallbackInput(output)

		if len(result) > 0 {
			for _, item := range result {
				var event *entity.AgentEvent
				if item.Role == schema.Tool {
					event = &entity.AgentEvent{
						EventType:    singleagent.EventTypeOfToolsMessage,
						ToolsMessage: []*schema.Message{item},
					}
				} else {
					event = &entity.AgentEvent{
						EventType: singleagent.EventTypeOfFuncCall,
						FuncCall:  item,
					}
				}
				r.sw.Send(event, nil)
			}
		}

	case keyOfSuggestParser:
		sg := convSuggestionNodeCallbackOutput(output)

		if len(sg) > 0 {
			for _, item := range sg {
				suggestionEvent := &entity.AgentEvent{
					EventType: singleagent.EventTypeOfSuggest,
					Suggest:   item,
				}
				r.sw.Send(suggestionEvent, nil)
			}
		}

	default:
		return ctx
	}

	return ctx
}

func (r *replyChunkCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo,
	output *schema.StreamReader[callbacks.CallbackOutput],
) context.Context {
	logs.InfoX(pkg.ModelName, "info-OnEndWithStreamOutput, info=%v, output=%v", conv.DebugJsonToStr(info), conv.DebugJsonToStr(output))
	switch info.Component {
	case compose.ComponentOfGraph, components.ComponentOfChatModel:
		if info.Name != keyOfReActAgentChatModel && info.Name != keyOfLLM {
			output.Close()
			return ctx
		}
		sr := schema.StreamReaderWithConvert(output, func(t callbacks.CallbackOutput) (*schema.Message, error) {
			cbOut := model.ConvCallbackOutput(t)
			return cbOut.Message, nil
		})

		r.sw.Send(&entity.AgentEvent{
			EventType:       singleagent.EventTypeOfChatModelAnswer,
			ChatModelAnswer: sr,
		}, nil)
		return ctx
	case compose.ComponentOfToolsNode:
		toolsMessage, err := r.concatToolsNodeOutput(ctx, output)
		if err != nil {
			r.sw.Send(nil, err)
			return ctx
		}

		r.sw.Send(&entity.AgentEvent{
			EventType:    singleagent.EventTypeOfToolsMessage,
			ToolsMessage: toolsMessage,
		}, nil)
		return ctx
	default:
		return ctx
	}
}

func convToolsNodeCallbackInput(input callbacks.CallbackInput) *schema.Message {
	switch t := input.(type) {
	case *schema.Message:
		return t
	default:
		return nil
	}
}

func convInterruptInfo(ctx context.Context, interruptInfo *compose.InterruptInfo) *singleagent.InterruptInfo {
	var output *compose.InterruptInfo
	output = interruptInfo.SubGraphs[keyOfReActAgent]
	var extra any

	for i := range output.RerunNodesExtra {
		extra = output.RerunNodesExtra[i]
		break
	}
	toolsNodeExtra, ok := extra.(*compose.ToolsInterruptAndRerunExtra)
	logs.InfoX(pkg.ModelName, "toolsNodeExtra=%v, err=%v", toolsNodeExtra, ok)

	var toolCallID string

	//wfResumeData := make(map[string]*crossworkflow.ToolInterruptEvent)
	toolResultData := make(map[string]*model2.ToolInterruptEvent)
	var interruptEventType singleagent.InterruptEventType
	//for k, v := range toolsNodeExtra.RerunExtraMap {
	//	toolCallID = k
	//
	//	interruptEventType = convInterruptEventType(v)
	//
	//	if interruptEventType == singleagent.InterruptEventType_OauthPlugin {
	//		toolResultData[k] = v.(*plugin.ToolInterruptEvent)
	//	} else {
	//		wfResumeData[k] = v.(*crossworkflow.ToolInterruptEvent)
	//	}
	//	break
	//}

	interrupt := &singleagent.InterruptInfo{
		AllToolInterruptData: toolResultData,
		//AllWfInterruptData:   wfResumeData,
		ToolCallID:    toolCallID,
		InterruptType: interruptEventType,
	}
	return interrupt
}

func convToolsPreRetrieverCallbackInput(output callbacks.CallbackOutput) []*schema.Message {
	switch t := output.(type) {
	case []*schema.Message:
		return t
	default:
		return nil
	}
}

func convToolsNodeCallbackOutput(output callbacks.CallbackOutput) []*schema.Message {
	switch t := output.(type) {
	case []*schema.Message:
		return t
	default:
		return nil
	}
}

func convSuggestionNodeCallbackOutput(output callbacks.CallbackInput) []*schema.Message {
	var sg []*schema.Message

	switch so := output.(type) {
	case *schema.Message:
		if so.Content != "" {
			var suggestions []string

			err := json.Unmarshal([]byte(so.Content), &suggestions)

			if err == nil && len(suggestions) > 0 {
				for _, suggestion := range suggestions {
					sm := &schema.Message{
						Role:         so.Role,
						Content:      suggestion,
						ResponseMeta: so.ResponseMeta,
					}
					sg = append(sg, sm)
				}
			}
		}
	default:
		return sg
	}

	return sg
}

func (r *replyChunkCallback) concatToolsNodeOutput(ctx context.Context, output *schema.StreamReader[callbacks.CallbackOutput]) ([]*schema.Message, error) {
	var toolsMsgChunks [][]*schema.Message
	var sr *schema.StreamReader[*schema.Message]
	var sw *schema.StreamWriter[*schema.Message]
	defer func() {
		if sw != nil {
			sw.Close()
		}
	}()
	var streamInitialized bool
	returnDirectToolsMap := make(map[int]bool)
	isReturnDirectToolsFirstCheck := true
	isToolsMsgChunksInit := false
	for {
		cbOut, err := output.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			if sw != nil {
				sw.Send(nil, err)
			}
			return nil, err
		}

		msgs := convToolsNodeCallbackOutput(cbOut)

		if !isToolsMsgChunksInit {
			isToolsMsgChunksInit = true
			toolsMsgChunks = make([][]*schema.Message, len(msgs))
		}

		for mIndex, msg := range msgs {

			if msg == nil {
				continue
			}
			if r.returnDirectlyTools.Cardinality() > 0 {
				if isReturnDirectToolsFirstCheck {
					isReturnDirectToolsFirstCheck = false
					if r.returnDirectlyTools.Contains(msg.ToolName) {
						returnDirectToolsMap[mIndex] = true
					}
				}

				if _, ok := returnDirectToolsMap[mIndex]; ok {
					if !streamInitialized {
						sr, sw = schema.Pipe[*schema.Message](5)
						r.sw.Send(&entity.AgentEvent{
							EventType:             singleagent.EventTypeOfToolsAsChatModelStream,
							ToolAsChatModelAnswer: sr,
						}, nil)
						streamInitialized = true
					}
					sw.Send(msg, nil)
				}
			}
			if toolsMsgChunks[mIndex] == nil {
				toolsMsgChunks[mIndex] = []*schema.Message{msg}
			} else {
				toolsMsgChunks[mIndex] = append(toolsMsgChunks[mIndex], msg)
			}
		}
	}

	toolMessages := make([]*schema.Message, 0, len(toolsMsgChunks))

	for _, msgChunks := range toolsMsgChunks {
		msg, err := schema.ConcatMessages(msgChunks)
		if err != nil {
			return nil, err
		}
		toolMessages = append(toolMessages, msg)
	}

	return toolMessages, nil
}
