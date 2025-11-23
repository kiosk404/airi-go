package dao

import (
	"context"
	"errors"
	"time"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/modules/component/prompt/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/component/prompt/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/modules/component/prompt/infra/repo/gorm_gen/query"
	"github.com/kiosk404/airi-go/backend/modules/component/prompt/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type PromptDAO struct {
	IDGen   idgen.IDGenerator
	dbQuery *query.Query
}

func NewPromptDAO(db *gorm.DB, generator idgen.IDGenerator) *PromptDAO {
	return &PromptDAO{
		IDGen:   generator,
		dbQuery: query.Use(db),
	}
}

func (d *PromptDAO) promptResourceDO2PO(p *entity.PromptResource) *model.PromptResource {
	return &model.PromptResource{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		PromptText:  ptr.Of(p.PromptText),
		Status:      p.Status,
		CreatorID:   p.CreatorID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func (d *PromptDAO) promptResourcePO2DO(p *model.PromptResource) *entity.PromptResource {
	return &entity.PromptResource{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		PromptText:  ptr.From(p.PromptText),
		Status:      p.Status,
		CreatorID:   p.CreatorID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func (d *PromptDAO) CreatePromptResource(ctx context.Context, do *entity.PromptResource) (int64, error) {
	id, err := d.IDGen.GenID(ctx)
	if err != nil {
		return 0, errorx.New(errno.ErrPromptIDGenFailCode, errorx.KV("msg", "CreatePromptResource"))
	}

	p := d.promptResourceDO2PO(do)

	now := time.Now().Unix()

	p.ID = id
	p.Status = 1
	p.CreatedAt = now
	p.UpdatedAt = now

	promptModel := d.dbQuery.PromptResource
	err = promptModel.WithContext(ctx).Create(p)
	if err != nil {
		return 0, errorx.WrapByCode(err, errno.ErrPromptCreateCode)
	}

	return id, nil
}

func (d *PromptDAO) GetPromptResource(ctx context.Context, promptID int64) (*entity.PromptResource, error) {
	promptModel := d.dbQuery.PromptResource
	promptWhere := []gen.Condition{
		promptModel.ID.Eq(promptID),
	}

	promptResource, err := promptModel.WithContext(ctx).Where(promptWhere...).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.WrapByCode(err, errno.ErrPromptDataNotFoundCode)
	}

	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPromptGetCode)
	}

	do := d.promptResourcePO2DO(promptResource)

	return do, nil
}

func (d *PromptDAO) UpdatePromptResource(ctx context.Context, promptID int64, name, description, promptText *string) error {
	updateMap := make(map[string]any, 5)

	if name != nil {
		updateMap["name"] = *name
	}

	if description != nil {
		updateMap["description"] = *description
	}

	if promptText != nil {
		updateMap["prompt_text"] = *promptText
	}

	promptModel := d.dbQuery.PromptResource
	promptWhere := []gen.Condition{
		promptModel.ID.Eq(promptID),
	}

	_, err := promptModel.WithContext(ctx).Where(promptWhere...).Updates(updateMap)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPromptUpdateCode)
	}

	return nil
}

func (d *PromptDAO) DeletePromptResource(ctx context.Context, ID int64) error {
	promptModel := d.dbQuery.PromptResource
	promptWhere := []gen.Condition{
		promptModel.ID.Eq(ID),
	}
	_, err := promptModel.WithContext(ctx).Where(promptWhere...).Delete()
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPromptDeleteCode)
	}

	return nil
}
