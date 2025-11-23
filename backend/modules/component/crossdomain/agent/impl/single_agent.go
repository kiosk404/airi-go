package impl

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/infra/contract/imagex"
	singleagent "github.com/kiosk404/airi-go/backend/modules/component/agent/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg"
	crossagent "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/agent"
	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/agent/model"
	agentrun "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/agentrun/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

var defaultSVC crossagent.SingleAgent

type impl struct {
	DomainSVC singleagent.SingleAgent
	ImagexSVC imagex.ImageX
}

func InitDomainService(c singleagent.SingleAgent, imagexClient imagex.ImageX) crossagent.SingleAgent {
	defaultSVC = &impl{
		DomainSVC: c,
		ImagexSVC: imagexClient,
	}

	return defaultSVC
}

func (c *impl) StreamExecute(ctx context.Context, agentRuntime *crossagent.AgentRuntime,
) (*schema.StreamReader[*model.AgentEvent], error) {

	singleAgentStreamExecReq := c.buildSingleAgentStreamExecuteReq(ctx, agentRuntime)

	streamEvent, err := c.DomainSVC.StreamExecute(ctx, singleAgentStreamExecReq)
	logs.InfoX(pkg.ModelName, "agent StreamExecute req:%v, streamEvent:%v, err:%v", conv.DebugJsonToStr(singleAgentStreamExecReq), streamEvent, err)
	return streamEvent, err
}

func (c *impl) buildSingleAgentStreamExecuteReq(ctx context.Context, agentRuntime *crossagent.AgentRuntime,
) *model.ExecuteRequest {

	return &model.ExecuteRequest{
		Identity:        c.buildIdentity(agentRuntime),
		Input:           agentRuntime.Input,
		History:         agentRuntime.HistoryMsg,
		UserID:          agentRuntime.UserID,
		CustomVariables: agentRuntime.CustomVariables,
		PreCallTools: slices.Transform(agentRuntime.PreRetrieveTools, func(tool *agentrun.Tool) *agentrun.ToolsRetriever {
			return &agentrun.ToolsRetriever{
				PluginID:  tool.PluginID,
				ToolName:  tool.ToolName,
				ToolID:    tool.ToolID,
				Arguments: tool.Arguments,
				Type: func(toolType agentrun.ToolType) agentrun.ToolType {
					switch toolType {
					case agentrun.ToolTypeWorkflow:
						return agentrun.ToolTypeWorkflow
					case agentrun.ToolTypePlugin:
						return agentrun.ToolTypePlugin
					}
					return agentrun.ToolTypePlugin
				}(tool.Type),
				//PluginFrom: tool.PluginFrom,
			}
		}),
		ResumeInfo: agentRuntime.ResumeInfo,

		ConversationID: agentRuntime.ConversationId,
	}
}

func (c *impl) buildIdentity(agentRuntime *crossagent.AgentRuntime) *model.AgentIdentity {
	return &model.AgentIdentity{
		AgentID: agentRuntime.AgentID,
		Version: agentRuntime.AgentVersion,
		IsDraft: agentRuntime.IsDraft,
	}
}

func (c *impl) ObtainAgentByIdentity(ctx context.Context, identity *model.AgentIdentity) (*model.SingleAgent, error) {
	agentInfo, err := c.DomainSVC.ObtainAgentByIdentity(ctx, identity)
	if err != nil {
		return nil, err
	}
	if agentInfo == nil {
		return nil, nil
	}
	return agentInfo.SingleAgent, nil
}

func (c *impl) GetSingleAgentDraft(ctx context.Context, agentID int64) (*model.SingleAgent, error) {
	agentInfo, err := c.DomainSVC.GetSingleAgentDraft(ctx, agentID)
	if err != nil {
		return nil, err
	}
	if agentInfo == nil {
		return nil, nil
	}
	return agentInfo.SingleAgent, nil
}
