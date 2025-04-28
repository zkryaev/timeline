package query

import (
	"fmt"
	"strconv"
)

type IntParam struct {
	Name     string
	Val      int
	Required bool
	loaded   bool
}

func NewParamInt(name string, required bool) *IntParam {
	return &IntParam{Name: name, Required: required}
}

func (p *IntParam) load(val string) error {
	valInt, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("strconv.Atoi: %w", err)
	}
	p.Val = valInt
	p.loaded = true
	return nil
}

func (p IntParam) isRequired() bool {
	return p.Required
}

func (p IntParam) getName() string {
	return p.Name
}

func (p *IntParam) EmptyValue() bool {
	return !p.loaded
}

type BoolParam struct {
	Name     string
	Val      bool
	Required bool
	loaded   bool
}

func NewParamBool(name string, required bool) *BoolParam {
	return &BoolParam{Name: name, Required: required}
}

func (p *BoolParam) load(val string) error {
	valBool, err := strconv.ParseBool(val)
	if err != nil {
		return fmt.Errorf("strconv.Atoi: %w", err)
	}
	p.Val = valBool
	p.loaded = true
	return nil
}

func (p BoolParam) isRequired() bool {
	return p.Required
}

func (p BoolParam) getName() string {
	return p.Name
}

func (p BoolParam) EmptyValue() bool {
	return !p.loaded
}

type StringParam struct {
	Name     string
	Val      string
	Required bool
	loaded   bool
}

func NewParamString(name string, required bool) *StringParam {
	return &StringParam{Name: name, Required: required}
}

func (p *StringParam) load(val string) error {
	p.Val = val
	p.loaded = true
	return nil
}

func (p StringParam) getName() string {
	return p.Name
}

func (p StringParam) isRequired() bool {
	return p.Required
}

func (p StringParam) EmptyValue() bool {
	return !p.loaded
}

type Float32Param struct {
	Name     string
	Val      float32
	Required bool
	loaded   bool
}

func NewParamFloat32(name string, required bool) *Float32Param {
	return &Float32Param{Name: name, Required: required}
}

func (p *Float32Param) load(val string) error {
	valFloat, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return fmt.Errorf("strconv.Atoi: %w", err)
	}
	p.Val = float32(valFloat)
	p.loaded = true
	return nil
}

func (p Float32Param) getName() string {
	return p.Name
}

func (p Float32Param) isRequired() bool {
	return p.Required
}

func (p Float32Param) EmptyValue() bool {
	return !p.loaded
}
