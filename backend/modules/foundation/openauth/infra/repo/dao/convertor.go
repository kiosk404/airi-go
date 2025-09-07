package dao

import (
	"github.com/kiosk404/airi-go/backend/modules/foundation/openauth/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/foundation/openauth/infra/repo/gorm_gen/model"
)

func ModelUserDO2PO(record *entity.ApiKey) *model.APIKey {
	return &model.APIKey{
		ID:         record.ID,
		APIKey:     record.ApiKey,
		Name:       record.Name,
		UserID:     record.UserID,
		ExpiredAt:  record.ExpiredAt,
		CreatedAt:  record.CreatedAt,
		LastUsedAt: record.LastUsedAt,
		UpdatedAt:  record.UpdatedAt,
	}
}
