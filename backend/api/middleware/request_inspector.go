package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const RequestAuthTypeStr = "RequestAuthTypeStr"

type RequestAuthType = int32

const (
	RequestAuthTypeWebAPI     RequestAuthType = 0
	RequestAuthTypeOpenAPI    RequestAuthType = 1
	RequestAuthTypeStaticFile RequestAuthType = 2
)

func RequestInspectorMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		authType := RequestAuthTypeWebAPI // default is web api, session auth
		if isNeedOpenapiAuth(c) {
			authType = RequestAuthTypeOpenAPI
		} else if isStaticFile(c) {
			authType = RequestAuthTypeStaticFile
		}

		c.Set(RequestAuthTypeStr, authType)
		c.Next()
	}
}

var staticFilePath = map[string]bool{
	"/static":      true,
	"/":            true,
	"/sign":        true,
	"/favicon.png": true,
}

func isStaticFile(c *gin.Context) bool {
	path := c.Request.URL.Path
	if staticFilePath[path] {
		return true
	}

	if strings.HasPrefix(path, "/static/") ||
		strings.HasPrefix(path, "/explore/") ||
		strings.HasPrefix(path, "/space/") {
		return true
	}

	if path == "/information/auth/success" {
		return true
	}

	return false
}
