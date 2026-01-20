namespace go llm.domain.manage

include "common.thrift"

struct AdminModel {
    1: optional i64 model_id (api.js_conv='true', go.tag='json:"model_id"')
    2: optional string name
    3: optional string desc
    4: optional string icon_uri
    5: optional string icon_url
    6: optional Ability ability
    7: optional Protocol protocol
    8: optional ProtocolConfig protocol_config
    9: optional map<common.Scenario, ScenarioConfig> scenario_configs
    10: optional ParamConfig param_config
}

struct Ability {
    1: optional i64 max_context_tokens (api.js_conv='true', go.tag='json:"max_context_tokens"')
    2: optional i64 max_input_tokens (api.js_conv='true', go.tag='json:"max_input_tokens"')
    3: optional i64 max_output_tokens (api.js_conv='true', go.tag='json:"max_output_tokens"')
    4: optional bool function_call
    5: optional bool json_mode
    6: optional bool multi_modal
    7: optional AbilityMultiModal ability_multi_modal
}

struct AbilityMultiModal {
    1: optional bool image
    2: optional bool function_call
    3: optional bool video
    4: optional bool audio
    5: optional bool prefill_resp
}

struct ProtocolConfig {
    1: optional string base_url
    2: optional string api_key
    3: optional string model
    4: optional ProtocolConfigOpenAI protocol_config_openai
    5: optional ProtocolConfigClaude protocol_config_claude
    6: optional ProtocolConfigDeepSeek protocol_config_deepseek
    7: optional ProtocolConfigOllama protocol_config_ollama
    8: optional ProtocolConfigQwen protocol_config_qwen
    9: optional ProtocolConfigGemini protocol_config_gemini
}

struct ProtocolConfigOpenAI {
    1: optional bool by_azure
    2: optional string api_version
    3: optional string response_format_type
    4: optional string response_format_json_schema
}
struct ProtocolConfigClaude {
    1: optional bool by_bedrock
    // bedrock config
    2: optional string access_key
    3: optional string secret_access_key
    4: optional string session_token
    5: optional string region
    6: optional i32    budget_tokens

}
struct ProtocolConfigDeepSeek {
    1: optional string response_format_type
}

struct ProtocolConfigGemini {
    1: optional i32 backend
    2: optional string project
    3: optional string location
    4: optional string api_version
    5: optional string timeout_ms
    6: optional bool include_thoughts
    7: optional i32 thinking_budget
    8: optional map<string, list<string>> headers;
}

struct ProtocolConfigOllama {
    1: optional string format
    2: optional i64 keep_alive_ms (api.js_conv='true', go.tag='json:"keep_alive_ms"')
}

struct ProtocolConfigQwen {
    1: optional string response_format_type
    2: optional string response_format_json_schema
}

struct ScenarioConfig {
    1: optional common.Scenario scenario
    3: optional Quota quota
    4: optional bool unavailable
}

struct ParamConfig {
    1: optional list<ParamSchema> param_schemas
}

struct ParamSchema {
    1: optional string name // 实际名称
    2: optional string label // 展示名称
    3: optional string desc
    4: optional ParamType type
    5: optional string min
    6: optional string max
    7: optional map<DefaultType, string> default_value
    8: optional list<ParamOption> options
}

struct ParamOption {
    1: optional string value // 实际值
    2: optional string label // 展示值
}

struct Quota {
    1: optional i64 qpm (api.js_conv='true', go.tag='json:"qpm"')
    2: optional i64 tpm (api.js_conv='true', go.tag='json:"tpm"')
}

typedef string Protocol (ts.enum="true")

const Protocol protocol_openai = "openai"
const Protocol protocol_claude = "claude"
const Protocol protocol_deepseek = "deepseek"
const Protocol protocol_ollama = "ollama"
const Protocol protocol_gemini = "gemini"
const Protocol protocol_qwen = "qwen"

typedef string ParamType (ts.enum="true")
const ParamType param_type_float = "float"
const ParamType param_type_int = "int"
const ParamType param_type_boolean = "boolean"
const ParamType param_type_string = "string"

typedef string DefaultType (ts.enum="true")

const DefaultType default_type_value = "default"
const DefaultType default_type_creative = "creative"
const DefaultType default_type_balance = "balance"
const DefaultType default_type_precise = "precise"