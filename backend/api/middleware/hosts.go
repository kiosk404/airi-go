package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/pkg/ctxcache"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

func SetHostMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxcache.Store(c, consts.HostKeyInCtx, string(c.Request.Host))
		ctxcache.Store(c, consts.RequestSchemeKeyInCtx, string(c.Request.URL.Scheme))
		c.Next()
	}
}
