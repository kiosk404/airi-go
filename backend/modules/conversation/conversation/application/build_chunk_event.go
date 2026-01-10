package application

import (
	"github.com/gin-contrib/sse"
	"github.com/kiosk404/airi-go/backend/api/model/conversation/run"
	"github.com/kiosk404/airi-go/backend/pkg/json"
)

func buildDoneEvent(event string) *sse.Event {
	return &sse.Event{
		Event: event,
		Data:  "[DONE]",
	}
}

func buildErrorEvent(errCode int64, errMsg string) *sse.Event {
	errData := run.ErrorData{
		Code: errCode,
		Msg:  errMsg,
	}
	ed, _ := json.Marshal(errData)

	return &sse.Event{
		Event: run.RunEventError,
		Data:  ed,
	}
}

func buildMessageChunkEvent(event string, chunkMsg []byte) *sse.Event {
	return &sse.Event{
		Event: event,
		Data:  chunkMsg,
	}
}
