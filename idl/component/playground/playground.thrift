include "../../base.thrift"
include "../../app/bot_common.thrift"
include "./prompt_resource.thrift"

namespace go playground

struct UpdateDraftBotInfoAgwResponse {
    1: required UpdateDraftBotInfoAgwData data,

    253: required i64                   code,
    254: required string                msg,
    255: required base.BaseResp BaseResp (api.none="true")
}

struct UpdateDraftBotInfoAgwData {
    1: optional bool   has_change       // Is there any change?
    2:          bool   check_not_pass   // True: The machine audit verification failed
    3: optional Branch branch           // Which branch is it currently on?
    4: optional bool   same_with_online
    5: optional string check_not_pass_msg // The machine audit verification failed the copy.
}

// branch
enum Branch {
    Undefined     = 0
    PersonalDraft = 1 // draft
    Base          = 2 // Space draft
    Publish       = 3 // Online version, used in diff scenarios
}

struct UpdateDraftBotInfoAgwRequest {
    1: optional bot_common.BotInfoForUpdate bot_info
    2: optional i64   base_commit_version (api.js_conv='true',agw.js_conv="str")

    255: base.Base Base (api.none="true")
}

struct GetDraftBotInfoAgwRequest {
    1: required i64  bot_id  (api.js_conv='true',agw.js_conv="str") // Draft bot_id
    2: optional string  version  // Check the history, the id of the historical version, corresponding to the id of the bot_draft_history
    3: optional string  commit_version // Query specifies commit_version version, pre-release use, seems to be the same thing as version, but the acquisition logic is different

    255: base.Base Base (api.none="true")
}

struct GetDraftBotInfoAgwResponse {
    1: required GetDraftBotInfoAgwData data,

    253: required i64                   code,
    254: required string                msg,
    255: required base.BaseResp BaseResp (api.none="true")
}


struct GetDraftBotInfoAgwData {
    1: required bot_common.BotInfo bot_info // core bot data
    2: optional BotOptionData bot_option_data // bot option information
    3: optional bool            has_unpublished_change // Are there any unpublished changes?
    4: optional BotMarketStatus bot_market_status      // The product status after the bot is put on the shelves
    5: optional bool            in_collaboration       // Is the bot in multiplayer cooperative mode?
    6: optional bool            same_with_online       // Is the content committed consistent with the online content?
    7: optional bool            editable               // For frontend, permission related, can the current user edit this bot
    8: optional bool            deletable              // For frontend, permission related, can the current user delete this bot
    9: optional UserInfo        publisher              // Is the publisher of the latest release version
    10:         bool has_publish // Has it been published?
    11:         i64 space_id    (api.js_conv='true',agw.js_conv="str")  // Space ID
    12: optional Branch              branch          // What branch did you get the content of?
    13: optional string              commit_version  // If branch=PersonalDraft, the version number of checkout/rebase; if branch = base, the committed version
    14: optional string              committer_name  // For the front end, the most recent author
    15: optional string              commit_time     // For frontend, commit time
    16: optional string              publish_time    // For frontend, release time
}

struct BotOptionData {
//    1: optional map<i64,ModelDetail>        model_detail_map      // model details
//    2: optional map<i64,PluginDetal>        plugin_detail_map     // plugin details
//    3: optional map<i64,PluginAPIDetal>     plugin_api_detail_map // Plugin API Details
//    4: optional map<i64,WorkflowDetail>     workflow_detail_map   // Workflow Details
//    5: optional map<string,KnowledgeDetail> knowledge_detail_map  // Knowledge Details
//    6: optional list<shortcut_command.ShortcutCommand>   shortcut_command_list  // Quick command list
}

enum BotMarketStatus {
    Offline = 0 // offline
    Online  = 1 // put on the shelves
}

struct UserInfo {
    1: i64    user_id   (api.js_conv='true',agw.js_conv="str")  // user id
    2: string name     // user name
    3: string icon_url // user icon
}


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
    UpdateDraftBotInfoAgwResponse UpdateDraftBotInfoAgw(1:UpdateDraftBotInfoAgwRequest request)(api.post='/api/playground_api/draftbot/update_draft_bot_info', api.category="draftbot",agw.preserve_base="true")
    GetDraftBotInfoAgwResponse GetDraftBotInfoAgw(1:GetDraftBotInfoAgwRequest request)(api.post='/api/playground_api/draftbot/get_draft_bot_info', api.category="draftbot",agw.preserve_base="true")

    // public popup_info
    GetBotPopupInfoResponse GetBotPopupInfo (1:GetBotPopupInfoRequest request)(api.post='/api/playground_api/operate/get_bot_popup_info', api.category="account",agw.preserve_base="true")
    UpdateBotPopupInfoResponse UpdateBotPopupInfo (1:UpdateBotPopupInfoRequest request)(api.post='/api/playground_api/operate/update_bot_popup_info', api.category="account",agw.preserve_base="true")
    ReportUserBehaviorResponse ReportUserBehavior(1:ReportUserBehaviorRequest request)(api.post='/api/playground_api/report_user_behavior', api.category="playground_api",agw.preserve_base="true")

    // prompt resource
    prompt_resource.GetOfficialPromptResourceListResponse GetOfficialPromptResourceList(1:prompt_resource.GetOfficialPromptResourceListRequest request)(api.post='/api/playground_api/get_official_prompt_list', api.category="prompt_resource",agw.preserve_base="true")
    prompt_resource.GetPromptResourceInfoResponse GetPromptResourceInfo(1:prompt_resource.GetPromptResourceInfoRequest request)(api.get='/api/playground_api/get_prompt_resource_info', api.category="prompt_resource",agw.preserve_base="true")
    prompt_resource.UpsertPromptResourceResponse UpsertPromptResource(1:prompt_resource.UpsertPromptResourceRequest request)(api.post='/api/playground_api/upsert_prompt_resource', api.category="prompt_resource",agw.preserve_base="true")
    prompt_resource.DeletePromptResourceResponse DeletePromptResource(1:prompt_resource.DeletePromptResourceRequest request)(api.post='/api/playground_api/delete_prompt_resource', api.category="prompt_resource",agw.preserve_base="true")
}
