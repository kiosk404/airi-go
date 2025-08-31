namespace go airi

include "./app/bot_open_api.thrift"
include "./permission/openapiauth_service.thrift"
include "./passport/passport.thrift"

service BotOpenApiService extends bot_open_api.BotOpenApiService {}
service OpenAPIAuthService extends openapiauth_service.OpenAPIAuthService {}
service PassportService extends passport.PassportService {}