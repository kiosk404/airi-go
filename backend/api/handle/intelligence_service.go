package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/api/model/app/intelligence"
	searchapp "github.com/kiosk404/airi-go/backend/modules/data/search/application"
)

// GetDraftIntelligenceList .
// @router /api/intelligence_api/search/get_draft_intelligence_list [POST]
func GetDraftIntelligenceList(c *gin.Context) {
	var err error
	ctx := c.Request.Context()
	var req intelligence.GetDraftIntelligenceListRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := searchapp.SearchSVC.GetDraftIntelligenceList(ctx, &req)
	if err != nil {
		internalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
