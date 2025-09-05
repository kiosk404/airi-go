package goroutineutil

import (
	"context"
	"runtime"

	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

func GoWithDeferFunc(ctx context.Context, f func()) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				logs.Error("goroutine panic: %s: %s", e, buf)
			}
		}()
		f()
	}()
}

func GoWithDefaultRecovery(ctx context.Context, f func()) {
	GoWithDeferFunc(ctx, f)
}
