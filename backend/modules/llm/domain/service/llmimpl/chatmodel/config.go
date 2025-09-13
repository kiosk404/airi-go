package chatmodel

import (
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/json"
	"github.com/pkg/errors"
)

type Config struct {
	entity.CommonParam
	entity.ProtocolConfig

	Custom map[string]string `json:"custom,omitempty" yaml:"custom"`
}

func NewConfig(model *entity.Model, opts ...entity.Option) (*Config, error) {
	if err := checkModelBeforeBuild(model); err != nil {
		return nil, err
	}
	commonParam := model.ParamConfig.GetCommonParamDefaultVal()
	entity.ApplyOptions(entity.NewDefaultParams(commonParam), opts...)

	return &Config{
		CommonParam:    commonParam,
		ProtocolConfig: model.ProtocolConfig,
		Custom:         make(map[string]string),
	}, nil
}

func checkModelBeforeBuild(model *entity.Model) error {
	if model == nil {
		return errors.Errorf("[checkModelBeforeBuild] failed as model:%s", json.MarshalStringIgnoreErr(model))
	}
	return nil
}
