package dao

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/data/upload/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/data/upload/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/modules/data/upload/infra/repo/gorm_gen/query"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
	"gorm.io/gorm"
)

type FilesDAO struct {
	DB    *gorm.DB
	Query *query.Query
}

func NewFilesDAO(db *gorm.DB) *FilesDAO {
	return &FilesDAO{
		DB:    db,
		Query: query.Use(db),
	}
}

func (dao *FilesDAO) Create(ctx context.Context, file *entity.File) error {
	f := dao.fromEntityToModel(file)
	return dao.Query.File.WithContext(ctx).Create(f)
}

func (dao *FilesDAO) BatchCreate(ctx context.Context, files []*entity.File) error {
	if len(files) == 0 {
		return nil
	}
	return dao.Query.File.WithContext(ctx).CreateInBatches(slices.Transform(files, dao.fromEntityToModel), len(files))
}

func (dao *FilesDAO) Delete(ctx context.Context, id int64) error {
	_, err := dao.Query.File.WithContext(ctx).Where(dao.Query.File.ID.Eq(id)).Delete()
	return err
}

func (dao *FilesDAO) GetByID(ctx context.Context, id int64) (*entity.File, error) {
	file, err := dao.Query.File.WithContext(ctx).Where(dao.Query.File.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return dao.fromModelToEntity(file), nil
}

func (dao *FilesDAO) MGetByIDs(ctx context.Context, ids []int64) ([]*entity.File, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	files, err := dao.Query.File.WithContext(ctx).Where(dao.Query.File.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return slices.Transform(files, dao.fromModelToEntity), nil
}

func (dao *FilesDAO) fromModelToEntity(model *model.File) *entity.File {
	if model == nil {
		return nil
	}
	return &entity.File{
		ID:          model.ID,
		Name:        model.Name,
		FileSize:    model.FileSize,
		TosURI:      model.TosURI,
		Status:      entity.FileStatus(model.Status),
		Comment:     model.Comment,
		Source:      entity.FileSource(model.Source),
		CreatorID:   model.CreatorID,
		ContentType: model.ContentType,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func (dao *FilesDAO) fromEntityToModel(entity *entity.File) *model.File {
	return &model.File{
		ID:          entity.ID,
		Name:        entity.Name,
		FileSize:    entity.FileSize,
		TosURI:      entity.TosURI,
		Status:      int32(entity.Status),
		Comment:     entity.Comment,
		Source:      int32(entity.Source),
		CreatorID:   entity.CreatorID,
		ContentType: entity.ContentType,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
