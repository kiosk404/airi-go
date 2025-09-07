package application

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/kiosk404/airi-go/backend/api/model/app/bot_open_api"
	openapimodel "github.com/kiosk404/airi-go/backend/api/model/permission/openapiauth"
	"github.com/kiosk404/airi-go/backend/application/base/ctxutil"
	"github.com/kiosk404/airi-go/backend/infra/contract/idgen"
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/modules/foundation/openauth/domain/entity"
	openapiauth "github.com/kiosk404/airi-go/backend/modules/foundation/openauth/domain/service"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/lang/slices"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

var (
	openapiAuthDomainSVC openapiauth.APIAuth
)

func InitService(db rdb.Provider, idGenSVC idgen.IDGenerator) *OpenAuthApplicationService {
	openapiAuthDomainSVC = openapiauth.NewService(
		idGenSVC, db)

	OpenAuthApplication.OpenAPIDomainSVC = openapiAuthDomainSVC

	return OpenAuthApplication
}

type OpenAuthApplicationService struct {
	OpenAPIDomainSVC openapiauth.APIAuth
}

var OpenAuthApplication = &OpenAuthApplicationService{}

func (s *OpenAuthApplicationService) GetPersonalAccessTokenAndPermission(ctx context.Context, req *openapimodel.GetPersonalAccessTokenAndPermissionRequest) (*openapimodel.GetPersonalAccessTokenAndPermissionResponse, error) {

	resp := new(openapimodel.GetPersonalAccessTokenAndPermissionResponse)

	userID := ctxutil.GetUIDFromCtx(ctx)

	appReq := &entity.GetApiKey{
		ID: req.ID,
	}
	apiKeyResp, err := openapiAuthDomainSVC.Get(ctx, appReq)

	if err != nil {
		logs.Error("OpenAuthApplicationService.GetPersonalAccessTokenAndPermission failed, err=%v", err)
		return resp, errors.New("GetPersonalAccessTokenAndPermission failed")
	}
	if apiKeyResp == nil {
		return resp, errors.New("GetPersonalAccessTokenAndPermission failed")
	}

	if apiKeyResp.UserID != *userID {
		return resp, errors.New("permission not match")
	}
	resp.Data = &openapimodel.GetPersonalAccessTokenAndPermissionResponseData{
		PersonalAccessToken: &openapimodel.PersonalAccessToken{
			ID:        apiKeyResp.ID,
			Name:      apiKeyResp.Name,
			ExpireAt:  apiKeyResp.ExpiredAt,
			CreatedAt: apiKeyResp.CreatedAt,
			UpdatedAt: apiKeyResp.UpdatedAt,
		},
	}
	return resp, nil
}

func (s *OpenAuthApplicationService) CreatePersonalAccessToken(ctx context.Context, req *openapimodel.CreatePersonalAccessTokenAndPermissionRequest) (*openapimodel.CreatePersonalAccessTokenAndPermissionResponse, error) {
	resp := new(openapimodel.CreatePersonalAccessTokenAndPermissionResponse)
	userID := ctxutil.GetUIDFromCtx(ctx)

	appReq := &entity.CreateApiKey{
		Name:   req.Name,
		Expire: req.ExpireAt,
		UserID: *userID,
		AkType: entity.AkTypeCustomer,
	}

	if req.DurationDay == "customize" {
		appReq.Expire = req.ExpireAt
	} else {
		expireDay, err := strconv.ParseInt(req.DurationDay, 10, 64)
		if err != nil {
			return resp, errors.New("invalid expireDay")
		}
		appReq.Expire = time.Now().Add(time.Duration(expireDay) * time.Hour * 24).Unix()
	}

	apiKeyResp, err := openapiAuthDomainSVC.Create(ctx, appReq)
	if err != nil {
		logs.Error("OpenAuthApplicationService.CreatePersonalAccessToken failed, err=%v", err)
		return resp, errors.New("CreatePersonalAccessToken failed")
	}
	resp.Data = &openapimodel.CreatePersonalAccessTokenAndPermissionResponseData{
		PersonalAccessToken: &openapimodel.PersonalAccessToken{
			ID:        apiKeyResp.ID,
			Name:      apiKeyResp.Name,
			ExpireAt:  apiKeyResp.ExpiredAt,
			CreatedAt: apiKeyResp.CreatedAt,
			UpdatedAt: apiKeyResp.UpdatedAt,
		},
		Token: apiKeyResp.ApiKey,
	}
	return resp, nil
}

func (s *OpenAuthApplicationService) ImpersonateAiriUserAccessToken(ctx context.Context, req *bot_open_api.ImpersonateAiriUserRequest) (*bot_open_api.ImpersonateAiriUserResponse, error) {
	resp := new(bot_open_api.ImpersonateAiriUserResponse)
	userID := ctxutil.GetUIDFromCtx(ctx)

	expiredSecond := time.Now().Add(time.Duration(time.Second * 60 * 15)).Unix()

	appReq := &entity.CreateApiKey{
		UserID: *userID,
		AkType: entity.AkTypeTemporary,
		Expire: expiredSecond,
		Name:   "temporary access token",
	}

	apiKeyResp, err := openapiAuthDomainSVC.Create(ctx, appReq)
	if err != nil {
		logs.Error("OpenAuthApplicationService.CreatePersonalAccessToken failed, err=%v", err)
		return resp, errors.New("CreatePersonalAccessToken failed")
	}
	resp.Data = &bot_open_api.ImpersonateAiriUserResponseData{
		AccessToken: apiKeyResp.ApiKey,
		ExpiresIn:   expiredSecond,
		TokenType:   "Bearer",
	}
	return resp, nil
}

func (s *OpenAuthApplicationService) ListPersonalAccessTokens(ctx context.Context, req *openapimodel.ListPersonalAccessTokensRequest) (*openapimodel.ListPersonalAccessTokensResponse, error) {

	resp := new(openapimodel.ListPersonalAccessTokensResponse)

	userID := ctxutil.GetUIDFromCtx(ctx)
	appReq := &entity.ListApiKey{
		UserID: *userID,
		Page:   *req.Page,
		Limit:  *req.Size,
	}

	apiKeyResp, err := openapiAuthDomainSVC.List(ctx, appReq)
	if err != nil {
		logs.Error("OpenAuthApplicationService.ListPersonalAccessTokens failed, err=%v", err)
		return resp, errors.New("ListPersonalAccessTokens failed")
	}

	if apiKeyResp == nil {
		return resp, nil
	}
	resp.Data = &openapimodel.ListPersonalAccessTokensResponseData{
		HasMore: apiKeyResp.HasMore,
		PersonalAccessTokens: slices.Transform(apiKeyResp.ApiKeys, func(a *entity.ApiKey) *openapimodel.PersonalAccessTokenWithCreatorInfo {
			lastUsedAt := a.LastUsedAt
			if lastUsedAt == 0 {
				lastUsedAt = -1
			}
			return &openapimodel.PersonalAccessTokenWithCreatorInfo{
				ID:         a.ID,
				Name:       a.Name,
				ExpireAt:   a.ExpiredAt,
				CreatedAt:  a.CreatedAt,
				UpdatedAt:  a.UpdatedAt,
				LastUsedAt: lastUsedAt,
			}
		}),
	}
	return resp, nil
}

func (s *OpenAuthApplicationService) DeletePersonalAccessTokenAndPermission(ctx context.Context, req *openapimodel.DeletePersonalAccessTokenAndPermissionRequest) (*openapimodel.DeletePersonalAccessTokenAndPermissionResponse, error) {
	resp := new(openapimodel.DeletePersonalAccessTokenAndPermissionResponse)

	userID := ctxutil.GetUIDFromCtx(ctx)

	appReq := &entity.DeleteApiKey{
		ID:     req.ID,
		UserID: *userID,
	}
	err := openapiAuthDomainSVC.Delete(ctx, appReq)
	if err != nil {
		logs.Error("OpenAuthApplicationService.DeletePersonalAccessTokenAndPermission failed, err=%v", err)
		return resp, errors.New("DeletePersonalAccessTokenAndPermission failed")
	}
	return resp, nil
}

func (s *OpenAuthApplicationService) UpdatePersonalAccessTokenAndPermission(ctx context.Context, req *openapimodel.UpdatePersonalAccessTokenAndPermissionRequest) (*openapimodel.UpdatePersonalAccessTokenAndPermissionResponse, error) {
	resp := new(openapimodel.UpdatePersonalAccessTokenAndPermissionResponse)
	userID := ctxutil.GetUIDFromCtx(ctx)

	upErr := openapiAuthDomainSVC.Save(ctx, &entity.SaveMeta{
		ID:     req.ID,
		Name:   ptr.Of(req.Name),
		UserID: *userID,
	})

	return resp, upErr
}

func (s *OpenAuthApplicationService) UpdateLastUsedAt(ctx context.Context, apiID int64, userID int64) error {
	upErr := openapiAuthDomainSVC.Save(ctx, &entity.SaveMeta{
		ID:         apiID,
		LastUsedAt: ptr.Of(time.Now().Unix()),
		UserID:     userID,
	})
	return upErr
}

func (s *OpenAuthApplicationService) CheckPermission(ctx context.Context, token string) (*entity.ApiKey, error) {
	appReq := &entity.CheckPermission{
		ApiKey: token,
	}
	apiKey, err := openapiAuthDomainSVC.CheckPermission(ctx, appReq)
	if err != nil {
		logs.Error("OpenAuthApplicationService.CheckPermission failed, err=%v", err)
		return nil, errors.New("CheckPermission failed")
	}
	return apiKey, nil
}
