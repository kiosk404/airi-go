package consts

import (
	"time"
)

const (
	SessionDataKeyInCtx = "session_data_key_in_ctx"
	OpenapiAuthKeyInCtx = "openapi_auth_key_in_ctx"

	HostKeyInCtx          = "HOST_KEY_IN_CTX"
	RequestSchemeKeyInCtx = "REQUEST_SCHEME_IN_CTX"

	APIConnectorID = int64(9527)

	MinIOProxyEndpoint = "MINIO_PROXY_ENDPOINT"
)

const (
	SessionMaxAgeSecond    = 30 * 24 * 60 * 60
	DefaultSessionDuration = SessionMaxAgeSecond * time.Second
)

const (
	DisableUserRegistration  = "DISABLE_USER_REGISTRATION"
	AllowRegistrationAccount = "ALLOW_REGISTRATION_ACCOUNT"
)
