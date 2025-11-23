package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpwrap "github.com/kiosk404/airi-go/backend/api/http"
	"github.com/kiosk404/airi-go/backend/api/model/foundation/user"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/application"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/entity"
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
