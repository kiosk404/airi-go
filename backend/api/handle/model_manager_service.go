package handle

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/api/model/modelapi"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/application/singleagent"
	modelmgrapp "github.com/kiosk404/airi-go/backend/modules/llm/application"
	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
)

// CreateModel .
// @router /api/admin/model/create [POST]
func CreateModel(c *gin.Context) {
	var err error
	var req modelapi.CreateModelReq
	var modelID int64
	ctx := c.Request.Context()
	// 绑定请求参数
	if err = c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	if modelID, err = modelmgrapp.ModelMgrSVC.CreateModel(ctx, req); err != nil {
		invalidParamRequestResponse(c, fmt.Sprintf("generate model failed: %v", err))
		return
	}

	resp := new(modelapi.CreateModelResp)
	resp.ID = conv.Int64ToStr(modelID)
	c.JSON(http.StatusOK, resp)
}

// UpdateModel .
// @router /api/admin/model/update [POST]
func UpdateModel(c *gin.Context) {
	var err error
	var req modelapi.UpdateModelReq
	ctx := c.Request.Context()
	if err = c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	if err = modelmgrapp.ModelMgrSVC.UpdateModel(ctx, req); err != nil {
		invalidParamRequestResponse(c, fmt.Sprintf("update model failed: %v", err))
		return
	}
	resp := new(modelapi.UpdateModelResp)
	c.JSON(http.StatusOK, resp)
}

// DeleteModel .
// @router /api/admin/model/delete [POST]
func DeleteModel(c *gin.Context) {
	var err error
	var req modelapi.DeleteModelReq
	ctx := c.Request.Context()
	if err = c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	if err = modelmgrapp.ModelMgrSVC.DeleteModel(ctx, req); err != nil {
		invalidParamRequestResponse(c, fmt.Sprintf("delete model failed: %v", err))
		return
	}
	resp := new(modelapi.DeleteModelResp)
	c.JSON(http.StatusOK, resp)
}

// SetDefaultModel .
// @router /api/admin/model/set_default [POST]
func SetDefaultModel(c *gin.Context) {
	var err error
	var req modelapi.SetDefaultModelReq
	ctx := c.Request.Context()
	if err = c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	if err = modelmgrapp.ModelMgrSVC.SetDefaultModel(ctx, &req); err != nil {
		invalidParamRequestResponse(c, fmt.Sprintf("set default model failed: %v", err))
		return
	}

	if err = singleagent.SingleAgentSVC.UpdateAgentModelInfo(ctx, req.GetID()); err != nil {
		invalidParamRequestResponse(c, fmt.Sprintf("set default model failed: %v", err))
		return
	}

	resp := new(modelapi.SetDefaultModelResp)
	c.JSON(http.StatusOK, resp)
}

// GetModelList .
// @router /api/admin/model/list [GET]
func GetModelList(c *gin.Context) {
	var err error
	var req modelapi.GetModelListReq
	ctx := c.Request.Context()
	modeList, err := modelmgrapp.ModelMgrSVC.GetInUseModelList(ctx, req.GetModelType())
	if err != nil {
		invalidParamRequestResponse(c, fmt.Sprintf("get builtin model list failed: %v", err))
		return
	}
	resp := new(modelapi.GetModelListResp)
	resp.ProviderModelList = modeList
	c.JSON(http.StatusOK, resp)
}
