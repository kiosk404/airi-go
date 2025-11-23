package ctxutil

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/foundation/openauth/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/ctxcache"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

func GetApiAuthFromCtx(ctx context.Context) *entity.ApiKey {
	data, ok := ctxcache.Get[*entity.ApiKey](ctx, consts.OpenapiAuthKeyInCtx)

	if !ok {
		return nil
	}
	return data
}

func MustGetUIDFromApiAuthCtx(ctx context.Context) int64 {
	apiKeyInfo := GetApiAuthFromCtx(ctx)
	if apiKeyInfo == nil {
		panic("mustGetUIDFromApiAuthCtx: apiKeyInfo is nil")
	}
	return apiKeyInfo.UserID
}
