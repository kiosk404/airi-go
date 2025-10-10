include "../../base.thrift"
include "../../app/bot_common.thrift"

namespace go playground

// Onboarding json structure
struct OnboardingContent {
    1: optional string       prologue            // Introductory remarks (C-end usage scenarios, only 1; background scenarios, possibly multiple)
    2: optional list<string> suggested_questions // suggestion question
    3: optional bot_common.SuggestedQuestionsShowMode suggested_questions_show_mode
}

struct GetBotPopupInfoRequest {
    1: required list<BotPopupType> bot_popup_types
    2: required i64                bot_id          (agw.js_conv="str" api.js_conv="true")

    255: base.Base Base (api.none="true")
}

struct GetBotPopupInfoResponse {
    1:   required BotPopupInfoData data

    253: required i64              code
    254: required string           msg
    255: required base.BaseResp BaseResp (api.none="true")
}

struct BotPopupInfoData {
    1: required map<BotPopupType,i64> bot_popup_count_info
}

enum BotPopupType {
    AutoGenBeforePublish = 1
}


struct UpdateBotPopupInfoResponse {
    253: required i64    code
    254: required string msg
    255: required base.BaseResp BaseResp (api.none="true")
}


struct UpdateBotPopupInfoRequest {
    1: required BotPopupType bot_popup_type
    2: required i64          bot_id         (agw.js_conv="str" api.js_conv="true")

    255: base.Base Base (api.none="true")
}

enum ResourceType {
    DraftBot = 1
    Project = 2
}

enum BehaviorType {
    Visit = 1
    Edit = 2
}

struct ReportUserBehaviorRequest {
    1: required i64 ResourceID (api.body = "resource_id",api.js_conv="true")
    2: required ResourceType ResourceType (api.body="resource_type")
    3: required BehaviorType BehaviorType (api.body="behavior_type")
    4: optional i64 SpaceID (agw.js_conv="str",api.js_conv="true",api.body="space_id",agw.key="space_id") // This requirement must be passed on

    255: base.Base Base (api.none="true")
}

struct ReportUserBehaviorResponse {
    253: required i64    code
    254: required string msg
    255: required base.BaseResp BaseResp
}

service PlaygroundService {
    // public popup_info
    GetBotPopupInfoResponse GetBotPopupInfo (1:GetBotPopupInfoRequest request)(api.post='/api/playground_api/operate/get_bot_popup_info', api.category="account",agw.preserve_base="true")
    UpdateBotPopupInfoResponse UpdateBotPopupInfo (1:UpdateBotPopupInfoRequest request)(api.post='/api/playground_api/operate/update_bot_popup_info', api.category="account",agw.preserve_base="true")
    ReportUserBehaviorResponse ReportUserBehavior(1:ReportUserBehaviorRequest request)(api.post='/api/playground_api/report_user_behavior', api.category="playground_api",agw.preserve_base="true")
}
