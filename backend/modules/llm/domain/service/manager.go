package service

import (
	"context"
	"fmt"

	"github.com/kiosk404/airi-go/backend/modules/llm/domain/component/conf"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	llm_errorx "github.com/kiosk404/airi-go/backend/modules/llm/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=mocks/manage.go -package=mocks . IManage
type IManage interface {
	ListModels(ctx context.Context, req entity.ListModelsRequest) (models []*entity.Model, total int64, hasMore bool, nextPageToken int64, err error)
	GetModelByID(ctx context.Context, id int64) (model *entity.Model, err error)
	MGetModelByID(ctx context.Context, req *entity.MGetModelReq) ([]*entity.Model, error)
}

type ManageImpl struct {
	conf conf.IConfigManage
}

var _ IManage = (*ManageImpl)(nil)

func (m *ManageImpl) ListModels(ctx context.Context, req entity.ListModelsRequest) (models []*entity.Model, total int64, hasMore bool, nextPageToken int64, err error) {
	return m.conf.ListModels(ctx, req)
}

func (m *ManageImpl) GetModelByID(ctx context.Context, id int64) (model *entity.Model, err error) {
	model, err = m.conf.GetModel(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewByCode(llm_errorx.ResourceNotFoundCode, errorx.WithExtraMsg(fmt.Sprintf("model id:%d not exist in db", id)))
		}
		return nil, errorx.NewByCode(llm_errorx.CommonMySqlErrorCode, errorx.WithExtraMsg(err.Error()))
	}
	return model, nil
}

func (m *ManageImpl) MGetModelByID(ctx context.Context, req *entity.MGetModelReq) ([]*entity.Model, error) {
	resp := make([]*entity.Model, 0, len(req.ModelIDs))
	modelsMap, err := m.conf.GetModelSet(ctx)
	if err != nil {
		return nil, err
	}
	for _, id := range req.ModelIDs {
		if md, found := modelsMap[id]; found {
			resp = append(resp, md)
		}
	}
	return resp, nil
}
