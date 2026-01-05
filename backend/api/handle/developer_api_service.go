package handle

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/api/model/app/developer_api"
	"github.com/kiosk404/airi-go/backend/application/ctxutil"
	"github.com/kiosk404/airi-go/backend/modules/component/agent/application/singleagent"
	uploadapp "github.com/kiosk404/airi-go/backend/modules/data/upload/application"
	"github.com/kiosk404/airi-go/backend/modules/data/upload/pkg/errno"
	modelmgr "github.com/kiosk404/airi-go/backend/modules/llm/application"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

// CheckDraftBotCommit .
// @router /api/draftbot/commit_check [POST]
func CheckDraftBotCommit(c *gin.Context) {

}

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

// DeleteBotDelete .
// @router /api/draftbot/delete [POST]
func DeleteBotDelete(c *gin.Context) {
	var err error
	var req developer_api.DeleteDraftBotRequest
	ctx := c.Request.Context()

	if err = c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := singleagent.SingleAgentSVC.DeleteAgentDraft(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDraftBotDisplayInfo .
// @router /api/draftbot/get_display_info [POST]
func GetDraftBotDisplayInfo(c *gin.Context) {
	var err error
	var req developer_api.GetDraftBotDisplayInfoRequest
	ctx := c.Request.Context()
	// 绑定并校验参数
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp, err := singleagent.SingleAgentSVC.GetAgentDraftDisplayInfo(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateDraftBotDisplayInfo .
// @router /api/draftbot/update_display_info [POST]
func UpdateDraftBotDisplayInfo(c *gin.Context) {
	var err error
	var req developer_api.UpdateDraftBotDisplayInfoRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := singleagent.SingleAgentSVC.UpdateAgentDraftDisplayInfo(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UploadFile .
// @router /api/bot/upload_file [POST]
func UploadFile(c *gin.Context) {
	var err error
	ctx := c.Request.Context()

	var req developer_api.UploadFileRequest
	if err = c.ShouldBind(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp := new(developer_api.UploadFileResponse)
	fileContent, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	userID := ctxutil.GetUIDFromCtx(ctx)
	if userID == nil {
		internalServerErrorResponse(c, errorx.New(errno.ErrUploadPermissionCode, errorx.KV("msg", "session required")))
		return
	}
	secret := createSecret(ptr.From(userID), req.FileHead.FileType)
	fileName := fmt.Sprintf("%d_%d_%s.%s", ptr.From(userID), time.Now().UnixNano(), secret, req.FileHead.FileType)
	objectName := fmt.Sprintf("%s/%s", req.FileHead.BizType.String(), fileName)
	resp, err = uploadapp.SVC.UploadFile(ctx, fileContent, objectName)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetTypeList .
// @router /api/bot/get_type_list [POST]
func GetTypeList(c *gin.Context) {
	var err error
	var req developer_api.GetTypeListRequest
	ctx := c.Request.Context()
	// 绑定并校验参数
	if err = c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := modelmgr.ModelMgrSVC.GetModelList(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

const baseWord = "1Aa2Bb3Cc4Dd5Ee6Ff7Gg8Hh9Ii0JjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"

func createSecret(uid int64, fileType string) string {
	num := 10
	input := fmt.Sprintf("upload_%d_Ma*9)fhi_%d_gou_%s_rand_%d", uid, time.Now().Unix(), fileType, rand.Intn(100000))
	// Do md5, take the first 20,//mapIntToBase62 map the number to Base62
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s", input)))
	hashString := base64.StdEncoding.EncodeToString(hash[:])
	if len(hashString) > num {
		hashString = hashString[:num]
	}

	result := ""
	for _, char := range hashString {
		index := int(char) % 62
		result += string(baseWord[index])
	}
	return result
}
