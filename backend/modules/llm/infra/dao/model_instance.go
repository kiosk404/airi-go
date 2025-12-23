package dao

import (
	"context"
	"time"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	instancemodle "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/repo/gorm_gen/query"
	"github.com/kiosk404/airi-go/backend/modules/llm/pkg"
	"github.com/kiosk404/airi-go/backend/modules/llm/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/json"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"gorm.io/gorm"
)

type ModelMgrDao struct {
	dbQuery *query.Query
	IDGen   idgen.IDGenerator
}

func NewModelMgrDao(db *gorm.DB, idGen idgen.IDGenerator) *ModelMgrDao {
	return &ModelMgrDao{
		IDGen:   idGen,
		dbQuery: query.Use(db),
	}
}

func (m *ModelMgrDao) CreateModel(ctx context.Context, instance *entity.ModelInstance) (id int64, err error) {
	mgrDao := m.dbQuery.ModelInstance

	id, err = m.IDGen.GenID(ctx)
	if err != nil {
		return 0, errorx.WrapByCode(err, errno.ErrModelIDGenFailCode, errorx.KV("msg", "CreateProjectVariable"))
	}

	po := m.modelInstancePo2Do(instance)
	po.ID = id

	err = mgrDao.WithContext(ctx).Create(po)

	return id, err
}

func (m *ModelMgrDao) GetModel(ctx context.Context, id int64) (do *entity.ModelInstance, err error) {
	mgrDao := m.dbQuery.ModelInstance

	po, err := mgrDao.WithContext(ctx).Where(mgrDao.ID.Eq(id)).First()
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrModelNotFoundCode, errorx.KV("msg", "GetModel"))
	}

	return m.modelInstanceDo2Po(po), err
}

func (m *ModelMgrDao) ListModels(ctx context.Context) (do []*entity.ModelInstance, err error) {
	mgrDao := m.dbQuery.ModelInstance

	pos, err := mgrDao.WithContext(ctx).Find()
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrModelNotFoundCode, errorx.KV("msg", "ListModels"))
	}

	do = make([]*entity.ModelInstance, 0, len(pos))
	for _, po := range pos {
		do = append(do, m.modelInstanceDo2Po(po))
	}

	return do, nil
}

func (m *ModelMgrDao) DeleteModel(ctx context.Context, id int64) (err error) {
	mgrDao := m.dbQuery.ModelInstance

	result, err := mgrDao.WithContext(ctx).Where(mgrDao.ID.Eq(id)).Delete()
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrModelNotFoundCode, errorx.KV("msg", "DeleteModel"))
	}
	logs.InfoX(pkg.ModelName, "DeleteModel result effected: %v", result.RowsAffected)

	return nil
}

func (m *ModelMgrDao) UpdateModel(ctx context.Context, instance *entity.ModelInstance) (err error) {
	mgrDao := m.dbQuery.ModelInstance

	po := m.modelInstancePo2Do(instance)

	_, err = mgrDao.WithContext(ctx).Where(mgrDao.ID.Eq(instance.ID)).Updates(po)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrModelNotFoundCode, errorx.KV("msg", "UpdateModel"))
	}

	return nil
}

func (m *ModelMgrDao) ListModelByType(ctx context.Context, modelType entity.ModelType, limit int) (do []*entity.ModelInstance, err error) {
	mgrDao := m.dbQuery.ModelInstance
	if limit <= 0 {
		limit = 10
	}

	pos, err := mgrDao.WithContext(ctx).Where(mgrDao.Type.Eq(modelType.Int32())).Limit(limit).Find()
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrModelNotFoundCode, errorx.KV("msg", "ListModelByType"))
	}

	do = make([]*entity.ModelInstance, 0, len(pos))
	for _, po := range pos {
		do = append(do, m.modelInstanceDo2Po(po))
	}

	return do, nil
}

func (m *ModelMgrDao) modelInstancePo2Do(do *entity.ModelInstance) *model.ModelInstance {
	if do == nil {
		return nil
	}

	extraStr := "{}"
	if do.Extra.ModelExtra != nil {
		extraByte, err := json.Marshal(do.Extra.ModelExtra)
		if err != nil {
			logs.ErrorX(pkg.ModelName, "marshal extra failed, err: %w", err)
		}

		extraStr = string(extraByte)
	}

	return &model.ModelInstance{
		ID:          do.ID,
		Type:        do.Type.ModelType.Int32(),
		Provider:    ptr.Of(do.Provider),
		DisplayInfo: ptr.Of(do.DisplayInfo),
		Connection:  do.Connection.Connection,
		Capability:  ptr.Of(do.Capability),
		Parameters:  do.Parameters,
		Extra:       ptr.Of(extraStr),
	}
}

func (m *ModelMgrDao) modelInstanceDo2Po(do *model.ModelInstance) *entity.ModelInstance {
	if do == nil {
		return nil
	}

	return &entity.ModelInstance{
		ID:   do.ID,
		Type: entity.ModelType{ModelType: ptr.Of(instancemodle.ModelType(do.Type))},
		Provider: instancemodle.ModelProvider{
			Name:        do.Provider.Name,
			IconURI:     do.Provider.IconURI,
			IconURL:     do.Provider.IconURL,
			Description: do.Provider.Description,
			ModelClass:  do.Provider.ModelClass,
		},
		DisplayInfo: instancemodle.DisplayInfo{
			Name:         do.DisplayInfo.Name,
			Description:  do.DisplayInfo.Description,
			OutputTokens: do.DisplayInfo.OutputTokens,
			MaxTokens:    do.DisplayInfo.MaxTokens,
		},
		Connection: entity.Connection{Connection: do.Connection},
		Capability: instancemodle.ModelAbility{
			CotDisplay:         do.Capability.CotDisplay,
			FunctionCall:       do.Capability.FunctionCall,
			ImageUnderstanding: do.Capability.ImageUnderstanding,
			VideoUnderstanding: do.Capability.VideoUnderstanding,
			AudioUnderstanding: do.Capability.AudioUnderstanding,
			SupportMultiModal:  do.Capability.SupportMultiModal,
			PrefillResp:        do.Capability.PrefillResp,
		},
		Parameters: do.Parameters,
		Extra:      entity.ModelExtra{ModelExtra: ptr.Of(ModelExtra(do.Extra))},
		CreatedAt:  time.Unix(do.CreatedAt, 0),
		UpdatedAt:  time.Unix(do.UpdatedAt, 0),
	}
}

func ModelExtra(do *string) instancemodle.ModelExtra {
	var extra instancemodle.ModelExtra
	if err := json.Unmarshal([]byte(*do), &extra); err != nil {
		return instancemodle.ModelExtra{}
	}
	return extra
}
