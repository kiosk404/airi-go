package service

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/storage"
	uploadEntity "github.com/kiosk404/airi-go/backend/modules/action/upload/entity"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/repo"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/pkg"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/types/errno"
)

func NewUserDomain(ctx context.Context, oss storage.Storage,
	idGen idgen.IDGenerator, userRepo repo.UserRepository) User {
	return &userImpl{
		IconOSS:  oss,
		IDGen:    idGen,
		UserRepo: userRepo,
	}
}

var _ User = &userImpl{}

type userImpl struct {
	IconOSS  storage.Storage
	IDGen    idgen.IDGenerator
	UserRepo repo.UserRepository
}

func (u *userImpl) Create(ctx context.Context, req *CreateUserRequest) (user *entity.User, err error) {
	exist, err := u.UserRepo.CheckAccountExist(ctx, req.Account)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, errorx.New(errno.ErrUserAccountAlreadyExistCode, errorx.KV("account", req.Account))
	}

	if req.UniqueName != "" {
		exist, err = u.UserRepo.CheckUniqueNameExist(ctx, req.UniqueName)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, errorx.New(errno.ErrUserUniqueNameAlreadyExistCode, errorx.KV("name", req.UniqueName))
		}
	}

	// Hashing passwords using the Argon id algorithm
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	name := req.Name

	userID, err := u.IDGen.GenID(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate id error: %w", err)
	}

	now := time.Now().UnixMilli()

	spaceID := req.SpaceID
	if spaceID <= 0 {
		var sid int64
		sid, err = u.IDGen.GenID(ctx)
		if err != nil {
			return nil, fmt.Errorf("gen space_id failed: %w", err)
		}

		spaceID = sid
	}

	newUser := &model.User{
		ID:           userID,
		IconURI:      uploadEntity.UserIconURI,
		Name:         name,
		UniqueName:   u.getUniqueNameFormAccount(ctx, req.Account),
		Account:      req.Account,
		Password:     hashedPassword,
		Description:  req.Description,
		UserVerified: false,
		Locale:       req.Locale,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = u.UserRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("insert user failed: %w", err)
	}

	iconURL, err := u.IconOSS.GetObjectUrl(ctx, newUser.IconURI)
	if err != nil {
		return nil, fmt.Errorf("get icon url failed: %w", err)
	}

	return userPo2Do(newUser, iconURL), nil
}

func (u *userImpl) Login(ctx context.Context, account, password string) (user *entity.User, err error) {
	userModel, exist, err := u.UserRepo.GetUsersByAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errorx.New(errno.ErrUserInfoInvalidateCode)
	}

	// Verify the password using the Argon id algorithm
	valid, err := verifyPassword(password, userModel.Password)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errorx.New(errno.ErrUserInfoInvalidateCode)
	}

	uniqueSessionID, err := u.IDGen.GenID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session id: %w", err)
	}

	sessionDO := &entity.Session{
		UserID:    userModel.ID,
		SessionID: uniqueSessionID,
	}

	sessionKey, err := NewSessionService().GenerateSessionKey(ctx, sessionDO)
	if err != nil {
		return nil, err
	}

	// Update user session key
	err = u.UserRepo.UpdateSessionKey(ctx, userModel.ID, sessionKey)
	if err != nil {
		return nil, err
	}

	userModel.SessionKey = sessionKey

	resURL, err := u.IconOSS.GetObjectUrl(ctx, userModel.IconURI)
	if err != nil {
		return nil, err
	}

	return userPo2Do(userModel, resURL), nil
}

func (u *userImpl) Logout(ctx context.Context, userID int64) (err error) {
	err = u.UserRepo.ClearSessionKey(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userImpl) ResetPassword(ctx context.Context, account, password string) (err error) {
	// Hashing passwords using the Argon id algorithm
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	err = u.UserRepo.UpdatePassword(ctx, account, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (u *userImpl) CreateSession(ctx context.Context, userID int64) (sessionKey string, err error) {
	uniqueSessionID, err := u.IDGen.GenID(ctx)
	if err != nil {
		return "", errorx.New(errno.ErrUserInfoInvalidateCode)
	}
	sessionDO := &entity.Session{
		UserID:    userID,
		SessionID: uniqueSessionID,
	}

	sessionKey, err = NewSessionService().GenerateSessionKey(ctx, sessionDO)
	if err != nil {
		return "", errorx.WrapByCode(err, errno.ErrUserSessionIntervalErrCode, errorx.WithExtraMsg("failed to generate session key"))
	}
	err = u.UserRepo.UpdateSessionKey(ctx, userID, sessionKey)
	if err != nil {
		return "", err
	}

	return sessionKey, nil
}

func (u *userImpl) GetUserInfo(ctx context.Context, userID int64) (user *entity.User, err error) {
	if userID <= 0 {
		return nil, errorx.New(errno.ErrUserInvalidParamCode,
			errorx.KVf("msg", "invalid user id : %d", userID))
	}

	userModel, err := u.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resURL, err := u.IconOSS.GetObjectUrl(ctx, userModel.IconURI)
	if err != nil {
		return nil, err
	}

	return userPo2Do(userModel, resURL), nil
}

func (u *userImpl) UpdateAvatar(ctx context.Context, userID int64, ext string, imagePayload []byte) (url string, err error) {
	avatarKey := "user_avatar/" + strconv.FormatInt(userID, 10) + "." + ext
	err = u.IconOSS.PutObject(ctx, avatarKey, imagePayload)
	if err != nil {
		return "", err
	}

	err = u.UserRepo.UpdateAvatar(ctx, userID, avatarKey)
	if err != nil {
		return "", err
	}

	url, err = u.IconOSS.GetObjectUrl(ctx, avatarKey)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (u *userImpl) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (user *entity.User, err error) {
	updates := map[string]interface{}{
		"updated_at": time.Now().UnixMilli(),
	}

	if req.UniqueName != nil {
		resp, err := u.ValidateProfileUpdate(ctx, &ValidateProfileUpdateRequest{
			UniqueName: req.UniqueName,
		})
		if err != nil {
			return nil, err
		}

		if resp.Code != ValidateSuccess {
			return nil, errorx.New(errno.ErrUserInvalidParamCode, errorx.KV("msg", resp.Msg))
		}

		updates["unique_name"] = ptr.From(req.UniqueName)
	}

	if req.Name != nil {
		updates["name"] = ptr.From(req.Name)
	}

	if req.Description != nil {
		updates["description"] = ptr.From(req.Description)
	}

	if req.Locale != nil {
		updates["locale"] = ptr.From(req.Locale)
	}

	err = u.UserRepo.UpdateProfile(ctx, req.UserID, updates)
	userModel, err := u.UserRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	resURL, err := u.IconOSS.GetObjectUrl(ctx, userModel.IconURI)
	if err != nil {
		return nil, err
	}

	return userPo2Do(userModel, resURL), err
}

func (u *userImpl) ValidateProfileUpdate(ctx context.Context, req *ValidateProfileUpdateRequest) (resp *ValidateProfileUpdateResponse, err error) {
	if req.UniqueName == nil && req.Account == nil {
		return nil, errorx.New(errno.ErrUserInvalidParamCode, errorx.KV("msg", "missing parameter"))
	}

	if req.UniqueName != nil {
		uniqueName := ptr.From(req.UniqueName)
		charNum := utf8.RuneCountInString(uniqueName)

		if charNum < 4 || charNum > 20 {
			return &ValidateProfileUpdateResponse{
				Code: UniqueNameTooShortOrTooLong,
				Msg:  "unique name length should be between 4 and 20",
			}, nil
		}

		exist, err := u.UserRepo.CheckUniqueNameExist(ctx, uniqueName)
		if err != nil {
			return nil, err
		}

		if exist {
			return &ValidateProfileUpdateResponse{
				Code: UniqueNameExist,
				Msg:  "unique name existed",
			}, nil
		}
	}

	return &ValidateProfileUpdateResponse{
		Code: ValidateSuccess,
		Msg:  "success",
	}, nil
}

func (u *userImpl) GetUserProfiles(ctx context.Context, userID int64) (user *entity.User, err error) {
	userInfos, err := u.MGetUserProfiles(ctx, []int64{userID})
	if err != nil {
		return nil, err
	}

	if len(userInfos) == 0 {
		return nil, errorx.New(errno.ErrUserResourceNotFound, errorx.KV("type", "user"),
			errorx.KV("id", conv.Int64ToStr(userID)))
	}

	return userInfos[0], nil
}

func (u *userImpl) MGetUserProfiles(ctx context.Context, userIDs []int64) (users []*entity.User, err error) {
	userModels, err := u.UserRepo.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	users = make([]*entity.User, 0, len(userModels))
	for _, um := range userModels {
		// Get image URL
		resURL, err := u.IconOSS.GetObjectUrl(ctx, um.IconURI)
		if err != nil {
			continue // If getting the image URL fails, skip the user
		}

		users = append(users, userPo2Do(um, resURL))
	}

	return users, nil
}

func (u *userImpl) ValidateSession(ctx context.Context, sessionKey string) (session *entity.Session, exist bool, err error) {
	// authentication session key
	sessionModel, err := NewSessionService().ValidateSession(ctx, sessionKey)
	if err != nil {
		return nil, false, errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "access denied"))
	}

	// Retrieve user information from the database
	userModel, exist, err := u.UserRepo.GetUserBySessionKey(ctx, sessionKey)
	if err != nil {
		return nil, false, err
	}

	if !exist {
		return nil, false, nil
	}

	return &entity.Session{
		UserID:    userModel.ID,
		Locale:    userModel.Locale,
		CreatedAt: sessionModel.CreatedAt,
		ExpiresAt: sessionModel.ExpiresAt,
	}, true, nil
}

func (u *userImpl) getUniqueNameFormAccount(ctx context.Context, account string) string {
	username := account
	exist, err := u.UserRepo.CheckUniqueNameExist(ctx, username)
	if err != nil {
		logs.WarnX(pkg.UserModel, "check unique name exist failed: %v", err)
		return account
	}

	if exist {
		logs.WarnX(pkg.UserModel, "unique name %s already exist", username)
		return account
	}

	return username
}
