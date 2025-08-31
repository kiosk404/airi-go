package openapiauth

import (
	"context"
	"time"

	"github.com/kiosk404/airi-go/backend/domain/openauth/openapiauth/entity"
	"github.com/kiosk404/airi-go/backend/domain/openauth/openapiauth/internal/dal"
	"github.com/kiosk404/airi-go/backend/domain/openauth/openapiauth/internal/dal/model"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"gorm.io/gorm"
)

type apiAuthImpl struct {
	IDGen idgen.IDGenerator
	DB    *gorm.DB
	dao   *dal.ApiKeyDAO
}

type Components struct {
	IDGen idgen.IDGenerator
	DB    *gorm.DB
}

func NewService(c *Components) APIAuth {
	return &apiAuthImpl{
		IDGen: c.IDGen,
		DB:    c.DB,
		dao:   dal.NewApiKeyDAO(c.IDGen, c.DB),
	}
}

func (a *apiAuthImpl) Create(ctx context.Context, req *entity.CreateApiKey) (*entity.ApiKey, error) {
	apiKeyData, err := a.dao.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return apiKeyData, nil
}

func (a *apiAuthImpl) Delete(ctx context.Context, req *entity.DeleteApiKey) error {

	return a.dao.Delete(ctx, req.ID, req.UserID)

}
func (a *apiAuthImpl) Get(ctx context.Context, req *entity.GetApiKey) (*entity.ApiKey, error) {

	apiKey, err := a.dao.Get(ctx, req.ID)
	logs.Info("apiKey=%v, err:%v", apiKey, err)
	if err != nil {
		return nil, err
	}
	if apiKey == nil {
		return nil, nil
	}
	return a.buildPoData2ApiKey([]*model.APIKey{apiKey})[0], nil
}

func (a *apiAuthImpl) buildPoData2ApiKey(apiKey []*model.APIKey) []*entity.ApiKey {

	apiKeyData := slices.Transform(apiKey, func(a *model.APIKey) *entity.ApiKey {
		return &entity.ApiKey{
			ID:         a.ID,
			Name:       a.Name,
			ApiKey:     a.APIKey,
			UserID:     a.UserID,
			ExpiredAt:  a.ExpiredAt,
			CreatedAt:  a.CreatedAt,
			LastUsedAt: a.LastUsedAt,
		}
	})

	return apiKeyData
}

func (a *apiAuthImpl) List(ctx context.Context, req *entity.ListApiKey) (*entity.ListApiKeyResp, error) {
	resp := &entity.ListApiKeyResp{
		ApiKeys: make([]*entity.ApiKey, 0),
		HasMore: false,
	}
	apiKey, hasMore, err := a.dao.List(ctx, req.UserID, int(req.Limit), int(req.Page))
	if err != nil {
		return nil, err
	}
	resp.ApiKeys = a.buildPoData2ApiKey(apiKey)
	resp.HasMore = hasMore

	return resp, nil
}
func (a *apiAuthImpl) CheckPermission(ctx context.Context, req *entity.CheckPermission) (*entity.ApiKey, error) {

	apiKey, err := a.dao.FindByKey(ctx, req.ApiKey)
	if err != nil {
		return nil, err
	}
	if apiKey.APIKey != req.ApiKey {
		return nil, nil
	}
	apiKeyDo := &entity.ApiKey{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		UserID:    apiKey.UserID,
		ExpiredAt: apiKey.ExpiredAt,
		CreatedAt: apiKey.CreatedAt,
	}
	return apiKeyDo, nil
}

func (a *apiAuthImpl) Save(ctx context.Context, sm *entity.SaveMeta) error {

	updateColumn := make(map[string]any)
	if sm.Name != nil {
		updateColumn["name"] = sm.Name
	}
	if sm.LastUsedAt != nil {
		updateColumn["last_used_at"] = sm.LastUsedAt
	}
	updateColumn["updated_at"] = time.Now().Unix()
	err := a.dao.Update(ctx, sm.ID, sm.UserID, updateColumn)

	return err
}
