package local

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/kiosk404/airi-go/backend/infra/contract/imagex"
	"github.com/kiosk404/airi-go/backend/pkg/ctxcache"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

func NewStorageImageX(ctx context.Context, pathDir string) (imagex.ImageX, error) {
	l, err := getLocalClient(pathDir)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (l *LocalClient) GetUploadHost(ctx context.Context) string {
	currentHost, ok := ctxcache.Get[string](ctx, consts.HostKeyInCtx)
	if !ok {
		return ""
	}
	return currentHost + consts.ApplyUploadActionURI
}

func (l *LocalClient) GetServerID() string {
	return ""
}

func (l *LocalClient) GetUploadAuth(ctx context.Context, opt ...imagex.UploadAuthOpt) (*imagex.SecurityToken, error) {
	scheme := strings.ToLower(os.Getenv(consts.StorageUploadHTTPScheme))
	if scheme == "" {
		scheme = "http"
	}
	return &imagex.SecurityToken{
		AccessKeyID:     "",
		SecretAccessKey: "",
		SessionToken:    "",
		ExpiredTime:     time.Now().Add(time.Hour).Format("2006-01-02 15:04:05"),
		CurrentTime:     time.Now().Format("2006-01-02 15:04:05"),
		HostScheme:      scheme,
	}, nil
}

func (l *LocalClient) GetResourceURL(ctx context.Context, uri string, opts ...imagex.GetResourceOpt) (*imagex.ResourceURL, error) {
	url, err := l.GetObjectUrl(ctx, uri)
	if err != nil {
		return nil, err
	}
	return &imagex.ResourceURL{
		URL: url,
	}, nil
}

func (l *LocalClient) Upload(ctx context.Context, data []byte, opts ...imagex.UploadAuthOpt) (*imagex.UploadResult, error) {
	return nil, nil
}

func (l *LocalClient) GetUploadAuthWithExpire(ctx context.Context, expire time.Duration, opt ...imagex.UploadAuthOpt) (*imagex.SecurityToken, error) {
	return nil, nil
}
