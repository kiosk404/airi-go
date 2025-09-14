package llmfactory

import (
	"context"

	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llmimpl"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llminterface"
)

//go:generate mockgen -destination=mocks/factory.go -package=mocks . IFactory
type IFactory interface {
	CreateLLM(ctx context.Context, model *entity.Model, opts ...entity.Option) (llminterface.ILLM, error)
}

type FactoryImpl struct{}

var _ IFactory = (*FactoryImpl)(nil)
var ModelF *FactoryImpl

func (f *FactoryImpl) CreateLLM(ctx context.Context, model *entity.Model, opts ...entity.Option) (llminterface.ILLM, error) {
	// 用 factory 创建llm接口的实现
	return llmimpl.NewLLM(ctx, model, opts...)
}
