package dao

import (
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/infra/repo/gorm_gen/model"
)

func ModelUserDO2PO(record *entity.User) *model.User {
	return &model.User{
		ID:           record.UserID,
		Name:         record.Name,
		UniqueName:   record.UniqueName,
		Account:      record.Account,
		Description:  record.Description,
		IconURI:      record.IconURI,
		UserVerified: record.UserVerified,
		Locale:       record.Locale,
		SessionKey:   record.SessionKey,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
	}
}
