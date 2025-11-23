package application

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	search "github.com/kiosk404/airi-go/backend/modules/data/search/domain/service"
	user "github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/service"
)

type ServiceComponents struct {
	IDGen    idgen.IDGenerator
	DB       rdb.Provider
	OSS      storage.Storage
	EventBus search.ResourceEventBus
	UserSVC  user.User
}

func InitService(ctx context.Context, components *ServiceComponents) (*PluginApplicationService, error) {

	return PluginApplicationSVC, nil
}
