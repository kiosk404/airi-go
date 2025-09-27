package safego

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

func Recovery(ctx context.Context) {
	e := recover()
	if e == nil {
		return
	}

	if ctx == nil {
		ctx = context.Background() // nolint: byted_context_not_reinitialize -- false positive
	}

	err := fmt.Errorf("%v", e)
	logs.Error(fmt.Sprintf("[catch panic] err = %v \n stacktrace:\n%s", err, debug.Stack()))
}

func Go(ctx context.Context, fn func()) {
	go func() {
		defer Recovery(ctx)

		fn()
	}()
}
