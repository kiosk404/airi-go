package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/pkg/ctxcache"
)

func ContextCacheMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		cacheCtx := ctxcache.Init(ctx)
		c.Request = c.Request.WithContext(cacheCtx)
		c.Next()
	}
}
