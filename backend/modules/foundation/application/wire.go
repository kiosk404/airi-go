//go:build wireinject
// +build wireinject

package application

import (
	"context"
	"github.com/google/wire"
	"github.com/kiosk404/airi-go/backend/api/model/foundation/user"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
)

var (
	userApp = wire.NewSet(
		InitUserService,
		wire.Bind(new(user.UserService), new(*UserApplicationService)),
	)
)

func InitUserApplication(ctx context.Context,
	provider rdb.Provider,
	oss storage.Storage,
	idGen idgen.IDGenerator) (user.UserService, error) {
	wire.Build(userApp)
	return nil, nil
}
