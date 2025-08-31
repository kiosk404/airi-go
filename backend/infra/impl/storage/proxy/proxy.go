package proxy

import (
	"context"
	"net"
	"net/url"
	"os"

	"github.com/kiosk404/airi-go/backend/pkg/ctxcache"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

func CheckIfNeedReplaceHost(ctx context.Context, originURLStr string) (ok bool, proxyURL string) {
	// url parse
	originURL, err := url.Parse(originURLStr)
	if err != nil {
		logs.Warn("[CheckIfNeedReplaceHost] url parse failed, err: %v", err)
		return false, ""
	}

	proxyPort := os.Getenv(consts.MinIOProxyEndpoint) // :8889
	if proxyPort == "" {
		return false, ""
	}

	currentHost, ok := ctxcache.Get[string](ctx, consts.HostKeyInCtx)
	if !ok {
		return false, ""
	}

	currentScheme, ok := ctxcache.Get[string](ctx, consts.RequestSchemeKeyInCtx)
	if !ok {
		return false, ""
	}

	host, _, err := net.SplitHostPort(currentHost)
	if err != nil {
		host = currentHost
	}

	minioProxyHost := host + proxyPort
	originURL.Host = minioProxyHost
	originURL.Scheme = currentScheme
	logs.Debug("[CheckIfNeedReplaceHost] reset originURL.String = %s", originURL.String())
	return true, originURL.String()
}
