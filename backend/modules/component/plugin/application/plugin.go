package application

import (
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/repo"
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/domain/service"
	search "github.com/kiosk404/airi-go/backend/modules/data/search/domain/service"
	user "github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/service"
)

var PluginApplicationSVC = &PluginApplicationService{}

type PluginApplicationService struct {
	DomainSVC service.PluginService
	eventbus  search.ResourceEventBus
	oss       storage.Storage
	userSVC   user.User

	toolRepo   repo.ToolRepository
	pluginRepo repo.PluginRepository
}
