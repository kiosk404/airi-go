package agentflow

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	crossplugin "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin"
	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/consts"
	model2 "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/model"
	pluginEntity "github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
)

type toolConfig struct {
	userID        string
	agentIdentity *entity.AgentIdentity
	toolConf      []*bot_common.PluginInfo

	conversationID int64
}

func newPluginTools(ctx context.Context, conf *toolConfig) ([]tool.InvokableTool, error) {
	req := &model2.MGetAgentToolsRequest{
		AgentID: conf.agentIdentity.AgentID,
		IsDraft: conf.agentIdentity.IsDraft,
		VersionAgentTools: slices.Transform(conf.toolConf, func(a *bot_common.PluginInfo) model2.VersionAgentTool {
			return model2.VersionAgentTool{
				ToolID:       a.GetApiId(),
				AgentVersion: ptr.Of(conf.agentIdentity.Version),
				PluginFrom:   a.PluginFrom,
				PluginID:     a.GetPluginId(),
			}
		}),
	}
	agentTools, err := crossplugin.DefaultSVC().MGetAgentTools(ctx, req)
	if err != nil {
		return nil, err
	}

	projectInfo := &model2.ProjectInfo{
		ProjectID:      conf.agentIdentity.AgentID,
		ProjectType:    consts.ProjectTypeOfAgent,
		ProjectVersion: ptr.Of(conf.agentIdentity.Version),
	}

	tools := make([]tool.InvokableTool, 0, len(agentTools))
	for _, ti := range agentTools {
		tools = append(tools, &pluginInvokableTool{
			userID:      conf.userID,
			isDraft:     conf.agentIdentity.IsDraft,
			projectInfo: projectInfo,
			toolInfo:    ti,
			pluginFrom:  ti.Source,

			conversationID: conf.conversationID,
		})
	}

	return tools, nil
}

type pluginInvokableTool struct {
	userID      string
	isDraft     bool
	toolInfo    *pluginEntity.ToolInfo
	projectInfo *model2.ProjectInfo

	pluginFrom *bot_common.PluginFrom

	conversationID int64
}

func (p *pluginInvokableTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	paramInfos, err := p.toolInfo.Operation.ToEinoSchemaParameterInfo(ctx)
	if err != nil {
		return nil, err
	}

	if len(paramInfos) == 0 {
		return &schema.ToolInfo{
			Name:        p.toolInfo.GetName(),
			Desc:        p.toolInfo.GetDesc(),
			ParamsOneOf: nil,
		}, nil
	}

	return &schema.ToolInfo{
		Name:        p.toolInfo.GetName(),
		Desc:        p.toolInfo.GetDesc(),
		ParamsOneOf: schema.NewParamsOneOfByParams(paramInfos),
	}, nil
}

func (p *pluginInvokableTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	req := &model2.ExecuteToolRequest{
		UserID:          p.userID,
		PluginID:        p.toolInfo.PluginID,
		ToolID:          p.toolInfo.ID,
		ExecDraftTool:   false,
		PluginFrom:      p.pluginFrom,
		ArgumentsInJson: argumentsInJSON,
		ExecScene: func() consts.ExecuteScene {
			if p.isDraft {
				return consts.ExecSceneOfDraftAgent
			}
			return consts.ExecSceneOfOnlineAgent
		}(),
	}

	opts := []model2.ExecuteToolOpt{
		model2.WithInvalidRespProcessStrategy(consts.InvalidResponseProcessStrategyOfReturnDefault),
		model2.WithToolVersion(p.toolInfo.GetVersion()),
		model2.WithProjectInfo(p.projectInfo),
		model2.WithPluginHTTPHeader(p.conversationID),
	}

	resp, err := crossplugin.DefaultSVC().ExecuteTool(ctx, req, opts...)
	if err != nil {
		return "", err
	}

	return resp.TrimmedResp, nil
}
