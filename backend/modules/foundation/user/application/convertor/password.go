package convertor

import (
	domain "github.com/kiosk404/airi-go/backend/api/model/foundation/domain/user"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

func UserDO2DTO(do *entity.User) *domain.UserInfoDetail {
	if do == nil {
		return nil
	}
	return &domain.UserInfoDetail{
		Name:      do.UniqueName,
		NickName:  do.Name,
		AvatarURL: do.IconURL,
		Account:   do.Account,
		Mobile:    nil,
		UserID:    ptr.Of(conv.ToString(do.UserID)),
	}
}
