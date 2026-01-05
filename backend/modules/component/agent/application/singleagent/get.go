package singleagent

import (
	"context"
	"fmt"

	"github.com/kiosk404/airi-go/backend/api/model/playground"
	"github.com/kiosk404/airi-go/backend/application/ctxutil"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg/errno"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr"
	model "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
)

func (s *SingleAgentApplicationService) GetAgentBotInfo(ctx context.Context, req *playground.GetDraftBotInfoAgwRequest) (*playground.GetDraftBotInfoAgwResponse, error) {

	uid := ctxutil.MustGetUIDFromCtx(ctx)

	agentInfo, err := s.DomainSVC.GetSingleAgent(ctx, req.GetBotID(), req.GetVersion())
	if err != nil {
		return nil, err
	}

	if agentInfo == nil {
		return nil, errorx.New(errno.ErrAgentInvalidParamCode, errorx.KVf("msg", "agent %d not found", req.GetBotID()))
	}

	if agentInfo.CreatorID != uid {
		return nil, errorx.New(errno.ErrAgentInvalidParamCode, errorx.KVf("msg", "agent %d not found", req.GetBotID()))
	}

	vo, err := s.singleAgentDraftDo2Vo(ctx, agentInfo)
	if err != nil {
		return nil, err
	}

	//klInfos, err := s.fetchKnowledgeDetails(ctx, agentInfo)
	//if err != nil {
	//	return nil, err
	//}

	modelInfos, err := s.fetchModelDetails(ctx, agentInfo)
	if err != nil {
		return nil, err
	}
	//
	//toolInfos, err := s.fetchToolDetails(ctx, agentInfo, req)
	//if err != nil {
	//	return nil, err
	//}
	//
	//pluginInfos, err := s.fetchPluginDetails(ctx, agentInfo, toolInfos)
	//if err != nil {
	//	return nil, err
	//}
	//
	//workflowInfos, err := s.fetchWorkflowDetails(ctx, agentInfo)
	//if err != nil {
	//	return nil, err
	//}
	//
	//shortCutCmdResp, err := s.fetchShortcutCMD(ctx, agentInfo)
	//if err != nil {
	//	return nil, err
	//}
	//
	//workflowDetailMap, err := workflowDo2Vo(workflowInfos)
	//if err != nil {
	//	return nil, err
	//}

	return &playground.GetDraftBotInfoAgwResponse{
		Data: &playground.GetDraftBotInfoAgwData{
			BotInfo: vo,
			BotOptionData: &playground.BotOptionData{
				ModelDetailMap: modelInfoDo2Vo(modelInfos),
				//KnowledgeDetailMap:  knowledgeInfoDo2Vo(klInfos),
				//PluginAPIDetailMap:  toolInfoDo2Vo(toolInfos),
				//PluginDetailMap:     s.pluginInfoDo2Vo(ctx, pluginInfos),
				//WorkflowDetailMap:   workflowDetailMap,
				//ShortcutCommandList: shortCutCmdResp,
			},
			Editable:  ptr.Of(true),
			Deletable: ptr.Of(true),
		},
	}, nil
}

func modelInfoDo2Vo(modelInfos []*model.Model) map[int64]*playground.ModelDetail {
	return slices.ToMap(modelInfos, func(e *model.Model) (int64, *playground.ModelDetail) {
		return e.ID, toModelDetail(e)
	})
}

func toModelDetail(m *model.Model) *playground.ModelDetail {
	return &playground.ModelDetail{
		Name:         ptr.Of(m.DisplayInfo.Name),
		ModelName:    ptr.Of(m.Connection.BaseConnInfo.Model),
		ModelID:      ptr.Of(m.ID),
		ModelFamily:  ptr.Of(int64(m.Provider.ModelClass)),
		ModelIconURL: ptr.Of(m.Provider.IconURL),
	}
}

func (s *SingleAgentApplicationService) fetchModelDetails(ctx context.Context, agentInfo *entity.SingleAgent) ([]*model.Model, error) {
	if agentInfo.ModelInfo.ModelId == nil {
		return nil, nil
	}

	modelID := agentInfo.ModelInfo.GetModelId()

	modelInfo, err := modelmgr.DefaultSVC().GetModelByID(ctx, modelID)

	if err != nil {
		return nil, fmt.Errorf("fetch model(%d) details failed: %v", modelID, err)
	}

	return []*model.Model{modelInfo}, nil
}

//
//func (s *SingleAgentApplicationService) fetchKnowledgeDetails(ctx context.Context, agentInfo *entity.SingleAgent) ([]*knowledgeModel.Knowledge, error) {
//	knowledgeIDs := make([]int64, 0, len(agentInfo.Knowledge.KnowledgeInfo))
//	for _, v := range agentInfo.Knowledge.KnowledgeInfo {
//		id, err := conv.StrToInt64(v.GetId())
//		if err != nil {
//			return nil, fmt.Errorf("invalid knowledge id: %s", v.GetId())
//		}
//		knowledgeIDs = append(knowledgeIDs, id)
//	}
//
//	if len(knowledgeIDs) == 0 {
//		return nil, nil
//	}
//
//	listResp, err := s.appContext.KnowledgeDomainSVC.ListKnowledge(ctx, &knowledge.ListKnowledgeRequest{
//		IDs: knowledgeIDs,
//	})
//	if err != nil {
//		return nil, fmt.Errorf("fetch knowledge details failed: %v", err)
//	}
//
//	return listResp.KnowledgeList, err
//}
//
//func (s *SingleAgentApplicationService) fetchToolDetails(ctx context.Context, agentInfo *entity.SingleAgent, req *playground.GetDraftBotInfoAgwRequest) ([]*pluginEntity.ToolInfo, error) {
//	return s.appContext.PluginDomainSVC.MGetAgentTools(ctx, &model.MGetAgentToolsRequest{
//		SpaceID: agentInfo.SpaceID,
//		AgentID: req.GetBotID(),
//		IsDraft: true,
//		VersionAgentTools: slices.Transform(agentInfo.Plugin, func(a *bot_common.PluginInfo) model.VersionAgentTool {
//			return model.VersionAgentTool{
//				ToolID:     a.GetApiId(),
//				PluginFrom: a.PluginFrom,
//				PluginID:   a.GetPluginId(),
//			}
//		}),
//	})
//}
//
//func (s *SingleAgentApplicationService) fetchPluginDetails(ctx context.Context, agentInfo *entity.SingleAgent, toolInfos []*pluginEntity.ToolInfo) ([]*pluginEntity.PluginInfo, error) {
//	vLocalPlugins := make([]model.VersionPlugin, 0, len(agentInfo.Plugin))
//	vPluginMap := make(map[string]bool, len(agentInfo.Plugin))
//	vSaasPlugin := make([]model.VersionPlugin, 0, len(agentInfo.Plugin))
//	for _, v := range toolInfos {
//		k := fmt.Sprintf("%d:%s:%s", v.PluginID, v.GetVersion(), v.GetPluginFrom())
//		if vPluginMap[k] {
//			continue
//		}
//		vPluginMap[k] = true
//		if v.GetPluginFrom() == bot_common.PluginFrom_FromSaas {
//			vSaasPlugin = append(vSaasPlugin, model.VersionPlugin{
//				PluginID: v.PluginID,
//				Version:  v.GetVersion(),
//			})
//		} else {
//			vLocalPlugins = append(vLocalPlugins, model.VersionPlugin{
//				PluginID: v.PluginID,
//				Version:  v.GetVersion(),
//			})
//		}
//	}
//	pluginInfos := make([]*pluginEntity.PluginInfo, 0, len(vLocalPlugins)+len(vSaasPlugin))
//	if len(vLocalPlugins) > 0 {
//		localPluginInfos, err := s.appContext.PluginDomainSVC.MGetVersionPlugins(ctx, vLocalPlugins)
//		if err != nil {
//			return nil, fmt.Errorf("fetch local plugin details failed: %v", err)
//		}
//		pluginInfos = append(pluginInfos, localPluginInfos...)
//	}
//
//	if len(vSaasPlugin) > 0 {
//		saasPluginInfos, err := s.appContext.PluginDomainSVC.GetSaasPluginInfo(ctx, slices.Transform(vSaasPlugin, func(v model.VersionPlugin) int64 {
//			return v.PluginID
//		}))
//		if err != nil {
//			return nil, fmt.Errorf("fetch saas plugin details failed: %v", err)
//		}
//		pluginInfos = append(pluginInfos, saasPluginInfos...)
//	}
//	return pluginInfos, nil
//}
//
//func (s *SingleAgentApplicationService) fetchWorkflowDetails(ctx context.Context, agentInfo *entity.SingleAgent) ([]*workflowEntity.Workflow, error) {
//	if len(agentInfo.Workflow) == 0 {
//		return nil, nil
//	}
//
//	policy := &vo.MGetPolicy{
//		MetaQuery: vo.MetaQuery{
//			IDs: slices.Transform(agentInfo.Workflow, func(a *bot_common.WorkflowInfo) int64 {
//				return a.GetWorkflowId()
//			}),
//		},
//		QType: workflowModel.FromLatestVersion,
//	}
//	ret, _, err := s.appContext.WorkflowDomainSVC.MGet(ctx, policy)
//	if err != nil {
//		return nil, fmt.Errorf("fetch workflow details failed: %v", err)
//	}
//	return ret, nil
//}
