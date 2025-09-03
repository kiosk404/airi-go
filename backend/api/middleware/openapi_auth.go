package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	httpwarp "github.com/kiosk404/airi-go/backend/api/http"
	"github.com/kiosk404/airi-go/backend/application/openauth"
	"github.com/kiosk404/airi-go/backend/pkg/ctxcache"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/types/consts"
	"github.com/kiosk404/airi-go/backend/types/errno"
)

const HeaderAuthorizationKey = "Authorization"

var needAuthPath = map[string]bool{
	"/v3/chat":                         true,
	"/v1/conversations":                true,
	"/v1/conversation/create":          true,
	"/v1/conversation/message/list":    true,
	"/v1/files/upload":                 true,
	"/v1/workflow/run":                 true,
	"/v1/workflow/stream_run":          true,
	"/v1/workflow/stream_resume":       true,
	"/v1/workflow/get_run_history":     true,
	"/v1/bot/get_online_info":          true,
	"/v1/workflows/chat":               true,
	"/v1/workflow/conversation/create": true,
	"/v3/chat/cancel":                  true,
}

var needAuthFunc = map[string]bool{
	"^/v1/conversations/[0-9]+/clear$": true, // v1/conversations/:conversation_id/clear
	"^/v1/bots/[0-9]+$":                true,
	"^/v1/conversations/[0-9]+$":       true,

	"^/v1/workflows/[0-9]+$": true,
	"^/v1/apps/[0-9]+$":      true,
}

func parseBearerAuthToken(authHeader string) string {
	if len(authHeader) == 0 {
		return ""
	}
	parts := strings.Split(authHeader, "Bearer")
	if len(parts) != 2 {
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if len(token) == 0 {
		return ""
	}

	return token
}

func isNeedOpenapiAuth(c *gin.Context) bool {
	isNeedAuth := false

	uriPath := c.Request.URL.Path

	for rule, res := range needAuthFunc {
		if regexp.MustCompile(rule).MatchString(uriPath) {
			isNeedAuth = res
			break
		}
	}

	if needAuthPath[uriPath] {
		isNeedAuth = true
	}

	return isNeedAuth
}

func OpenapiAuthMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestAuthType := c.GetInt64(RequestAuthTypeStr)
		if requestAuthType != int64(RequestAuthTypeOpenAPI) {
			c.Next()
			return
		}

		ctx := c.Request.Context()

		// open api auth
		if len(c.Request.Header.Get(HeaderAuthorizationKey)) == 0 {
			httpwarp.InternalError(c,
				errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "missing authorization in header")))
			return
		}

		apiKey := parseBearerAuthToken(c.Request.Header.Get(HeaderAuthorizationKey))
		if len(apiKey) == 0 {
			httpwarp.InternalError(c,
				errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "missing api_key in request")))
			return
		}

		md5Hash := md5.Sum([]byte(apiKey))
		md5Key := hex.EncodeToString(md5Hash[:])
		apiKeyInfo, err := openauth.OpenAuthApplication.CheckPermission(ctx, md5Key)

		if err != nil {
			logs.Error("OpenAuthApplication.CheckPermission failed, err=%v", err)
			httpwarp.InternalError(c,
				errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", err.Error())))
			return
		}

		if apiKeyInfo == nil {
			httpwarp.InternalError(c,
				errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "api key invalid")))
			return
		}

		apiKeyInfo.ConnectorID = consts.APIConnectorID
		logs.Info("OpenapiAuthMW: apiKeyInfo=%v", conv.DebugJsonToStr(apiKeyInfo))
		ctxcache.Store(ctx, consts.OpenapiAuthKeyInCtx, apiKeyInfo)
		err = openauth.OpenAuthApplication.UpdateLastUsedAt(ctx, apiKeyInfo.ID, apiKeyInfo.UserID)
		if err != nil {
			logs.Error("OpenAuthApplication.UpdateLastUsedAt failed, err=%v", err)
		}
		c.Next()
	}
}
