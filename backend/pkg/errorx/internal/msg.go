package internal

import (
	"errors"
	"fmt"
)

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Unwrap() error {
	return w.cause
}

func (w *withMessage) Error() string {
	return fmt.Sprintf("%s\ncause=%s", w.msg, w.cause.Error())
}

func wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}

	return err
}

func Wrapf(err error, format string, args ...interface{}) error {
	return withStackTraceIfNotExists(wrapf(err, format, args...))
}

// GetStatusError 获取错误链中最顶层的 StatusError.
// 如果有获取code或其他扩展字段的需求，再考虑对外暴露
func GetStatusError(err error) *statusError {
	if err == nil {
		return nil
	}

	var ws *statusError
	if errors.As(err, &ws) {
		return ws
	}

	return nil
}

// FromStatusError converts err to StatusError.
// 解析RPC返回的error, 如果是statusError转换而来, 则返回ok为true
func FromStatusError(err error) (statusErr *statusError, ok bool) {
	if err == nil {
		return nil, false
	}

	if se := GetStatusError(err); se != nil {
		return se, true
	}

	statusErr = &statusError{
		statusCode: statusErr.Code(),
		message:    statusErr.message,
		ext: Extension{
			IsAffectStability: false,
		},
	}

	return statusErr, true
}
