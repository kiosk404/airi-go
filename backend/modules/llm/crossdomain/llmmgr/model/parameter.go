package model

import (
	"fmt"
)

type ModelParameter struct {
	Name       string                  `json:"name"`
	Label      string                  `json:"label"`
	Desc       string                  `json:"desc"`
	Type       ModelParamType          `json:"type"`
	Min        string                  `json:"min"`
	Max        string                  `json:"max"`
	Precision  int32                   `json:"precision"`
	DefaultVal *ModelParamDefaultValue `json:"default_val"`
	Options    []*Option               `json:"options"`
	ParamClass *ModelParamClass        `json:"param_class"`
}

type ModelParamType int64

const (
	ModelParamType_Float   ModelParamType = 1
	ModelParamType_Int     ModelParamType = 2
	ModelParamType_Boolean ModelParamType = 3
	ModelParamType_String  ModelParamType = 4
)

func (p ModelParamType) String() string {
	switch p {
	case ModelParamType_Float:
		return "Float"
	case ModelParamType_Int:
		return "Int"
	case ModelParamType_Boolean:
		return "Boolean"
	case ModelParamType_String:
		return "String"
	}
	return "<UNSET>"
}

func ModelParamTypeFromString(s string) (ModelParamType, error) {
	switch s {
	case "Float":
		return ModelParamType_Float, nil
	case "Int":
		return ModelParamType_Int, nil
	case "Boolean":
		return ModelParamType_Boolean, nil
	case "String":
		return ModelParamType_String, nil
	}
	return ModelParamType(0), fmt.Errorf("not a valid ModelParamType string")
}

type ModelParamDefaultValue struct {
	DefaultVal string `json:"default_val"`
	Creative   string `json:"creative,omitempty"`
	Balance    string `json:"balance,omitempty"`
	Precise    string `json:"precise,omitempty"`
}

func NewModelParamDefaultValue() *ModelParamDefaultValue {
	return &ModelParamDefaultValue{}
}

func (p *ModelParamDefaultValue) InitDefault() {
}

func (p *ModelParamDefaultValue) GetDefaultVal() (v string) {
	return p.DefaultVal
}

type Option struct {
	Label string `thrift:"label,1" json:"label"`
	Value string `thrift:"value,2" json:"value"`
}

func NewOption() *Option {
	return &Option{}
}

func (p *Option) InitDefault() {
}

func (p *Option) GetLabel() (v string) {
	return p.Label
}

func (p *Option) GetValue() (v string) {
	return p.Value
}
func (p *Option) SetLabel(val string) {
	p.Label = val
}
func (p *Option) SetValue(val string) {
	p.Value = val
}

func (p *Option) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Option(%+v)", *p)
}

type ModelParamClass struct {
	ClassID int32  `thrift:"class_id,1" json:"class_id"`
	Label   string `thrift:"label,2" json:"label"`
}

func NewModelParamClass() *ModelParamClass {
	return &ModelParamClass{}
}

func (p *ModelParamClass) InitDefault() {
}

func (p *ModelParamClass) GetClassID() (v int32) {
	return p.ClassID
}

func (p *ModelParamClass) GetLabel() (v string) {
	return p.Label
}
func (p *ModelParamClass) SetClassID(val int32) {
	p.ClassID = val
}
func (p *ModelParamClass) SetLabel(val string) {
	p.Label = val
}

func (p *ModelParamClass) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ModelParamClass(%+v)", *p)
}
