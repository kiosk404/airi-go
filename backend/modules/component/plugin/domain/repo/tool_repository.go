package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/model"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/infra/dao"
)

type ToolRepository interface {
	CreateDraftTool(ctx context.Context, tool *entity.ToolInfo) (toolID int64, err error)
	UpsertDraftTools(ctx context.Context, pluginID int64, tools []*entity.ToolInfo) (err error)
	UpdateDraftTool(ctx context.Context, tool *entity.ToolInfo) (err error)
	GetDraftTool(ctx context.Context, toolID int64) (tool *entity.ToolInfo, exist bool, err error)
	MGetDraftTools(ctx context.Context, toolIDs []int64, opts ...ToolSelectedOptions) (tools []*entity.ToolInfo, err error)

	GetDraftToolWithAPI(ctx context.Context, pluginID int64, api dao.UniqueToolAPI) (tool *entity.ToolInfo, exist bool, err error)
	MGetDraftToolWithAPI(ctx context.Context, pluginID int64, apis []dao.UniqueToolAPI, opts ...ToolSelectedOptions) (tools map[dao.UniqueToolAPI]*entity.ToolInfo, err error)
	DeleteDraftTool(ctx context.Context, toolID int64) (err error)

	GetOnlineTool(ctx context.Context, toolID int64) (tool *entity.ToolInfo, exist bool, err error)
	MGetOnlineTools(ctx context.Context, toolIDs []int64, opts ...ToolSelectedOptions) (tools []*entity.ToolInfo, err error)

	GetVersionTool(ctx context.Context, vTool model.VersionTool) (tool *entity.ToolInfo, exist bool, err error)
	MGetVersionTools(ctx context.Context, vTools []model.VersionTool) (tools []*entity.ToolInfo, err error)

	BindDraftAgentTools(ctx context.Context, agentID int64, bindTools []*model.BindToolInfo) (err error)
	DuplicateDraftAgentTools(ctx context.Context, fromAgentID, toAgentID int64) (err error)
	GetDraftAgentTool(ctx context.Context, agentID, toolID int64) (tool *entity.ToolInfo, exist bool, err error)
	GetDraftAgentToolWithToolName(ctx context.Context, agentID int64, toolName string) (tool *entity.ToolInfo, exist bool, err error)
	MGetDraftAgentTools(ctx context.Context, agentID int64, toolIDs []int64) (tools []*entity.ToolInfo, err error)
	UpdateDraftAgentTool(ctx context.Context, req *UpdateDraftAgentToolRequest) (err error)
	GetSpaceAllDraftAgentTools(ctx context.Context, agentID int64) (tools []*entity.ToolInfo, err error)
	GetAgentPluginIDs(ctx context.Context, agentID int64) (pluginIDs []int64, err error)

	GetVersionAgentTool(ctx context.Context, agentID int64, vAgentTool model.VersionAgentTool) (tool *entity.ToolInfo, exist bool, err error)
	GetVersionAgentToolWithToolName(ctx context.Context, req *GetVersionAgentToolWithToolNameRequest) (tool *entity.ToolInfo, exist bool, err error)
	MGetVersionAgentTool(ctx context.Context, agentID int64, vAgentTools []model.VersionAgentTool) (tools []*entity.ToolInfo, err error)
	BatchCreateVersionAgentTools(ctx context.Context, agentID int64, agentVersion string, tools []*entity.ToolInfo) (err error)

	GetPluginAllDraftTools(ctx context.Context, pluginID int64, opts ...ToolSelectedOptions) (tools []*entity.ToolInfo, err error)
	GetPluginAllOnlineTools(ctx context.Context, pluginID int64) (tools []*entity.ToolInfo, err error)
	ListPluginDraftTools(ctx context.Context, pluginID int64, pageInfo dao.PageInfo) (tools []*entity.ToolInfo, total int64, err error)

	// SaaS plugin tools
	BatchGetSaasPluginToolsInfo(ctx context.Context, pluginIDs []int64) (tools map[int64][]*entity.ToolInfo, plugins map[int64]*entity.PluginInfo, err error)
}

type GetVersionAgentToolWithToolNameRequest struct {
	AgentID      int64
	ToolName     string
	AgentVersion *string
}

type UpdateDraftAgentToolRequest struct {
	AgentID  int64
	ToolName string
	Tool     *entity.ToolInfo
}
