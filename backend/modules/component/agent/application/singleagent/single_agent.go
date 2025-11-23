package singleagent

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	"github.com/kiosk404/airi-go/backend/api/model/playground"
	"github.com/kiosk404/airi-go/backend/application/ctxutil"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	singleagent "github.com/kiosk404/airi-go/backend/modules/component/agent/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

type SingleAgentApplicationService struct {
	appContext *ServiceComponents
	DomainSVC  singleagent.SingleAgent
}

func newApplicationService(s *ServiceComponents, domain singleagent.SingleAgent) *SingleAgentApplicationService {
	return &SingleAgentApplicationService{
		appContext: s,
		DomainSVC:  domain,
	}
}

const onboardingInfoMaxLength = 65535

func (s *SingleAgentApplicationService) generateOnboardingStr(onboardingInfo *bot_common.OnboardingInfo) (string, error) {
	onboarding := playground.OnboardingContent{}
	if onboardingInfo != nil {
		onboarding.Prologue = ptr.Of(onboardingInfo.GetPrologue())
		onboarding.SuggestedQuestions = onboardingInfo.GetSuggestedQuestions()
		onboarding.SuggestedQuestionsShowMode = onboardingInfo.SuggestedQuestionsShowMode
	}

	onboardingInfoStr, err := sonic.MarshalString(onboarding)
	if err != nil {
		return "", err
	}

	return onboardingInfoStr, nil
}

func (s *SingleAgentApplicationService) applyAgentUpdates(target *entity.SingleAgent, patch *bot_common.BotInfoForUpdate) (*entity.SingleAgent, error) {
	if patch.Name != nil {
		target.Name = *patch.Name
	}

	if patch.Description != nil {
		target.Desc = *patch.Description
	}

	if patch.IconUri != nil {
		target.IconURI = *patch.IconUri
	}

	if patch.OnboardingInfo != nil {
		target.OnboardingInfo = patch.OnboardingInfo
	}

	if patch.ModelInfo != nil {
		target.ModelInfo = patch.ModelInfo
	}

	if patch.PromptInfo != nil {
		target.Prompt = patch.PromptInfo
	}

	if patch.WorkflowInfoList != nil {
		target.Workflow = patch.WorkflowInfoList
	}

	if patch.PluginInfoList != nil {
		target.Plugin = patch.PluginInfoList
	}

	if patch.Knowledge != nil {
		target.Knowledge = patch.Knowledge
	}

	if patch.SuggestReplyInfo != nil {
		target.SuggestReply = patch.SuggestReplyInfo
	}

	if patch.BackgroundImageInfoList != nil {
		target.BackgroundImageInfoList = patch.BackgroundImageInfoList
	}

	if patch.Agents != nil && len(patch.Agents) > 0 && patch.Agents[0].JumpConfig != nil {
		target.JumpConfig = patch.Agents[0].JumpConfig
	}

	if patch.ShortcutSort != nil {
		target.ShortcutCommand = patch.ShortcutSort
	}

	if patch.DatabaseList != nil {
		for _, db := range patch.DatabaseList {
			if db.PromptDisabled == nil {
				db.PromptDisabled = ptr.Of(false) // default is false
			}
		}
		target.Database = patch.DatabaseList
	}

	return target, nil
}

func (s *SingleAgentApplicationService) UpdateSingleAgentDraft(ctx context.Context, req *playground.UpdateDraftBotInfoAgwRequest) (*playground.UpdateDraftBotInfoAgwResponse, error) {
	if req.BotInfo.OnboardingInfo != nil {
		infoStr, err := s.generateOnboardingStr(req.BotInfo.OnboardingInfo)
		if err != nil {
			return nil, errorx.New(errno.ErrAgentPermissionCode, errorx.KV("msg", "onboarding_info invalidate"))
		}

		if len(infoStr) > onboardingInfoMaxLength {
			return nil, errorx.New(errno.ErrAgentPermissionCode, errorx.KV("msg", "onboarding_info is too long"))
		}
	}

	agentID := req.BotInfo.GetBotId()
	currentAgentInfo, err := s.ValidateAgentDraftAccess(ctx, agentID)
	if err != nil {
		return nil, err
	}

	userID := ctxutil.MustGetUIDFromCtx(ctx)
	logs.InfoX(pkg.ModelName, "update single agent info %s draft by user %d", agentID, userID)

	updateAgentInfo, err := s.applyAgentUpdates(currentAgentInfo, req.BotInfo)
	if err != nil {
		return nil, err
	}

	if req.BotInfo.VariableList != nil {
		//var (
		//	varsMetaID int64
		//	vars       = variableEntity.NewVariablesWithAgentVariables(req.BotInfo.VariableList)
		//)
		//
		//varsMetaID, err = s.appContext.VariablesDomainSVC.UpsertBotMeta(ctx, agentID, "", userID, vars)
		//if err != nil {
		//	return nil, err
		//}
		//
		//updateAgentInfo.VariablesMetaID = &varsMetaID
	}

	err = s.DomainSVC.UpdateSingleAgentDraft(ctx, updateAgentInfo)
	if err != nil {
		return nil, err
	}

	//err = s.appContext.EventBus.PublishProject(ctx, &searchEntity.ProjectDomainEvent{
	//	OpType: searchEntity.Updated,
	//	Project: &searchEntity.ProjectDocument{
	//		ID:   agentID,
	//		Name: &updateAgentInfo.Name,
	//		Type: intelligence.IntelligenceType_Bot,
	//	},
	//})
	//if err != nil {
	//	return nil, err
	//}

	return &playground.UpdateDraftBotInfoAgwResponse{
		Data: &playground.UpdateDraftBotInfoAgwData{
			HasChange:    ptr.Of(true),
			CheckNotPass: false,
			Branch:       playground.BranchPtr(playground.Branch_PersonalDraft),
		},
	}, nil
}

func (s *SingleAgentApplicationService) ValidateAgentDraftAccess(ctx context.Context, agentID int64) (*entity.SingleAgent, error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrAgentPermissionCode, errorx.KV("msg", "session uid not found"))
	}

	do, err := s.DomainSVC.GetSingleAgentDraft(ctx, agentID)
	if err != nil {
		return nil, err
	}

	if do == nil {
		return nil, errorx.New(errno.ErrAgentPermissionCode, errorx.KVf("msg", "No agent draft(%d) found for the given agent ID", agentID))
	}

	return do, nil
}
