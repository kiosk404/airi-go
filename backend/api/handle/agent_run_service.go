package handle

import (
	"context"
	"encoding/json"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/api/model/conversation/run"
	"github.com/kiosk404/airi-go/backend/modules/conversation"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/pkg/errno"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	sseimpl "github.com/kiosk404/airi-go/backend/pkg/http/sse"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

// AgentRun .
// @router /api/conversation/chat [POST]
func AgentRun(c *gin.Context) {
	var err error
	var req run.AgentRunRequest

	ctx := c.Request.Context()

	// 绑定并校验参数
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	if checkErr := checkParams(ctx, &req); checkErr != nil {
		invalidParamRequestResponse(c, checkErr.Error())
		return
	}

	sseSender := sseimpl.NewSSESender(c)

	err = conversation.ConversationSVC.Run(ctx, sseSender, &req)
	if err != nil {
		errData := run.ErrorData{
			Code: errno.ErrConversationAgentRunError,
			Msg:  err.Error(),
		}
		ed, _ := json.Marshal(errData)
		_ = sseSender.Send(ctx, &sse.Event{
			Event: run.RunEventError,
			Data:  ed,
		})
	}
}

func checkParams(_ context.Context, ar *run.AgentRunRequest) error {
	if ar.BotID == 0 {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "bot id is required"))
	}

	if ar.Scene == nil {
		return errorx.New(errno.ErrConversationInvalidParamCode, errorx.KV("msg", "scene is required"))
	}

	if ar.ContentType == nil {
		ar.ContentType = ptr.Of(run.ContentTypeText)
	}
	return nil
}
