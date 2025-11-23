package llminterface

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
)

//go:generate mockgen -destination=mocks/llm.go -package=mocks . ILLM
type ILLM interface {
	// Generate 非流式
	Generate(ctx context.Context, input []*entity.Message, opts ...model.Option) (*entity.Message, error)
	// Stream 流式
	Stream(ctx context.Context, input []*entity.Message, opts ...model.Option) (
		entity.IStreamReader, error)
}
