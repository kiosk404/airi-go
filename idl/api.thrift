namespace go airi

include "./app/bot_open_api.thrift"
include "./permission/openapiauth_service.thrift"
include "./passport/passport.thrift"
include "./llm/manage.thrift"
include "./llm/runtime.thrift"

service BotOpenApiService extends bot_open_api.BotOpenApiService {}
service OpenAPIAuthService extends openapiauth_service.OpenAPIAuthService {}
service PassportService extends passport.PassportService {}
service LLMManageService extends manage.LLMManageService {}
service LLMRuntimeService extends runtime.LLMRuntimeService {}