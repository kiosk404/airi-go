package config

import (
	"context"

	llm_conf "github.com/kiosk404/airi-go/backend/modules/llm/domain/component/conf"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/conf"
)

type RuntimeImpl struct {
	cfg *entity.RuntimeConfig
}

func NewRuntime(ctx context.Context, factory conf.IConfigLoaderFactory) (llm_conf.IConfigRuntime, error) {
	loader, err := factory.NewConfigLoader("model_runtime_config.yaml")
	if err != nil {
		return nil, err
	}
	var cfg entity.RuntimeConfig
	if err = loader.Unmarshal(ctx, &cfg); err != nil {
		return nil, err
	}
	return &RuntimeImpl{
		cfg: &cfg,
	}, nil
}

func (r *RuntimeImpl) NeedCvtURLToBase64() bool {
	if r == nil || r.cfg == nil {
		return false
	}
	return r.cfg.NeedCvtURLToBase64
}
