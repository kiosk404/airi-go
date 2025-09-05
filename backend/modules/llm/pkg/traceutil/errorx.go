package traceutil

import (
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
)

func GetTraceStatusCode(err error) int32 {
	if statusErr, ok := errorx.FromStatusError(err); ok {
		return statusErr.Code()
	}
	return -1
}
