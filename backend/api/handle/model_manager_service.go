package handle

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/api/model/modelapi"
	modelmgrapp "github.com/kiosk404/airi-go/backend/modules/llm/application"
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
	resp.ID = modelID
	c.JSON(http.StatusOK, resp)
}

// DeleteModel .
// @router /api/admin/model/delete [POST]
func DeleteModel(c *gin.Context) {

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
