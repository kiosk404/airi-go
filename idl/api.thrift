namespace go airi

include "./app/bot_open_api.thrift"
include "./app/develer_api.thrift"
include "./foundation/openapiauth.thrift"
include "./foundation/user.thrift"
include "./llm/manage.thrift"
include "./llm/runtime.thrift"
include "./component/plugin/plugin_develop.thrift"

service BotOpenApiService extends bot_open_api.BotOpenApiService {}
service DeveloperApiService extends develer_api.DeveloperApiService{}
service OpenAPIAuthService extends openapiauth.OpenAPIAuthService {}
service UserService extends user.UserService {}
service LLMManageService extends manage.LLMManageService {}
service LLMRuntimeService extends runtime.LLMRuntimeService {}
service PluginDevelopService extends plugin_develop.PluginDevelopService {}