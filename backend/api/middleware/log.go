package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kiosk404/airi-go/backend/infra/http/ginutil"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

func AccessLogMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		path := c.Request.URL.RawPath
		latency := time.Since(start)
		method := c.Request.Method
		clientIP := c.ClientIP()

		handlerPkgPath := strings.Split(c.HandlerName(), "/")
		handleName := ""
		if len(handlerPkgPath) > 0 {
			handleName = handlerPkgPath[len(handlerPkgPath)-1]
		}

		requestType := ginutil.GetInt32(c, RequestAuthTypeStr)
		baseLog := fmt.Sprintf("| %s | %s | %d | %v | %s | %s | %v | %s | %d ",
			c.Request.URL.Scheme, c.Request.Host, status,
			latency, clientIP, method, path, handleName, requestType)

		switch {
		case status >= http.StatusInternalServerError:
			logs.Error("%s", baseLog)
		case status >= http.StatusBadRequest:
			logs.Warn("%s", baseLog)
		default:
			logs.Info("%s ", baseLog)
		}
	}
}

func SetLogIDMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		logID := uuid.New().String()
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "log-id", logID)

		c.Header("X-Log-ID", logID)
		c.Next()
	}
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b)) // nolint
}
