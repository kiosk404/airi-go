package modelmgr

import (
	"fmt"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service/llmimpl/chatmodel"
	"strconv"
)

type Model struct {
	ID                int64        `yaml:"id"`
	Name              string       `yaml:"name"`
	IconURI           string       `yaml:"icon_uri"`
	IconURL           string       `yaml:"icon_url"`
	DefaultParameters []*Parameter `yaml:"default_parameters"`
	Meta              ModelMeta    `yaml:"meta"`
}

func (m *Model) FindParameter(name ParameterName) (*Parameter, bool) {
	if len(m.DefaultParameters) == 0 {
		return nil, false
	}

	for _, param := range m.DefaultParameters {
		if param.Name == name {
			return param, true
		}
	}

	return nil, false
}

type Parameter struct {
	Name       ParameterName  `json:"name" yaml:"name"`
	Label      string         `json:"label,omitempty" yaml:"label,omitempty"`
	Desc       string         `json:"desc" yaml:"desc"`
	Type       ValueType      `json:"type" yaml:"type"`
	Min        string         `json:"min" yaml:"min"`
	Max        string         `json:"max" yaml:"max"`
	DefaultVal DefaultValue   `json:"default_val" yaml:"default_val"`
	Precision  int            `json:"precision,omitempty" yaml:"precision,omitempty"` // float precision, default 2
	Options    []*ParamOption `json:"options" yaml:"options"`                         // enum options
	Style      DisplayStyle   `json:"param_class" yaml:"style"`
}

func (p *Parameter) GetFloat(tp DefaultType) (float64, error) {
	if p.Type != ValueTypeFloat {
		return 0, fmt.Errorf("unexpected paramerter type, name=%v, expect=%v, given=%v",
			p.Name, ValueTypeFloat, p.Type)
	}

	if tp != DefaultTypeDefault && p.DefaultVal[tp] == "" {
		tp = DefaultTypeDefault
	}

	val, ok := p.DefaultVal[tp]
	if !ok {
		return 0, fmt.Errorf("unexpected default type, name=%v, type=%v", p.Name, tp)
	}

	return strconv.ParseFloat(val, 64)
}

func (p *Parameter) GetInt(tp DefaultType) (int64, error) {
	if p.Type != ValueTypeInt {
		return 0, fmt.Errorf("unexpected paramerter type, name=%v, expect=%v, given=%v",
			p.Name, ValueTypeInt, p.Type)
	}

	if tp != DefaultTypeDefault && p.DefaultVal[tp] == "" {
		tp = DefaultTypeDefault
	}
	val, ok := p.DefaultVal[tp]
	if !ok {
		return 0, fmt.Errorf("unexpected default type, name=%v, type=%v", p.Name, tp)
	}
	return strconv.ParseInt(val, 10, 64)
}

func (p *Parameter) GetBool(tp DefaultType) (bool, error) {
	if p.Type != ValueTypeBoolean {
		return false, fmt.Errorf("unexpected paramerter type, name=%v, expect=%v, given=%v",
			p.Name, ValueTypeBoolean, p.Type)
	}
	if tp != DefaultTypeDefault && p.DefaultVal[tp] == "" {
		tp = DefaultTypeDefault
	}
	val, ok := p.DefaultVal[tp]
	if !ok {
		return false, fmt.Errorf("unexpected default type, name=%v, type=%v", p.Name, tp)
	}
	return strconv.ParseBool(val)
}

func (p *Parameter) GetString(tp DefaultType) (string, error) {
	if tp != DefaultTypeDefault && p.DefaultVal[tp] == "" {
		tp = DefaultTypeDefault
	}

	val, ok := p.DefaultVal[tp]
	if !ok {
		return "", fmt.Errorf("unexpected default type, name=%v, type=%v", p.Name, tp)
	}
	return val, nil
}

type ModelMeta struct {
	Protocol   entity.Protocol   `yaml:"protocol"`    // Model Communication Protocol
	Capability *Capability       `yaml:"capability"`  // model capability
	ConnConfig *chatmodel.Config `yaml:"conn_config"` // model connection configuration
	Status     ModelStatus       `yaml:"status"`      // model state
}

type DefaultValue map[DefaultType]string

type DisplayStyle struct {
	Widget Widget `json:"class_id" yaml:"widget"`
	Label  string `json:"label" yaml:"label"`
}

type ParamOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Capability struct {
	// Model supports function calling
	FunctionCall bool `json:"function_call" yaml:"function_call" mapstructure:"function_call"`
	// Input modals
	InputModal []Modal `json:"input_modal,omitempty" yaml:"input_modal,omitempty" mapstructure:"input_modal,omitempty"`
	// Input tokens
	InputTokens int `json:"input_tokens" yaml:"input_tokens" mapstructure:"input_tokens"`
	// Model supports json mode
	JSONMode bool `json:"json_mode" yaml:"json_mode" mapstructure:"json_mode"`
	// Max tokens
	MaxTokens int `json:"max_tokens" yaml:"max_tokens" mapstructure:"max_tokens"`
	// Output modals
	OutputModal []Modal `json:"output_modal,omitempty" yaml:"output_modal,omitempty" mapstructure:"output_modal,omitempty"`
	// Output tokens
	OutputTokens int `json:"output_tokens" yaml:"output_tokens" mapstructure:"output_tokens"`
	// Model supports prefix caching
	PrefixCaching bool `json:"prefix_caching" yaml:"prefix_caching" mapstructure:"prefix_caching"`
	// Model supports reasoning
	Reasoning bool `json:"reasoning" yaml:"reasoning" mapstructure:"reasoning"`
	// Model supports prefill response
	PrefillResponse bool `json:"prefill_response" yaml:"prefill_response" mapstructure:"prefill_response"`
}
