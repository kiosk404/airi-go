package dao

import (
	"context"
	"errors"

	"github.com/kiosk404/airi-go/backend/api/crossdomain/singleagent"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/infra/repo/gorm_gen/query"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"gorm.io/gorm"
)

type SingleAgentVersionDAO struct {
	IDGen   idgen.IDGenerator
	dbQuery *query.Query
}

func NewSingleAgentVersion(db *gorm.DB, idGen idgen.IDGenerator) *SingleAgentVersionDAO {
	return &SingleAgentVersionDAO{
		IDGen:   idGen,
		dbQuery: query.Use(db),
	}
}

func (sa *SingleAgentVersionDAO) GetLatest(ctx context.Context, agentID int64) (*entity.SingleAgent, error) {
	singleAgentDAOModel := sa.dbQuery.SingleAgentVersion
	singleAgent, err := singleAgentDAOModel.WithContext(ctx).
		Where(singleAgentDAOModel.AgentID.Eq(agentID)).
		Order(singleAgentDAOModel.CreatedAt.Desc()).
		First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrAgentGetCode)
	}

	do := sa.singleAgentVersionPo2Do(singleAgent)

	return do, nil
}

func (sa *SingleAgentVersionDAO) Get(ctx context.Context, agentID int64, version string) (*entity.SingleAgent, error) {
	singleAgentDAOModel := sa.dbQuery.SingleAgentVersion
	singleAgent, err := singleAgentDAOModel.WithContext(ctx).
		Where(singleAgentDAOModel.AgentID.Eq(agentID), singleAgentDAOModel.Version.Eq(version)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrAgentGetCode)
	}

	do := sa.singleAgentVersionPo2Do(singleAgent)

	return do, nil
}

func (sa *SingleAgentVersionDAO) singleAgentVersionPo2Do(po *model.SingleAgentVersion) *entity.SingleAgent {
	return &entity.SingleAgent{
		SingleAgent: &singleagent.SingleAgent{
			AgentID:         po.AgentID,
			Name:            po.Name,
			Desc:            ptr.From(po.Description),
			IconURI:         po.IconURI,
			CreatedAt:       po.CreatedAt,
			UpdatedAt:       po.UpdatedAt,
			DeletedAt:       po.DeletedAt,
			ModelInfo:       po.ModelInfo,
			OnboardingInfo:  po.OnboardingInfo,
			Prompt:          po.Prompt,
			Plugin:          po.Plugin,
			Knowledge:       po.Knowledge,
			Workflow:        po.Workflow,
			SuggestReply:    po.SuggestReply,
			JumpConfig:      po.JumpConfig,
			Variables:       po.Variable,
			Database:        po.DatabaseConfig,
			ShortcutCommand: po.ShortcutCommand,
			Version:         po.Version,
		},
	}
}

func (sa *SingleAgentVersionDAO) singleAgentVersionDo2Po(do *entity.SingleAgent) *model.SingleAgentVersion {
	return &model.SingleAgentVersion{
		AgentID:         do.AgentID,
		Name:            do.Name,
		Description:     ptr.Of(do.Desc),
		IconURI:         do.IconURI,
		CreatedAt:       do.CreatedAt,
		UpdatedAt:       do.UpdatedAt,
		DeletedAt:       do.DeletedAt,
		ModelInfo:       do.ModelInfo,
		OnboardingInfo:  do.OnboardingInfo,
		Prompt:          do.Prompt,
		Plugin:          do.Plugin,
		Knowledge:       do.Knowledge,
		Workflow:        do.Workflow,
		SuggestReply:    do.SuggestReply,
		JumpConfig:      do.JumpConfig,
		Variable:        do.Variables,
		DatabaseConfig:  do.Database,
		ShortcutCommand: do.ShortcutCommand,
	}
}
