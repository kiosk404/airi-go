package handle

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	httpwrap "github.com/kiosk404/airi-go/backend/api/http"
	"github.com/kiosk404/airi-go/backend/api/model/foundation/user"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/application"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/types/consts"
)

// PassportWebRegisterPost .
// @router /foundation/v1/users/register [POST]
func PassportWebRegisterPost(c *gin.Context) {
	var err error
	var req user.UserRegisterRequest
	ctx := c.Request.Context()

	// 绑定并校验参数
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := application.UserApplicationSVC.WebAccountRegister(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.SetCookie(entity.SessionKey,
		resp.GetToken(),
		consts.SessionMaxAgeSecond,
		"/", httpwrap.GetOriginHost(c),
		false, true)

	c.JSON(http.StatusOK, resp)
}

// PassportWebLoginByPasswordPost .
// @router /foundation/v1/users/login [POST]
func PassportWebLoginByPasswordPost(c *gin.Context) {
	var err error
	var req user.LoginByPasswordRequest
	ctx := c.Request.Context()

	// 绑定并校验参数
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := application.UserApplicationSVC.WebAccountLoginByPassword(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	logs.Info("[PassportWebLoginPost] sessionKey: %s", resp.GetToken())

	c.SetCookie(entity.SessionKey,
		resp.GetToken(),
		consts.SessionMaxAgeSecond,
		"/", httpwrap.GetOriginHost(c),
		false, true)

	c.JSON(http.StatusOK, resp)
}

// PassportWebLogoutPost .
// @router /foundation/v1/users/logout [POST]
func PassportWebLogoutPost(c *gin.Context) {
	var err error
	var req user.LogoutRequest
	ctx := c.Request.Context()

	// 绑定并校验参数
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := application.UserApplicationSVC.WebLogout(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PassportAccountInfo .
// @router /foundation/v1/users/session [GET]
func PassportAccountInfo(c *gin.Context) {
	var err error
	ctx := c.Request.Context()

	req := user.NewGetUserInfoByTokenRequest()

	resp, err := application.UserApplicationSVC.GetUserInfoByToken(ctx, req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UserUpdateAvatar .
// @router /api/foundation/v1/users/:user_id:/upload_avatar/ [POST]
func UserUpdateAvatar(c *gin.Context) {
	var err error
	var req user.UserUpdateAvatarRequest
	ctx := c.Request.Context()

	// Get the uploaded file
	file, err := c.FormFile("avatar")
	if err != nil {
		logs.Error("Get Avatar Fail failed, err=%v", err)
		invalidParamRequestResponse(c, "missing avatar file")
		return
	}

	// Check file type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		invalidParamRequestResponse(c, "invalid file type, only image allowed")
		return
	}

	// Read file content
	src, err := file.Open()
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}
	defer src.Close()

	fileContent, err := io.ReadAll(src)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	req.Avatar = fileContent
	resp, err := application.UserApplicationSVC.UserUpdateAvatar(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
