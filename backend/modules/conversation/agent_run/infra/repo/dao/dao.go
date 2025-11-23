package dao

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/infra/repo/gorm_gen/query"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/pkg"
	agentrun "github.com/kiosk404/airi-go/backend/modules/conversation/crossdomain/agentrun/model"
	"github.com/kiosk404/airi-go/backend/pkg/json"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"gorm.io/gorm"
)

type RunRecordDAO struct {
	db    *gorm.DB
	query *query.Query
	idGen idgen.IDGenerator
}

func NewRunRecordDAO(db *gorm.DB, idGen idgen.IDGenerator) *RunRecordDAO {
	return &RunRecordDAO{
		db:    db,
		idGen: idGen,
		query: query.Use(db),
	}
}

func (dao *RunRecordDAO) Create(ctx context.Context, runMeta *entity.AgentRunMeta) (*entity.RunRecordMeta, error) {

	createPO, err := dao.buildCreatePO(ctx, runMeta)
	if err != nil {
		return nil, err
	}

	createErr := dao.query.RunRecord.WithContext(ctx).Create(createPO)
	if createErr != nil {
		return nil, createErr
	}

	return dao.buildPo2Do(createPO), nil
}

func (dao *RunRecordDAO) GetByID(ctx context.Context, id int64) (*entity.RunRecordMeta, error) {
	po, err := dao.query.RunRecord.WithContext(ctx).Where(dao.query.RunRecord.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return dao.buildPo2Do(po), nil
}

func (dao *RunRecordDAO) UpdateByID(ctx context.Context, id int64, updateMeta *entity.UpdateMeta) error {
	po := &model.RunRecord{
		ID: id,
	}
	if updateMeta.Status != "" {

		po.Status = string(updateMeta.Status)
	}
	if updateMeta.LastError != nil {
		errString, err := json.Marshal(updateMeta.LastError)
		if err != nil {
			return err
		}
		po.LastError = ptr.Of(string(errString))
	}
	if updateMeta.CompletedAt != 0 {

		po.CompletedAt = updateMeta.CompletedAt
	}
	if updateMeta.FailedAt != 0 {

		po.FailedAt = updateMeta.FailedAt
	}
	if updateMeta.Usage != nil {
		po.Usage = ptr.PtrConvert[agentrun.Usage, entity.Usage](updateMeta.Usage)
	}
	po.UpdatedAt = time.Now().UnixMilli()

	_, err := dao.query.RunRecord.WithContext(ctx).Where(dao.query.RunRecord.ID.Eq(id)).Updates(po)
	return err
}

func (dao *RunRecordDAO) Delete(ctx context.Context, id []int64) error {

	_, err := dao.query.RunRecord.WithContext(ctx).Where(dao.query.RunRecord.ID.In(id...)).UpdateColumns(map[string]interface{}{
		"updated_at": time.Now().UnixMilli(),
		"status":     entity.RunStatusDeleted,
	})

	return err
}

func (dao *RunRecordDAO) List(ctx context.Context, meta *entity.ListRunRecordMeta) ([]*entity.RunRecordMeta, error) {
	logs.InfoX(pkg.ModelName, "list run record req:%v, sectionID:%v, limit:%v", meta.ConversationID, meta.SectionID, meta.Limit)
	m := dao.query.RunRecord
	mq := m.WithContext(ctx)
	do := m.WithContext(ctx).Where(m.ConversationID.Eq(meta.ConversationID)).Debug().Where(m.Status.NotIn(string(entity.RunStatusDeleted)))
	if meta.BeforeID > 0 {
		runRecord, err := mq.Where(m.ID.Eq(meta.BeforeID)).First()
		if err != nil {
			return nil, err
		}
		do = do.Where(m.CreatedAt.Lt(runRecord.CreatedAt))
	}
	if meta.AfterID > 0 {
		runRecord, err := mq.Where(m.ID.Eq(meta.AfterID)).First()
		if err != nil {
			return nil, err
		}
		do = do.Where(m.CreatedAt.Gt(runRecord.CreatedAt))
	}
	if meta.SectionID > 0 {
		do = do.Where(m.SectionID.Eq(meta.SectionID))
	}
	if meta.Limit > 0 {
		do = do.Limit(int(meta.Limit))
	}
	if strings.ToLower(meta.OrderBy) == "asc" {
		do = do.Order(m.CreatedAt.Asc())
	} else {
		do = do.Order(m.CreatedAt.Desc())
	}

	runRecords, err := do.Find()
	return slices.Transform(runRecords, func(item *model.RunRecord) *entity.RunRecordMeta {
		return dao.buildPo2Do(item)
	}), err
}

func (dao *RunRecordDAO) buildCreatePO(ctx context.Context, runMeta *entity.AgentRunMeta) (*model.RunRecord, error) {

	runID, err := dao.idGen.GenID(ctx)

	if err != nil {
		return nil, err
	}
	reqOrigin, err := json.Marshal(runMeta)
	if err != nil {
		return nil, err
	}

	timeNow := time.Now().UnixMilli()
	creatorID, err := strconv.ParseInt(runMeta.UserID, 10, 64)
	if err != nil {
		return nil, err
	}
	return &model.RunRecord{
		ID:             runID,
		ConversationID: runMeta.ConversationID,
		SectionID:      runMeta.SectionID,
		AgentID:        runMeta.AgentID,
		Status:         string(entity.RunStatusCreated),
		ChatRequest:    ptr.Of(string(reqOrigin)),
		UserID:         runMeta.UserID,
		CreatedAt:      timeNow,
		CreatorID:      creatorID,
	}, nil
}

func (dao *RunRecordDAO) buildPo2Do(po *model.RunRecord) *entity.RunRecordMeta {
	runMeta := &entity.RunRecordMeta{
		ID:             po.ID,
		ConversationID: po.ConversationID,
		SectionID:      po.SectionID,
		AgentID:        po.AgentID,
		Status:         entity.RunStatus(po.Status),
		Ext:            ptr.From(po.Ext),
		CreatedAt:      po.CreatedAt,
		UpdatedAt:      po.UpdatedAt,
		CompletedAt:    po.CompletedAt,
		FailedAt:       po.FailedAt,
		Usage:          ptr.PtrConvert[entity.Usage, agentrun.Usage](po.Usage),
		CreatorID:      po.CreatorID,
	}

	return runMeta
}

func (dao *RunRecordDAO) Cancel(ctx context.Context, meta *entity.CancelRunMeta) (*entity.RunRecordMeta, error) {

	m := dao.query.RunRecord
	_, err := m.WithContext(ctx).Where(m.ID.Eq(meta.RunID)).UpdateColumns(map[string]interface{}{
		"updated_at": time.Now().UnixMilli(),
		"status":     entity.RunEventCancelled,
	})
	if err != nil {
		return nil, err
	}
	return dao.GetByID(ctx, meta.RunID)
}
