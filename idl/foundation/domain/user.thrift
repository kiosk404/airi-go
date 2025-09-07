namespace go foundation.domain.user


struct AppUserInfo {
    1: required string user_unique_name
}

// UserInfoDetail 用户详细信息，包含姓名、头像等
struct UserInfoDetail {
    1: required string name                     // 唯一名称
    2: required string nick_name                // 用户昵称
    3: required string avatar_url               // 用户头像url
    4: required string account                  // 用户账号
    5: optional string mobile                   // 手机号
    6: optional string description              // 用户描述
    7: optional string user_id (agw.js_conv="str", api.js_conv="true") // 用户在租户内的唯一标识
    8: optional AppUserInfo app_user_info       // 唯一标识

    10: i64 user_create_time                    // unix timestamp in seconds
}

typedef string UserStatus (ts.enum="true")
const UserStatus active = "active"
const UserStatus deactivated = "deactivated"
const UserStatus offboarded = "offboarded"

