package dao

import (
	"context"
	"errors"

	"github.com/kiosk404/airi-go/backend/api/crossdomain/singleagent"
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	"github.com/kiosk404/airi-go/backend/infra/contract/cache"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/infra/repo/gorm_gen/query"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"gorm.io/gorm"
)

type SingleAgentDraftDAO struct {
	idGen       idgen.IDGenerator
	dbQuery     *query.Query
	cacheClient cache.Cmdable
}

func NewSingleAgentDraftDAO(db *gorm.DB, idGen idgen.IDGenerator, cli cache.Cmdable) *SingleAgentDraftDAO {
	return &SingleAgentDraftDAO{
		idGen:       idGen,
		dbQuery:     query.Use(db),
		cacheClient: cli,
	}
}

func (sa *SingleAgentDraftDAO) Create(ctx context.Context, draft *entity.SingleAgent) (draftID int64, err error) {
	id, err := sa.idGen.GenID(ctx)
	if err != nil {
		return 0, errorx.WrapByCode(err, errno.ErrAgentIDGenFailCode, errorx.KV("msg", "CreatePromptResource"))
	}

	return sa.CreateWithID(ctx, id, draft)
}

func (sa *SingleAgentDraftDAO) CreateWithID(ctx context.Context, agentID int64, draft *entity.SingleAgent) (draftID int64, err error) {
	po := sa.singleAgentDraftDo2Po(draft)
	po.AgentID = agentID

	err = sa.dbQuery.SingleAgentDraft.WithContext(ctx).Create(po)
	if err != nil {
		return 0, errorx.WrapByCode(err, errno.ErrAgentCreateDraftCode)
	}

	return agentID, nil
}

func (sa *SingleAgentDraftDAO) Get(ctx context.Context, agentID int64) (*entity.SingleAgent, error) {
	singleAgentDAOModel := sa.dbQuery.SingleAgentDraft
	singleAgent, err := sa.dbQuery.SingleAgentDraft.WithContext(ctx).
		Where(singleAgentDAOModel.AgentID.Eq(agentID)).First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrAgentGetCode)
	}

	do := sa.singleAgentDraftPo2Do(singleAgent)

	return do, nil
}

func (sa *SingleAgentDraftDAO) MGet(ctx context.Context, agentIDs []int64) ([]*entity.SingleAgent, error) {
	sam := sa.dbQuery.SingleAgentDraft
	singleAgents, err := sam.WithContext(ctx).Where(sam.AgentID.In(agentIDs...)).Find()
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrAgentGetCode)
	}

	dos := make([]*entity.SingleAgent, 0, len(singleAgents))
	for _, singleAgent := range singleAgents {
		dos = append(dos, sa.singleAgentDraftPo2Do(singleAgent))
	}

	return dos, nil
}

func (sa *SingleAgentDraftDAO) Save(ctx context.Context, agentInfo *entity.SingleAgent) (err error) {
	po := sa.singleAgentDraftDo2Po(agentInfo)
	singleAgentDAOModel := sa.dbQuery.SingleAgentDraft

	err = singleAgentDAOModel.WithContext(ctx).Where(singleAgentDAOModel.AgentID.Eq(agentInfo.AgentID)).Save(po)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrAgentUpdateCode)
	}

	return nil
}

func (sa *SingleAgentDraftDAO) Update(ctx context.Context, agentInfo *entity.SingleAgent) (err error) {
	po := sa.singleAgentDraftDo2Po(agentInfo)
	singleAgentDAOModel := sa.dbQuery.SingleAgentDraft

	err = singleAgentDAOModel.WithContext(ctx).
		Where(singleAgentDAOModel.AgentID.Eq(agentInfo.AgentID)).Save(po)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrAgentUpdateCode)
	}

	return nil
}

func (sa *SingleAgentDraftDAO) Delete(ctx context.Context, agentID int64) (err error) {
	po := sa.dbQuery.SingleAgentDraft
	_, err = po.WithContext(ctx).Where(po.AgentID.Eq(agentID)).Delete()
	return err
}

func (sa *SingleAgentDraftDAO) singleAgentDraftPo2Do(po *model.SingleAgentDraft) *entity.SingleAgent {
	return &entity.SingleAgent{
		SingleAgent: &singleagent.SingleAgent{
			AgentID:                 po.AgentID,
			Name:                    po.Name,
			Desc:                    ptr.From(po.Description),
			IconURI:                 po.IconURI,
			CreatedAt:               po.CreatedAt,
			UpdatedAt:               po.UpdatedAt,
			DeletedAt:               po.DeletedAt,
			ModelInfo:               po.ModelInfo,
			OnboardingInfo:          po.OnboardingInfo,
			Prompt:                  po.Prompt,
			Plugin:                  po.Plugin,
			Knowledge:               po.Knowledge,
			Workflow:                po.Workflow,
			SuggestReply:            po.SuggestReply,
			JumpConfig:              po.JumpConfig,
			Variables:               po.Variable,
			BackgroundImageInfoList: po.BackgroundImageInfoList,
			Database:                po.DatabaseConfig,
			ShortcutCommand:         po.ShortcutCommand,
			BotMode:                 bot_common.BotMode(po.BotMode),
			LayoutInfo:              po.LayoutInfo,
		},
	}
}

func (sa *SingleAgentDraftDAO) singleAgentDraftDo2Po(do *entity.SingleAgent) *model.SingleAgentDraft {
	return &model.SingleAgentDraft{
		AgentID:                 do.AgentID,
		Name:                    do.Name,
		Description:             ptr.Of(do.Desc),
		IconURI:                 do.IconURI,
		CreatedAt:               do.CreatedAt,
		UpdatedAt:               do.UpdatedAt,
		DeletedAt:               do.DeletedAt,
		ModelInfo:               do.ModelInfo,
		OnboardingInfo:          do.OnboardingInfo,
		Prompt:                  do.Prompt,
		Plugin:                  do.Plugin,
		Knowledge:               do.Knowledge,
		Workflow:                do.Workflow,
		SuggestReply:            do.SuggestReply,
		JumpConfig:              do.JumpConfig,
		Variable:                do.Variables,
		BackgroundImageInfoList: do.BackgroundImageInfoList,
		DatabaseConfig:          do.Database,
		ShortcutCommand:         do.ShortcutCommand,
		BotMode:                 int32(do.BotMode),
		LayoutInfo:              do.LayoutInfo,
	}
}
