package dao

import (
	"github.com/kiosk404/airi-go/backend/api/model/component/plugin_develop/common"
	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/consts"
	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/model"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

type GetOAuthStatusResponse struct {
	IsOauth  bool
	Status   common.OAuthStatus
	OAuthURL string
}

type AgentPluginOAuthStatus struct {
	PluginID      int64
	PluginName    string
	PluginIconURL string
	Status        common.OAuthStatus
}

type GetAccessTokenRequest struct {
	UserID    string
	PluginID  *int64
	Mode      consts.AuthzSubType
	OAuthInfo *OAuthInfo
}

type PluginAuthInfo struct {
	AuthzType    *consts.AuthzType
	Location     *consts.HTTPParamLocation
	Key          *string
	ServiceToken *string
	OAuthInfo    *string
	AuthzSubType *consts.AuthzSubType
	AuthzPayload *string
}

type OAuthInfo struct {
	OAuthMode         consts.AuthzSubType
	AuthorizationCode *AuthorizationCodeInfo
}

type OAuthState struct {
	ClientName OAuthProvider `json:"client_name"`
	UserID     string        `json:"user_id"`
	PluginID   int64         `json:"plugin_id"`
	IsDraft    bool          `json:"is_draft"`
}

type AuthorizationCodeMeta struct {
	UserID   string
	PluginID int64
	IsDraft  bool
}

type AuthorizationCodeInfo struct {
	RecordID             int64
	Meta                 *AuthorizationCodeMeta
	Config               *model.OAuthAuthorizationCodeConfig
	AccessToken          string
	RefreshToken         string
	TokenExpiredAtMS     int64
	NextTokenRefreshAtMS *int64
	LastActiveAtMS       int64
}

func (a *AuthorizationCodeInfo) GetNextTokenRefreshAtMS() int64 {
	if a == nil {
		return 0
	}
	return ptr.FromOrDefault(a.NextTokenRefreshAtMS, 0)
}
