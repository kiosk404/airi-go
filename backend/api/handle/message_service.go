package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/api/model/conversation/message"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/application"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
)

// GetMessageList .
// @router /api/conversation/get_message_list [POST]
func GetMessageList(c *gin.Context) {
	var err error
	var req message.GetMessageListRequest
	ctx := c.Request.Context()
	// 绑定请求参数
	if err = c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 检查ML参数
	if req.BotID == "" {
		checkErr := errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "agent id is required"))
		invalidParamRequestResponse(c, checkErr.Error())
		return
	}

	resp, err := application.ConversationSVC.GetMessageList(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
