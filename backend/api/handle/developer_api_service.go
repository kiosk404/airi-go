package handle

import (
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/api/model/app/developer_api"
	"github.com/kiosk404/airi-go/backend/api/model/playground"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/application/singleagent"
)

// DraftBotCreate .
// @router /api/draftbot/create [POST]
func DraftBotCreate(c *gin.Context) {
	var err error
	var req developer_api.DraftBotCreateRequest
	ctx := c.Request.Context()

	// 绑定并校验参数
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	if req.Name == "" {
		invalidParamRequestResponse(c, "name is nil")
		return
	}

	if req.IconURI == "" {
		invalidParamRequestResponse(c, "icon uri is nil")
		return
	}

	if utf8.RuneCountInString(req.Name) > 50 {
		invalidParamRequestResponse(c, "name is too long")
		return
	}

	if utf8.RuneCountInString(req.Description) > 2000 {
		invalidParamRequestResponse(c, "description is too long")
		return
	}

	resp, err := singleagent.SingleAgentSVC.CreateSingleAgentDraft(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DraftBotUpdateInfo .
// @router /api/playground_api/draftbot/update_draft_bot_info [POST]
func DraftBotUpdateInfo(c *gin.Context) {
	var req playground.UpdateDraftBotInfoAgwRequest
	ctx := c.Request.Context()

	// 绑定并校验参数
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	if req.BotInfo == nil {
		invalidParamRequestResponse(c, "bot info is nil")
		return
	}

	if req.BotInfo.BotId == nil {
		invalidParamRequestResponse(c, "bot id is nil")
		return
	}

	resp, err := singleagent.SingleAgentSVC.UpdateSingleAgentDraft(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
