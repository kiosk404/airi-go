package impl

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/llm/application"
	"github.com/kiosk404/airi-go/backend/modules/llm/pkg"
	"github.com/kiosk404/airi-go/backend/pkg/ctxcache"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

var ctxCacheKey = "builtin_chat_model_in_context"

func GetBuiltinChatModel(ctx context.Context, envPrefix string) (bcm application.BaseChatModel, configured bool, err error) {
	bcm, ok := ctxcache.Get[application.BaseChatModel](ctx, ctxCacheKey)
	if ok {
		logs.DebugX(pkg.ModelName, "builtin chat model in context: %v", bcm)
		return bcm, true, nil
	}

	return
}
