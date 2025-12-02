package dao

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/infra/repo/gorm_gen/query"
	"gorm.io/gorm"
)

type ModelRunRecordDao struct {
	dbQuery *query.Query
}

func NewModelRunRecordDao(db *gorm.DB) *ModelRunRecordDao {
	return &ModelRunRecordDao{
		dbQuery: query.Use(db),
	}
}

func (mrd *ModelRunRecordDao) Create(ctx context.Context, runRecord *entity.ModelRequestRecord) (err error) {
	modelRunRecordDaoModel := mrd.dbQuery.ModelRequestRecord
	do := modelRunRecordDaoModel.WithContext(ctx).Debug()
	po := mrd.modelRunRecordDo2Po(runRecord)
	return do.Create(po)
}

func (mrd *ModelRunRecordDao) List(ctx context.Context, modelID string, limit int) (modelRunRecordList []*entity.ModelRequestRecord) {
	modelRunRecordDaoModel := mrd.dbQuery.ModelRequestRecord
	do := modelRunRecordDaoModel.WithContext(ctx).Debug()
	if limit <= 0 {
		limit = 1
	}
	poList, err := do.Where(modelRunRecordDaoModel.ModelID.Eq(modelID)).
		Order(modelRunRecordDaoModel.ID.Desc()).Limit(limit).Find()
	if err != nil {
		return nil
	}
	for _, po := range poList {
		modelRunRecordList = append(modelRunRecordList, mrd.modelRunRecordPo2Do(po))
	}
	return modelRunRecordList
}

func (mrd *ModelRunRecordDao) modelRunRecordPo2Do(runRecord *model.ModelRequestRecord) *entity.ModelRequestRecord {
	return &entity.ModelRequestRecord{
		ID:                  runRecord.ID,
		UserID:              runRecord.UserID,
		UsageScene:          runRecord.UsageScene,
		UsageSceneEntityID:  runRecord.UsageSceneEntityID,
		Protocol:            runRecord.Protocol,
		ModelIdentification: runRecord.ModelIdentification,
		ModelAk:             runRecord.ModelAk,
		ModelID:             runRecord.ModelID,
		ModelName:           runRecord.ModelName,
		InputToken:          runRecord.InputToken,
		OutputToken:         runRecord.OutputToken,
		LogId:               runRecord.Logid,
		ErrorCode:           runRecord.ErrorCode,
		ErrorMsg:            runRecord.ErrorMsg,
		CreatedAt:           runRecord.CreatedAt,
		UpdatedAt:           runRecord.UpdatedAt,
	}
}

func (mrd *ModelRunRecordDao) modelRunRecordDo2Po(runRecord *entity.ModelRequestRecord) *model.ModelRequestRecord {
	return &model.ModelRequestRecord{
		ID:                  runRecord.ID,
		UserID:              runRecord.UserID,
		UsageScene:          runRecord.UsageScene,
		UsageSceneEntityID:  runRecord.UsageSceneEntityID,
		Protocol:            runRecord.Protocol,
		ModelIdentification: runRecord.ModelIdentification,
		ModelAk:             runRecord.ModelAk,
		ModelID:             runRecord.ModelID,
		ModelName:           runRecord.ModelName,
		InputToken:          runRecord.InputToken,
		OutputToken:         runRecord.OutputToken,
		Logid:               runRecord.LogId,
		ErrorCode:           runRecord.ErrorCode,
		ErrorMsg:            runRecord.ErrorMsg,
		UpdatedAt:           runRecord.UpdatedAt,
		CreatedAt:           runRecord.CreatedAt,
	}
}
