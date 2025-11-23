package http

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func GetOriginHost(c *gin.Context) string {
	// 尝试从Origin头获取
	origin := c.Request.Header.Get("Origin")
	if origin != "" {
		u, err := url.Parse(origin)
		if err == nil {
			return u.Hostname()
		}
	}

	// 尝试从Host头获取
	host := c.Request.Header.Get("Host")
	if host != "" {
		return host
	}

	// 回退到请求URL的主机
	return c.Request.URL.Host
}
