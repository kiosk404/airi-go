package service

import (
	"context"
	"sync"
	"time"

	"github.com/kiosk404/airi-go/backend/modules/component/plugin/infra/dao"
)

var (
	initOnce           = sync.Once{}
	lastActiveInterval = 15 * 24 * time.Hour
	failedCache        = sync.Map{}
)

func (p *pluginServiceImpl) processOAuthAccessToken(ctx context.Context) {
}

func (p *pluginServiceImpl) GetAccessToken(ctx context.Context, oa *dao.OAuthInfo) (accessToken string, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) OAuthCode(ctx context.Context, code string, state *dao.OAuthState) (err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) RevokeAccessToken(ctx context.Context, meta *dao.AuthorizationCodeMeta) (err error) {
	return p.oauthRepo.DeleteAuthorizationCode(ctx, meta)
}

func (p *pluginServiceImpl) GetOAuthStatus(ctx context.Context, userID, pluginID int64) (resp *dao.GetOAuthStatusResponse, err error) {
	panic("implement me")
}

func (p *pluginServiceImpl) GetAgentPluginsOAuthStatus(ctx context.Context, userID, agentID int64) (status []*dao.AgentPluginOAuthStatus, err error) {
	panic("implement me")
}
