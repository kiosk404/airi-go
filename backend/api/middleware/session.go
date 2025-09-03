package middleware

import (
	"github.com/gin-gonic/gin"
	httpwarp "github.com/kiosk404/airi-go/backend/api/http"
	"github.com/kiosk404/airi-go/backend/application/user"
	"github.com/kiosk404/airi-go/backend/domain/user/entity"
	"github.com/kiosk404/airi-go/backend/pkg/ctxcache"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/http/ginutil"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/types/consts"
	"github.com/kiosk404/airi-go/backend/types/errno"
)

var noNeedSessionCheckPath = map[string]bool{
	"/health":                     true,
	"/api/passport/web/login/":    true,
	"/api/passport/web/register/": true,
}

func SessionAuthMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestAuthType := ginutil.GetInt32(c, RequestAuthTypeStr)
		if requestAuthType != int32(RequestAuthTypeWebAPI) {
			c.Next()
			return
		}

		if noNeedSessionCheckPath[c.Request.URL.Path] {
			c.Next()
			return
		}

		s, err := c.Cookie(entity.SessionKey)
		if len(s) == 0 || err != nil {
			logs.Error("[SessionAuthMW] session id is nil")
			httpwarp.InternalError(c,
				errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "missing session_key in cookie")))
			return
		}

		// sessionID -> sessionData
		session, err := user.UserApplicationSVC.ValidateSession(c, string(s))
		if err != nil {
			logs.Error("[SessionAuthMW] validate session failed, err: %v", err)
			httpwarp.InternalError(c, err)
			return
		}

		if session != nil {
			ctxcache.Store(c, consts.SessionDataKeyInCtx, session)
		}

		c.Next()
	}
}
