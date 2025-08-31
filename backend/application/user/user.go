package user

import (
	"context"
	"net/mail"
	"os"
	"slices"
	"strings"

	"github.com/kiosk404/airi-go/backend/api/model/passport"
	"github.com/kiosk404/airi-go/backend/application/base/ctxutil"
	"github.com/kiosk404/airi-go/backend/domain/user/entity"
	user "github.com/kiosk404/airi-go/backend/domain/user/service"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/types/consts"
	"github.com/kiosk404/airi-go/backend/types/errno"
	"golang.org/x/mod/sumdb/storage"
)

var UserApplicationSVC = &UserApplicationService{}

type UserApplicationService struct {
	oss       storage.Storage
	DomainSVC user.User
}

// Add a simple email verification function
func isValidAccount(email string) bool {
	// If the email string is not in the correct format, it will return an error.
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (u *UserApplicationService) PassportWebAccountRegister(ctx context.Context, locale string, req *passport.PassportWebAccountRegisterPostRequest) (
	resp *passport.PassportWebAccountRegisterPostResponse, sessionKey string, err error,
) {
	// Verify that the email format is legitimate
	if !isValidAccount(req.GetAccount()) {
		return nil, "", errorx.New(errno.ErrUserInvalidParamCode, errorx.KV("msg", "Invalid email"))
	}

	// Allow Register Checker
	if !u.allowRegisterChecker(req.GetAccount()) {
		return nil, "", errorx.New(errno.ErrNotAllowedRegisterCode)
	}

	userInfo, err := u.DomainSVC.Create(ctx, &user.CreateUserRequest{
		Account:  req.GetAccount(),
		Password: req.GetPassword(),

		Locale: locale,
	})
	if err != nil {
		return nil, "", err
	}

	userInfo, err = u.DomainSVC.Login(ctx, req.GetAccount(), req.GetPassword())
	if err != nil {
		return nil, "", err
	}

	return &passport.PassportWebAccountRegisterPostResponse{
		Data: userDo2PassportTo(userInfo),
		Code: 0,
	}, userInfo.SessionKey, nil
}

func (u *UserApplicationService) allowRegisterChecker(email string) bool {
	disableUserRegistration := os.Getenv(consts.DisableUserRegistration)
	if strings.ToLower(disableUserRegistration) != "true" {
		return true
	}

	allowedAccounts := os.Getenv(consts.AllowRegistrationAccount)
	if allowedAccounts == "" {
		return false
	}

	return slices.Contains(strings.Split(allowedAccounts, ","), strings.ToLower(email))
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

// PassportWebAccountLoginPost handle user email login requests
func (u *UserApplicationService) PassportWebAccountLoginPost(ctx context.Context, req *passport.PassportWebAccountLoginPostRequest) (
	resp *passport.PassportWebAccountLoginPostResponse, sessionKey string, err error,
) {
	userInfo, err := u.DomainSVC.Login(ctx, req.GetAccount(), req.GetPassword())
	if err != nil {
		return nil, "", err
	}

	return &passport.PassportWebAccountLoginPostResponse{
		Data: userDo2PassportTo(userInfo),
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

func (u *UserApplicationService) PassportAccountInfo(ctx context.Context, req *passport.PassportAccountInfoRequest) (
	resp *passport.PassportAccountInfoResponse, err error,
) {
	userID := ctxutil.MustGetUIDFromCtx(ctx)

	userInfo, err := u.DomainSVC.GetUserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &passport.PassportAccountInfoResponse{
		Data: userDo2PassportTo(userInfo),
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

func userDo2PassportTo(userDo *entity.User) *passport.User {
	var locale *string
	if userDo.Locale != "" {
		locale = ptr.Of(userDo.Locale)
	}

	return &passport.User{
		UserIDStr:      userDo.UserID,
		Name:           userDo.Name,
		ScreenName:     ptr.Of(userDo.Name),
		UserUniqueName: userDo.UniqueName,
		Account:        userDo.Account,
		Description:    userDo.Description,
		AvatarURL:      userDo.IconURL,
		AppUserInfo: &passport.AppUserInfo{
			UserUniqueName: userDo.UniqueName,
		},
		Locale: locale,

		UserCreateTime: userDo.CreatedAt / 1000,
	}
}
