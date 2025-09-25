package handle

import (
	"github.com/gin-gonic/gin"
	httputil "github.com/kiosk404/airi-go/backend/api/http"
)

func invalidParamRequestResponse(c *gin.Context, errMsg string) {
	httputil.BadRequest(c, errMsg)
}

func internalServerErrorResponse(c *gin.Context, err error) {
	httputil.InternalError(c, err)
}
