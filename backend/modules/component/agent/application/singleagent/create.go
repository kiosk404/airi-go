package singleagent

import (
	"context"
	"time"

	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	"github.com/kiosk404/airi-go/backend/api/model/app/developer_api"
	"github.com/kiosk404/airi-go/backend/api/model/llm/domain/common"
	"github.com/kiosk404/airi-go/backend/api/model/llm/manage"
	"github.com/kiosk404/airi-go/backend/application/ctxutil"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg/errno"
	singleagent "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/agent/model"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

func (s *SingleAgentApplicationService) CreateSingleAgentDraft(ctx context.Context, req *developer_api.DraftBotCreateRequest) (*developer_api.DraftBotCreateResponse, error) {
	model, err := modelmgr.DefaultSVC().GetOnlineDefaultModel(ctx)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, errorx.New(errno.ErrAgentNoModelInUseCode)
	}
	do, err := s.draftBotCreateRequestToSingleAgent(ctx, req)
	if err != nil {
		return nil, err
	}
	userID := ctxutil.MustGetUIDFromCtx(ctx)
	agentID, err := s.DomainSVC.CreateSingleAgentDraft(ctx, userID, do)
	if err != nil {
		return nil, err
	}

	logs.InfoX(pkg.ModelName, "create single draft %d from user %d", agentID, userID)
	return &developer_api.DraftBotCreateResponse{Data: &developer_api.DraftBotCreateData{
		BotID: agentID,
	}}, nil
}

func (s *SingleAgentApplicationService) draftBotCreateRequestToSingleAgent(ctx context.Context, req *developer_api.DraftBotCreateRequest) (*entity.SingleAgent, error) {
	sa, err := s.newDefaultSingleAgent(ctx)
	if err != nil {
		return nil, err
	}

	id, err := s.appContext.IDGen.GenID(ctx)
	sa.AgentID = id
	sa.Name = req.GetName()
	sa.Desc = req.GetDescription()
	sa.IconURI = req.GetIconURI()

	return sa, nil
}

func (s *SingleAgentApplicationService) newDefaultSingleAgent(ctx context.Context) (*entity.SingleAgent, error) {
	mi, err := s.defaultModelInfo(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	return &entity.SingleAgent{
		SingleAgent: &singleagent.SingleAgent{
			OnboardingInfo: &bot_common.OnboardingInfo{},
			ModelInfo:      mi,
			Prompt:         &bot_common.PromptInfo{},
			Plugin:         []*bot_common.PluginInfo{},
			Knowledge: &bot_common.Knowledge{
				TopK:           ptr.Of(int64(1)),
				MinScore:       ptr.Of(0.01),
				SearchStrategy: ptr.Of(bot_common.SearchStrategy_SemanticSearch),
				RecallStrategy: &bot_common.RecallStrategy{
					UseNl2sql:  ptr.Of(true),
					UseRerank:  ptr.Of(true),
					UseRewrite: ptr.Of(true),
				},
			},
			Workflow:     []*bot_common.WorkflowInfo{},
			SuggestReply: &bot_common.SuggestReplyInfo{},
			JumpConfig:   &bot_common.JumpConfig{},
			Database:     []*bot_common.Database{},

			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil
}

func (s *SingleAgentApplicationService) defaultModelInfo(ctx context.Context) (*bot_common.ModelInfo, error) {
	model, err := modelmgr.DefaultSVC().GetOnlineDefaultModel(ctx)
	if err != nil {
		return nil, err
	}

	if model == nil {
		return nil, errorx.New(errno.ErrAgentResourceNotFound, errorx.KV("type", "model"), errorx.KV("id", "default"))
	}

	dm := model

	return &bot_common.ModelInfo{
		ModelId:          ptr.Of(dm.ID),
		Temperature:      dm.GetDefaultTemperature(),
		MaxTokens:        dm.GetDefaultMaxTokens(),
		TopP:             dm.GetDefaultTopP(),
		FrequencyPenalty: dm.GetDefaultFrequencyPenalty(),
		PresencePenalty:  dm.GetDefaultPresencePenalty(),
		TopK:             dm.GetDefaultTopK(),
		ModelStyle:       bot_common.ModelStylePtr(bot_common.ModelStyle_Balance),
		ShortMemoryPolicy: &bot_common.ShortMemoryPolicy{
			ContextMode:  bot_common.ContextModePtr(bot_common.ContextMode_FunctionCall_2),
			HistoryRound: ptr.Of[int32](3),
		},
	}, nil
}

func defaultModelRequest(scenario common.Scenario) *manage.ListModelsRequest {
	pageSize := int32(10)
	pageToken := "10"
	return &manage.ListModelsRequest{
		Scenario:  ptr.Of(scenario),
		PageSize:  ptr.Of(pageSize),
		PageToken: ptr.Of(pageToken),
	}
}
