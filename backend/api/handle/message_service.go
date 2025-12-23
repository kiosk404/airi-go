package handle

import (
	"context"
	"errors"
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

// BreakMessage .
// @router /api/conversation/break_message [POST]
func BreakMessage(c *gin.Context) {
	var err error
	var req message.BreakMessageRequest
	ctx := c.Request.Context()

	if err = c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	if checkErr := checkBMParams(ctx, &req); checkErr != nil {
		invalidParamRequestResponse(c, checkErr.Error())
		return
	}

	resp, err := application.ConversationSVC.BreakMessage(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteMessage .
// router /api/conversation/delete_message [POST]
func DeleteMessage(c *gin.Context) {
	var err error
	var req message.DeleteMessageRequest
	ctx := c.Request.Context()

	if err = c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	if checkErr := checkDMParams(ctx, &req); checkErr != nil {
		invalidParamRequestResponse(c, checkErr.Error())
		return
	}

	resp, err := application.ConversationSVC.DeleteMessage(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func checkDMParams(_ context.Context, req *message.DeleteMessageRequest) error {
	if req.MessageID <= 0 {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "message id is invalid"))
	}

	return nil
}

func checkBMParams(_ context.Context, req *message.BreakMessageRequest) error {
	if req.AnswerMessageID == nil {
		return errors.New("answer message id is required")
	}
	if *req.AnswerMessageID <= 0 {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "answer message id is invalid"))
	}

	return nil
}
