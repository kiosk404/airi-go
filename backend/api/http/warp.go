package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

type data struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

func BadRequest(c *gin.Context, errMsg string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, data{Code: http.StatusBadRequest, Msg: errMsg})
}

func InternalError(c *gin.Context, err error) {
	var customErr errorx.StatusError

	if errors.As(err, &customErr) && customErr.Code() != 0 {
		logs.Warn("[ErrorX] error:  %v %v \n", customErr.Code(), err)
		c.AbortWithStatusJSON(http.StatusOK, data{Code: customErr.Code(), Msg: customErr.Msg()})
		return
	}

	logs.Error("[InternalError]  error: %v \n", err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, data{Code: 500, Msg: "internal server error"})
}
