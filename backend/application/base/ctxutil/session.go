package ctxutil

import (
	"context"

	"github.com/kiosk404/airi-go/backend/domain/user/entity"
	"github.com/kiosk404/airi-go/backend/pkg/ctxcache"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

func GetUserSessionFromCtx(ctx context.Context) *entity.Session {
	data, ok := ctxcache.Get[*entity.Session](ctx, consts.SessionDataKeyInCtx)
	if !ok {
		return nil
	}

	return data
}

func MustGetUIDFromCtx(ctx context.Context) int64 {
	sessionData := GetUserSessionFromCtx(ctx)
	if sessionData == nil {
		panic("mustGetUIDFromCtx: sessionData is nil")
	}

	return sessionData.UserID
}

func GetUIDFromCtx(ctx context.Context) *int64 {
	sessionData := GetUserSessionFromCtx(ctx)
	if sessionData == nil {
		return nil
	}

	return &sessionData.UserID
}
