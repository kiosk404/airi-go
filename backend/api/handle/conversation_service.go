package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/api/model/conversation/conversation"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/application"
)

// ListConversationsApi .
// @router /v1/conversations [GET]
func ListConversationsApi(c *gin.Context) {
	var err error
	var req conversation.ListConversationsApiRequest
	ctx := c.Request.Context()

	if err := c.ShouldBindQuery(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := application.ConversationSVC.ListConversation(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ClearConversationHistory .
// @router /api/conversation/clear_message [POST]
func ClearConversationHistory(c *gin.Context) {
	var err error
	var req conversation.ClearConversationHistoryRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := application.ConversationSVC.ClearHistory(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// CreateConversation .
// @router /api/conversation/create [POST]
func CreateConversation(c *gin.Context) {
	var err error
	var req conversation.CreateConversationRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := application.ConversationSVC.CreateConversation(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
