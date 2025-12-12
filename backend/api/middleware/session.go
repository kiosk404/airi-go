package middleware

import (
	"regexp"

	"github.com/gin-gonic/gin"
	httpwarp "github.com/kiosk404/airi-go/backend/api/http"
	user "github.com/kiosk404/airi-go/backend/modules/foundation/user/application"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/ctxcache"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/http/ginutil"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

var noNeedSessionCheckPath = map[string]bool{
	"/favicon.ico":                      true,
	"/health":                           true,
	"/api/foundation/v1/users/login":    true,
	"/api/foundation/v1/users/register": true,
}

var noNeedSessionCheckPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^/static/files/.*$`),
}

func SessionAuthMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestAuthType := ginutil.GetInt32(c, RequestAuthTypeStr)
		if requestAuthType != int32(RequestAuthTypeWebAPI) {
			c.Next()
			return
		}

		if noNeedSessionCheck(c.Request.URL.Path) {
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
			ctxcache.Store(c.Request.Context(), consts.SessionDataKeyInCtx, session)
		}

		c.Next()
	}
}

func noNeedSessionCheck(path string) bool {
	// 检查精确匹配
	if noNeedSessionCheckPath[path] {
		return true
	}

	// 检查正则模式匹配
	for _, pattern := range noNeedSessionCheckPatterns {
		if pattern.MatchString(path) {
			return true
		}
	}

	return false
}
