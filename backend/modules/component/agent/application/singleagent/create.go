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
	"github.com/kiosk404/airi-go/backend/modules/llm/application/convertor"
	moduleentity "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

func (s *SingleAgentApplicationService) CreateSingleAgentDraft(ctx context.Context, req *developer_api.DraftBotCreateRequest) (*developer_api.DraftBotCreateResponse, error) {
	resp, err := s.appContext.ModelMgr.DomainSVC.ListModels(ctx, defaultModelRequest(common.ScenarioDefault))
	if err != nil {
		return nil, err
	}
	if len(resp.Models) == 0 {
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
	// todo: 事件推送

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
	modelResp, err := s.appContext.ModelMgr.DomainSVC.ListModels(ctx, defaultModelRequest(common.ScenarioDefault))
	if err != nil {
		return nil, err
	}

	if len(modelResp.Models) == 0 {
		return nil, errorx.New(errno.ErrAgentResourceNotFound, errorx.KV("type", "model"), errorx.KV("id", "default"))
	}

	dm := modelResp.Models[0]
	dme := convertor.ModelDTO2DO(dm)

	var temperature *float64
	if tp, ok := dme.FindParameter(moduleentity.Temperature); ok {
		t, err := tp.GetFloat(moduleentity.DefaultTypeBalance)
		if err != nil {
			return nil, err
		}

		temperature = ptr.Of(t)
	}

	var maxTokens *int32
	if tp, ok := dme.FindParameter(moduleentity.MaxTokens); ok {
		t, err := tp.GetInt(moduleentity.DefaultTypeBalance)
		if err != nil {
			return nil, err
		}
		maxTokens = ptr.Of(int32(t))
	} else if dme.CommonParam.MaxTokens != nil {
		maxTokens = ptr.Of(int32(*dme.CommonParam.MaxTokens))
	}

	var topP *float64
	if tp, ok := dme.FindParameter(moduleentity.TopP); ok {
		t, err := tp.GetFloat(moduleentity.DefaultTypeBalance)
		if err != nil {
			return nil, err
		}
		topP = ptr.Of(t)
	}

	var topK *int32
	if tp, ok := dme.FindParameter(moduleentity.TopK); ok {
		t, err := tp.GetInt(moduleentity.DefaultTypeBalance)
		if err != nil {
			return nil, err
		}
		topK = ptr.Of(int32(t))
	}

	var frequencyPenalty *float64
	if tp, ok := dme.FindParameter(moduleentity.FrequencyPenalty); ok {
		t, err := tp.GetFloat(moduleentity.DefaultTypeBalance)
		if err != nil {
			return nil, err
		}
		frequencyPenalty = ptr.Of(t)
	}

	var presencePenalty *float64
	if tp, ok := dme.FindParameter(moduleentity.PresencePenalty); ok {
		t, err := tp.GetFloat(moduleentity.DefaultTypeBalance)
		if err != nil {
			return nil, err
		}
		presencePenalty = ptr.Of(t)
	}

	return &bot_common.ModelInfo{
		ModelId:          dm.ModelID,
		Temperature:      temperature,
		MaxTokens:        maxTokens,
		TopP:             topP,
		FrequencyPenalty: frequencyPenalty,
		PresencePenalty:  presencePenalty,
		TopK:             topK,
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
