package repo

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/component/plugin/infra/dao"
)

type OAuthRepository interface {
	GetAuthorizationCode(ctx context.Context, meta *dao.AuthorizationCodeMeta) (info *dao.AuthorizationCodeInfo, exist bool, err error)
	UpsertAuthorizationCode(ctx context.Context, info *dao.AuthorizationCodeInfo) (err error)
	UpdateAuthorizationCodeLastActiveAt(ctx context.Context, meta *dao.AuthorizationCodeMeta, lastActiveAtMs int64) (err error)
	BatchDeleteAuthorizationCodeByIDs(ctx context.Context, ids []int64) (err error)
	DeleteAuthorizationCode(ctx context.Context, meta *dao.AuthorizationCodeMeta) (err error)
	GetAuthorizationCodeRefreshTokens(ctx context.Context, nextRefreshAt int64, limit int) (infos []*dao.AuthorizationCodeInfo, err error)
	DeleteExpiredAuthorizationCodeTokens(ctx context.Context, expireAt int64, limit int) (err error)
	DeleteInactiveAuthorizationCodeTokens(ctx context.Context, lastActiveAt int64, limit int) (err error)
}
