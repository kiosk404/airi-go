package application

import (
	"context"
	"os"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/kiosk404/airi-go/backend/api/model/passport"
	"github.com/kiosk404/airi-go/backend/application/base/ctxutil"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/application/convertor"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/repo"
	user "github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/service"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/types/consts"
	"github.com/kiosk404/airi-go/backend/types/errno"
)

func InitUserService(ctx context.Context, provider rdb.Provider, oss storage.Storage, idGen idgen.IDGenerator) *UserApplicationService {
	db := provider.NewSession(ctx)
	UserApplicationSVC.DomainSVC = user.NewUserDomain(ctx,
		oss, idGen, repo.NewUserRepo(db))
	UserApplicationSVC.oss = oss

	return UserApplicationSVC
}

var UserApplicationSVC = &UserApplicationService{}

type UserApplicationService struct {
	oss       storage.Storage
	DomainSVC user.User
}

func (u *UserApplicationService) PassportWebAccountRegister(ctx context.Context, locale string, req *passport.PassportWebAccountRegisterPostRequest) (
	resp *passport.PassportWebAccountRegisterPostResponse, sessionKey string, err error) {
	// Verify that the account format is legitimate
	if !isValidAccount(req.GetAccount()) {
		return nil, "", errorx.New(errno.ErrUserInvalidParamCode, errorx.KV("msg", "Invalid account"))
	}

	// Allow Register Checker
	if !u.allowRegisterChecker(req.GetAccount()) {
		return nil, "", errorx.New(errno.ErrNotAllowedRegisterCode)
	}

	userInfo, err := u.DomainSVC.Create(ctx, &user.CreateUserRequest{
		Account:  req.GetAccount(),
		Password: req.GetPassword(),
		Locale:   locale,
	})
	if err != nil {
		return nil, "", err
	}

	userInfo, err = u.DomainSVC.Login(ctx, req.GetAccount(), req.GetPassword())
	if err != nil {
		return nil, "", err
	}

	return &passport.PassportWebAccountRegisterPostResponse{
		Data: convertor.UserDo2PassportTo(userInfo),
		Code: 0,
	}, userInfo.SessionKey, nil
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

// PassportWebLogoutGet handle user logout requests
func (u *UserApplicationService) PassportWebLogoutGet(ctx context.Context, req *passport.PassportWebLogoutGetRequest) (
	resp *passport.PassportWebLogoutGetResponse, err error,
) {
	uid := ctxutil.MustGetUIDFromCtx(ctx)

	err = u.DomainSVC.Logout(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &passport.PassportWebLogoutGetResponse{
		Code: 0,
	}, nil
}

// PassportWebAccountLoginPost handle user account login requests
func (u *UserApplicationService) PassportWebAccountLoginPost(ctx context.Context, req *passport.PassportWebAccountLoginPostRequest) (
	resp *passport.PassportWebAccountLoginPostResponse, sessionKey string, err error,
) {
	userInfo, err := u.DomainSVC.Login(ctx, req.GetAccount(), req.GetPassword())
	if err != nil {
		return nil, "", err
	}

	return &passport.PassportWebAccountLoginPostResponse{
		Data: convertor.UserDo2PassportTo(userInfo),
		Code: 0,
	}, userInfo.SessionKey, nil
}

func (u *UserApplicationService) PassportWebAccountPasswordResetGet(ctx context.Context, req *passport.PassportWebAccountPasswordResetGetRequest) (
	resp *passport.PassportWebAccountPasswordResetGetResponse, err error,
) {
	err = u.DomainSVC.ResetPassword(ctx, req.GetAccount(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &passport.PassportWebAccountPasswordResetGetResponse{
		Code: 0,
	}, nil
}

// UserUpdateAvatar Update user avatar
func (u *UserApplicationService) UserUpdateAvatar(ctx context.Context, mimeType string, req *passport.UserUpdateAvatarRequest) (
	resp *passport.UserUpdateAvatarResponse, err error,
) {
	// Get file suffix by MIME type
	var ext string
	switch mimeType {
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

	return &passport.UserUpdateAvatarResponse{
		Data: &passport.UserUpdateAvatarResponseData{
			WebURI: url,
		},
		Code: 0,
	}, nil
}

// UserUpdateProfile Update user profile
func (u *UserApplicationService) UserUpdateProfile(ctx context.Context, req *passport.UserUpdateProfileRequest) (
	resp *passport.UserUpdateProfileResponse, err error,
) {
	userID := ctxutil.MustGetUIDFromCtx(ctx)

	err = u.DomainSVC.UpdateProfile(ctx, &user.UpdateProfileRequest{
		UserID:      userID,
		Name:        req.Name,
		UniqueName:  req.UserUniqueName,
		Description: req.Description,
		Locale:      req.Locale,
	})
	if err != nil {
		return nil, err
	}

	return &passport.UserUpdateProfileResponse{
		Code: 0,
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

func (u *UserApplicationService) PassportAccountInfo(ctx context.Context, req *passport.PassportAccountInfoRequest) (
	resp *passport.PassportAccountInfoResponse, err error,
) {
	userID := ctxutil.MustGetUIDFromCtx(ctx)

	userInfo, err := u.DomainSVC.GetUserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &passport.PassportAccountInfoResponse{
		Data: convertor.UserDo2PassportTo(userInfo),
		Code: 0,
	}, nil
}

// Add a simple account verification function
// 规则：
// 1. 长度在 6 到 20 个字符之间。
// 2. 不能包含中文字符。
// 3. 只允许字母、数字和下划线。
func isValidAccount(account string) bool {
	// 规则1: 检查长度
	if len(account) < 6 || len(account) > 20 {
		return false
	}

	// 规则2: 检查是否包含中文字符
	for _, r := range account {
		if unicode.Is(unicode.Han, r) {
			return false
		}
	}

	// 规则3: 检查字符集
	// 使用正则表达式来检查是否只包含字母、数字和下划线
	// ^[a-zA-Z0-9_]+$
	// ^: 匹配字符串的开头
	// [a-zA-Z0-9_]: 匹配任何字母、数字或下划线
	// +: 匹配前一个字符集一次或多次
	// $: 匹配字符串的结尾
	regex := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	return regex.MatchString(account)
}
