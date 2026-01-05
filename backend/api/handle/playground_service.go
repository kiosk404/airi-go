package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/api/model/playground"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/application/singleagent"
)

// GetDraftBotInfoAgw .
// @router /api/playground_api/draftbot/get_draft_bot_info [POST]
func GetDraftBotInfoAgw(c *gin.Context) {
	var err error
	var req playground.GetDraftBotInfoAgwRequest
	ctx := c.Request.Context()
	if err = c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	if req.BotID == 0 {
		invalidParamRequestResponse(c, "bot id is nil")
		return
	}
	resp, err := singleagent.SingleAgentSVC.GetAgentBotInfo(ctx, &req)
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
