package consts

import (
	"time"
)

const (
	SessionDataKeyInCtx = "session_data_key_in_ctx"
	OpenapiAuthKeyInCtx = "openapi_auth_key_in_ctx"

	HostKeyInCtx          = "HOST_KEY_IN_CTX"
	RequestSchemeKeyInCtx = "REQUEST_SCHEME_IN_CTX"

	APIConnectorID     = int64(9527)
	LocalStoragePath   = "LOCAL_STORAGE_PATH"
	StorageType        = "STORAGE_TYPE"
	MinIOAK            = "MINIO_AK"
	MinIOSK            = "MINIO_SK"
	MinIOEndpoint      = "MINIO_ENDPOINT"
	MinIOProxyEndpoint = "MINIO_PROXY_ENDPOINT"
	MinIOAPIHost       = "MINIO_API_HOST"
	StorageBucket      = "STORAGE_BUCKET"
)

const (
	SessionMaxAgeSecond    = 30 * 24 * 60 * 60
	DefaultSessionDuration = SessionMaxAgeSecond * time.Second
)

const (
	DisableUserRegistration  = "DISABLE_USER_REGISTRATION"
	AllowRegistrationAccount = "ALLOW_REGISTRATION_ACCOUNT"
)
