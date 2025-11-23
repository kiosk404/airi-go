package chatmodel

import (
	"time"

	"github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/json"
	"github.com/pkg/errors"
)

type Config struct {
	model.CommonParam
	model.ProtocolConfig
	Timeout time.Duration     `json:"timeout,omitempty" yaml:"timeout"`
	Custom  map[string]string `json:"custom,omitempty" yaml:"custom"`
}

func NewConfig(m *entity.Model, opts ...model.Option) (*Config, error) {
	if err := checkModelBeforeBuild(m); err != nil {
		return nil, err
	}
	commonParam := m.ParamConfig.GetCommonParamDefaultVal()
	model.ApplyOptions(model.NewDefaultParams(commonParam), opts...)

	return &Config{
		CommonParam:    commonParam,
		ProtocolConfig: m.GetProtocolConfig(),
		Custom:         make(map[string]string),
	}, nil
}

func checkModelBeforeBuild(model *entity.Model) error {
	if model == nil {
		return errors.Errorf("[checkModelBeforeBuild] failed as model:%s", json.MarshalStringIgnoreErr(model))
	}
	return nil
}
