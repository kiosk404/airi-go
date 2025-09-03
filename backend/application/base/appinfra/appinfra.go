package appinfra

import (
	"context"
	"fmt"

	"github.com/kiosk404/airi-go/backend/infra/contract/cache"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/impl/cache/local"
	idgenimpl "github.com/kiosk404/airi-go/backend/infra/impl/idgen"
	"github.com/kiosk404/airi-go/backend/infra/impl/rdb/sqlite"
	"github.com/kiosk404/airi-go/backend/infra/impl/storage"
	"gorm.io/gorm"
)

type AppDependencies struct {
	DB        *gorm.DB
	CacheCli  cache.Cmdable
	IDGenSVC  idgen.IDGenerator
	TOSClient storage.Storage
}

func Init(ctx context.Context) (*AppDependencies, error) {
	deps := &AppDependencies{}
	var err error
	if deps.DB, err = sqlite.New(); err != nil {
		return nil, fmt.Errorf("init db failed, err=%w", err)
	}
	if deps.CacheCli, err = local.New(); err != nil {
		return nil, fmt.Errorf("init cache failed, err=%w", err)
	}
	if deps.IDGenSVC, err = idgenimpl.New(deps.CacheCli); err != nil {
		return nil, fmt.Errorf("init id gen svc failed, err=%w", err)
	}
	if deps.TOSClient, err = storage.New(ctx); err != nil {
		return nil, fmt.Errorf("init storage failed, err=%w", err)
	}

	return deps, err
}
