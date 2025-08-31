package appinfra

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/cache"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"gorm.io/gorm"
)

type AppDependencies struct {
	DB       *gorm.DB
	CacheCli cache.Cmdable
	IDGenSVC idgen.IDGenerator
}

func Init(ctx context.Context) (*AppDependencies, error) {
	deps := &AppDependencies{}
	var err error

	return deps, err
}
