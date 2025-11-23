package service

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/repo"
	"github.com/kiosk404/airi-go/backend/pkg/utils/safego"
)

type Components struct {
	IDGen      idgen.IDGenerator
	DB         rdb.Provider
	OSS        storage.Storage
	PluginRepo repo.PluginRepository
	ToolRepo   repo.ToolRepository
	OAuthRepo  repo.OAuthRepository
}

func NewService(components *Components) PluginService {
	impl := &pluginServiceImpl{
		db:         components.DB,
		oss:        components.OSS,
		pluginRepo: components.PluginRepo,
		toolRepo:   components.ToolRepo,
		oauthRepo:  components.OAuthRepo,
	}

	initOnce.Do(func() {
		ctx := context.Background()
		safego.Go(ctx, func() {
			impl.processOAuthAccessToken(ctx)
		})
	})

	return impl
}

type pluginServiceImpl struct {
	db         rdb.Provider
	oss        storage.Storage
	pluginRepo repo.PluginRepository
	toolRepo   repo.ToolRepository
	oauthRepo  repo.OAuthRepository
}
