package storage

import (
	"context"
	"os"

	"github.com/kiosk404/airi-go/backend/infra/contract/imagex"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	"github.com/kiosk404/airi-go/backend/infra/impl/storage/local"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

type Storage = storage.Storage

func New(ctx context.Context) (Storage, error) {
	//return minio.New(
	//	ctx,
	//	os.Getenv(consts.MinIOEndpoint),
	//	os.Getenv(consts.MinIOAK),
	//	os.Getenv(consts.MinIOSK),
	//	os.Getenv(consts.StorageBucket),
	//	false,
	//)
	return local.New(ctx, os.Getenv(consts.LocalStoragePath))
}

func NewImageX(ctx context.Context) (imagex.ImageX, error) {
	return local.NewStorageImageX(ctx, os.Getenv(consts.LocalStoragePath))
}
