namespace go passport

struct AppUserInfo {
    1: required string user_unique_name
}

struct User {
    // Align with the original interface field name
    1: required i64 user_id_str (agw.js_conv="str", api.js_conv="true")
    2: required string name
    3: required string user_unique_name
    4: required string account
    5: required string description
    6: required string avatar_url
    7: optional string screen_name
    8: optional AppUserInfo app_user_info
    9: optional string locale

    10: i64 user_create_time // unix timestamp in seconds
}

struct PassportWebAccountRegisterPostRequest {
    11: required string password

    23: string account
}

struct PassportWebAccountRegisterPostResponse {

    1: required User data

    253: required i32            code
    254: required string         msg
}

struct PassportWebLogoutGetRequest {
}

struct PassportWebLogoutGetResponse {
    1: required string redirect_url

    253: required i32            code
    254: required string         msg
}


struct PassportWebAccountLoginPostRequest {
    6: required string account
    7: required string password
}

struct PassportWebAccountLoginPostResponse {
    1: required User data

    253: required i32            code
    254: required string         msg
}

struct PassportWebAccountPasswordResetGetRequest {
    1: string password
    2: string code
    3: string account
}

struct PassportWebAccountPasswordResetGetResponse {
    253: required i32            code
    254: required string         msg
}

struct PassportAccountInfoRequest {}
struct PassportAccountInfoResponse {
    1: required User data

    253: required i32            code
    254: required string         msg
}

struct UserUpdateAvatarRequest {
    3: required binary avatar (api.form="avatar")
}

struct UserUpdateAvatarResponseData {
    1: required string web_uri
}

struct UserUpdateAvatarResponse {
    1: required UserUpdateAvatarResponseData data

    253: required i32            code
    254: required string         msg
}

struct UserUpdateProfileRequest {
    2: optional string name
    3: optional string user_unique_name
    5: optional string description
    6: optional string locale
}

struct UserUpdateProfileResponse {

    253: required i32            code
    254: required string         msg
}

service PassportService {

    // Account password registration
    PassportWebAccountRegisterPostResponse PassportWebAccountRegisterPost(1: PassportWebAccountRegisterPostRequest req) (api.post="/api/passport/web/account/register/")

    // log out
    PassportWebLogoutGetResponse PassportWebLogoutGet(1: PassportWebLogoutGetRequest req) (api.get="/api/passport/web/logout/")

    // Account account password login
    PassportWebAccountLoginPostResponse PassportWebAccountLoginPost(1: PassportWebAccountLoginPostRequest req) (api.post="/api/passport/web/account/login/")


    // Reset password via Account
    PassportWebAccountPasswordResetGetResponse PassportWebAccountPasswordResetGet(1: PassportWebAccountPasswordResetGetRequest req) (api.get="/api/passport/web/account/password/reset/")

    // account information
    PassportAccountInfoResponse PassportAccountInfo(1: PassportAccountInfoRequest req) (api.post="/api/passport/account/info/")


    UserUpdateAvatarResponse UserUpdateAvatar(1: UserUpdateAvatarRequest req) (api.post="/api/web/user/update/upload_avatar/", api.serializer="form")

    UserUpdateProfileResponse UserUpdateProfile(1: UserUpdateProfileRequest req) (api.post="/api/user/update_profile")
}