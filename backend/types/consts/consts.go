package consts

import (
	"time"
)

const (
	RunMode = "RUN_MODE"

	SessionDataKeyInCtx = "session_data_key_in_ctx"
	OpenapiAuthKeyInCtx = "openapi_auth_key_in_ctx"

	HostKeyInCtx          = "HOST_KEY_IN_CTX"
	RequestSchemeKeyInCtx = "REQUEST_SCHEME_IN_CTX"

	LocalStoragePath   = "LOCAL_STORAGE_PATH"
	StorageType        = "STORAGE_TYPE"
	MinIOAK            = "MINIO_AK"
	MinIOSK            = "MINIO_SK"
	MinIOEndpoint      = "MINIO_ENDPOINT"
	MinIOProxyEndpoint = "MINIO_PROXY_ENDPOINT"
	MinIOAPIHost       = "MINIO_API_HOST"
	StorageBucket      = "STORAGE_BUCKET"

	FileUploadComponentType       = "FILE_UPLOAD_COMPONENT_TYPE"
	FileUploadComponentTypeImageX = "imagex"

	StorageUploadHTTPScheme = "STORAGE_UPLOAD_HTTP_SCHEME"
)

const (
	MySQLDomain   = "AIRI_GO_MYSQL_DOMAIN"
	MySQLPort     = "AIRI_GO_MYSQL_PORT"
	MySQLUser     = "AIRI_GO_MYSQL_USER"
	MySQLPassport = "AIRI_GO_MYSQL_PASSWORD"
	MySQLDatabase = "AIRI_GO_MYSQL_DATABASE"
)

const (
	MQTypeKey                = "AIRI_MQ_TYPE"
	MQServer                 = "MQ_NAME_SERVER"
	RMQSecretKey             = "RMQ_SECRET_KEY"
	RMQAccessKey             = "RMQ_ACCESS_KEY"
	RMQTopicApp              = "airi_search_app"
	RMQTopicResource         = "airi_search_resource"
	RMQTopicKnowledge        = "airi_knowledge"
	RMQConsumeGroupResource  = "cg_search_resource"
	RMQConsumeGroupApp       = "cg_search_app"
	RMQConsumeGroupKnowledge = "cg_knowledge"
)

const (
	SessionMaxAgeSecond    = 30 * 24 * 60 * 60
	DefaultSessionDuration = SessionMaxAgeSecond * time.Second
)

const (
	ApplyUploadActionURI = "/api/common/upload/apply_upload_action"
	UploadURI            = "/api/common/upload"
)

const (
	DisableUserRegistration  = "DISABLE_USER_REGISTRATION"
	AllowRegistrationAccount = "ALLOW_REGISTRATION_ACCOUNT"
)

const (
	SearchESVersion = "SEARCH_ES_VERSION"
	BleveIndexPath  = "BLEVE_INDEX_PATH"
)
