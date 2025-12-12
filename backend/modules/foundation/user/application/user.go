package application

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/bytedance/gg/gslice"
	domain "github.com/kiosk404/airi-go/backend/api/model/foundation/domain/user"
	"github.com/kiosk404/airi-go/backend/api/model/foundation/user"
	"github.com/kiosk404/airi-go/backend/application/ctxutil"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/application/convertor"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/repo"
	usersvc "github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

var UserApplicationSVC = &UserApplicationService{}

func InitService(ctx context.Context, provider rdb.Provider, oss storage.Storage, idGen idgen.IDGenerator) *UserApplicationService {
	db := provider.NewSession(ctx)
	UserApplicationSVC.DomainSVC = usersvc.NewUserDomain(ctx,
		oss, idGen, repo.NewUserRepo(db))
	UserApplicationSVC.oss = oss

	return UserApplicationSVC
}

type UserApplicationService struct {
	oss       storage.Storage
	DomainSVC usersvc.User
}

var _ user.UserService = &UserApplicationService{}

func (u *UserApplicationService) WebAccountRegister(ctx context.Context,
	req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	// Verify that the account format is legitimate
	if !isValidAccount(req.GetAccount()) {
		return nil, errorx.New(errno.ErrUserInvalidParamCode, errorx.KV("msg", "Invalid account"))
	}

	// Allow Register Checker
	if !u.allowRegisterChecker(req.GetAccount()) {
		return nil, errorx.New(errno.ErrNotAllowedRegisterCode)
	}

	userDO, err := u.DomainSVC.Create(ctx, &usersvc.CreateUserRequest{
		Account:  req.GetAccount(),
		Password: req.GetPassword(),
		Locale:   req.GetLocale(),
	})
	if err != nil {
		return nil, err
	}

	userInfo, err := u.DomainSVC.Login(ctx, req.GetAccount(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &user.UserRegisterResponse{
		UserInfo:   convertor.UserDO2DTO(userDO),
		Token:      ptr.Of(userInfo.SessionKey),
		ExpireTime: ptr.Of(int64(entity.SessionExpires)),
	}, nil
}

func (u *UserApplicationService) allowRegisterChecker(account string) bool {
	disableUserRegistration := os.Getenv(consts.DisableUserRegistration)
	if strings.ToLower(disableUserRegistration) != "true" {
		return true
	}

	allowedAccounts := os.Getenv(consts.AllowRegistrationAccount)
	if allowedAccounts == "" {
		return false
	}

	return slices.Contains(strings.Split(allowedAccounts, ","), strings.ToLower(account))
}

// WebLogout handle user logout requests
func (u *UserApplicationService) WebLogout(ctx context.Context,
	req *user.LogoutRequest) (resp *user.LogoutResponse, err error) {
	uid := ctxutil.MustGetUIDFromCtx(ctx)

	err = u.DomainSVC.Logout(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &user.LogoutResponse{}, nil
}

// WebAccountLoginByPassword handle user account login requests
func (u *UserApplicationService) WebAccountLoginByPassword(ctx context.Context,
	req *user.LoginByPasswordRequest) (resp *user.LoginByPasswordResponse, err error) {
	userDO, err := u.DomainSVC.Login(ctx, req.GetAccount(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &user.LoginByPasswordResponse{
		UserInfo:   convertor.UserDO2DTO(userDO),
		Token:      ptr.Of(userDO.SessionKey),
		ExpireTime: ptr.Of(int64(entity.SessionExpires)),
	}, nil
}

func (u *UserApplicationService) WebAccountPasswordReset(ctx context.Context,
	req *user.ResetPasswordRequest) (resp *user.ResetPasswordResponse, err error) {
	err = u.DomainSVC.ResetPassword(ctx, req.GetAccount(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &user.ResetPasswordResponse{}, nil
}

// UserUpdateAvatar Update user avatar
func (u *UserApplicationService) UserUpdateAvatar(ctx context.Context,
	req *user.UserUpdateAvatarRequest) (resp *user.UserUpdateAvatarResponse, err error) {
	// Get file suffix by MIME type
	var ext string
	switch req.GetMimeType() {
	case "image/jpeg", "image/jpg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	case "image/gif":
		ext = "gif"
	case "image/webp":
		ext = "webp"
	default:
		return nil, errorx.WrapByCode(err, errno.ErrUserInvalidParamCode,
			errorx.KV("msg", "unsupported image type"))
	}

	uid := ctxutil.MustGetUIDFromCtx(ctx)
	url, err := u.DomainSVC.UpdateAvatar(ctx, uid, ext, req.GetAvatar())
	if err != nil {
		return nil, err
	}
	return &user.UserUpdateAvatarResponse{
		WebURI: url,
	}, nil
}

// ModifyUserProfile Update user profile
func (u *UserApplicationService) ModifyUserProfile(ctx context.Context,
	req *user.ModifyUserProfileRequest) (resp *user.ModifyUserProfileResponse, err error) {
	userID := ctxutil.MustGetUIDFromCtx(ctx)

	userDO, err := u.DomainSVC.UpdateProfile(ctx, &usersvc.UpdateProfileRequest{
		UserID:      userID,
		Name:        req.NickName,
		UniqueName:  req.Name,
		Description: req.Description,
		Locale:      req.Locale,
	})
	if err != nil {
		return nil, err
	}
	return &user.ModifyUserProfileResponse{
		UserInfo: convertor.UserDO2DTO(userDO),
	}, nil
}

func (u *UserApplicationService) ValidateSession(ctx context.Context, sessionKey string) (*entity.Session, error) {
	session, exist, err := u.DomainSVC.ValidateSession(ctx, sessionKey)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "session not exist"))
	}

	return session, nil
}

func (u *UserApplicationService) GetUserInfo(ctx context.Context,
	req *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	userID := ctxutil.MustGetUIDFromCtx(ctx)

	userDO, err := u.DomainSVC.GetUserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &user.GetUserInfoResponse{
		UserInfo: convertor.UserDO2DTO(userDO),
	}, err
}

func (u *UserApplicationService) GetUserInfoByToken(ctx context.Context, request *user.GetUserInfoByTokenRequest) (r *user.GetUserInfoByTokenResponse, err error) {
	userID := ctxutil.MustGetUIDFromCtx(ctx)

	userDO, err := u.DomainSVC.GetUserProfiles(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &user.GetUserInfoByTokenResponse{
		UserInfo: convertor.UserDO2DTO(userDO),
	}, nil
}

func (u *UserApplicationService) MGetUserInfo(ctx context.Context, req *user.MGetUserInfoRequest) (resp *user.MGetUserInfoResponse, err error) {
	resp = user.NewMGetUserInfoResponse()
	if len(req.GetUserIds()) == 0 {
		return nil, errorx.NewByCode(errno.ErrUserInvalidParamCode, errorx.WithExtraMsg("user id is empty"))
	}
	userIDs, err := gslice.TryMap(req.GetUserIds(), func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	}).Get()
	if err != nil {
		return nil, errorx.NewByCode(errno.ErrUserInvalidParamCode, errorx.WithExtraMsg("user id is invalid"))
	}

	userDOs, err := u.DomainSVC.MGetUserProfiles(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	resp.UserInfos = slices.Map(userDOs, func(userDO *entity.User, _ int) *domain.UserInfoDetail {
		return convertor.UserDO2DTO(userDO)
	})
	return resp, nil
}

// Add a simple account verification function
// 规则：
// 1. 长度在 2 到 20 个字符之间。
// 2. 不能包含中文字符。
// 3. 只允许字母、数字和下划线 或 邮箱
func isValidAccount(account string) bool {
	// 规则1: 检查长度
	if len(account) < 2 || len(account) > 20 {
		return false
	}

	// 规则2: 检查是否包含中文字符
	for _, r := range account {
		if unicode.Is(unicode.Han, r) {
			return false
		}
	}

	// 规则3: 检查字符集
	// 使用正则表达式来检查是否只包含字母、数字和下划线  或 邮箱
	// ^[a-zA-Z0-9_]+$
	// ^: 匹配字符串的开头
	// [a-zA-Z0-9_]: 匹配任何字母、数字或下划线
	// +: 匹配前一个字符集一次或多次
	// $: 匹配字符串的结尾
	StrRegex := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return StrRegex.MatchString(account) || emailRegex.MatchString(account)
}
