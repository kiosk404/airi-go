package user

import (
	"context"

	"github.com/kiosk404/airi-go/backend/domain/user/repository"
	"github.com/kiosk404/airi-go/backend/domain/user/service"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"golang.org/x/mod/sumdb/storage"
	"gorm.io/gorm"
)

func InitService(ctx context.Context, db *gorm.DB, oss storage.Storage, idgen idgen.IDGenerator) *UserApplicationService {
	UserApplicationSVC.DomainSVC = service.NewUserDomain(ctx, &service.Components{
		IconOSS:  oss,
		IDGen:    idgen,
		UserRepo: repository.NewUserRepo(db),
	})
	UserApplicationSVC.oss = oss

	return UserApplicationSVC
}
