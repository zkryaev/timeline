package query

import (
	"fmt"
	"strconv"
)

type IntParam struct {
	Name string
	// Type     string
	Val      int
	Required bool
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
	return nil
}

func (p IntParam) isRequired() bool {
	return p.Required
}

func (p IntParam) getName() string {
	return p.Name
}

// func (p IntParam) getType() string {
// 	return p.Type
// }

// func (p IntParam) getInt() int {
// 	return p.Val
// }

// func (p IntParam) getBool() bool {
// 	return p.Val != 0
// }

// func (p IntParam) getString() string {
// 	return strconv.Itoa(p.Val)
// }

// func (p IntParam) getFloat32() float32 {
// 	return float32(p.Val)
// }

type BoolParam struct {
	Name string
	//Type     string
	Val      bool
	Required bool
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
	return nil
}

func (p BoolParam) isRequired() bool {
	return p.Required
}

func (p BoolParam) getName() string {
	return p.Name
}

// func (p BoolParam) getType() string {
// 	return p.Type
// }

// func (p BoolParam) getInt() int {
// 	if p.Val {
// 		return 1
// 	} else {
// 		return 0
// 	}
// }

// func (p BoolParam) getBool() bool {
// 	return p.Val
// }

// func (p BoolParam) getString() string {
// 	return strconv.FormatBool(p.Val)
// }

// func (p BoolParam) getFloat32() float32 {
// 	return float32(p.getInt())
// }

type StringParam struct {
	Name string
	//Type     string
	Val      string
	Required bool
}

func NewParamString(name string, required bool) *StringParam {
	return &StringParam{Name: name, Required: required}
}

func (p *StringParam) load(val string) error {
	p.Val = val
	return nil
}

func (p StringParam) getName() string {
	return p.Name
}

func (p StringParam) isRequired() bool {
	return p.Required
}

// func (p StringParam) getType() string {
// 	return p.Type
// }
// func (p StringParam) getInt() int {
// 	return 0
// }
// func (p StringParam) getBool() bool {
// 	if strings.Contains(p.Val, "true") {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func (p StringParam) getString() string {
// 	return p.Val
// }
// func (p StringParam) getFloat32() float32 {
// 	return 0
// }

type Float32Param struct {
	Name string
	//Type     string
	Val      float32
	Required bool
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
	return nil
}

func (p Float32Param) getName() string {
	return p.Name
}

func (p Float32Param) isRequired() bool {
	return p.Required
}

// func (p Float32Param) getType() string {
// 	return p.Type
// }
// func (p Float32Param) getInt() int {
// 	return int(p.Val)
// }

// func (p Float32Param) getBool() bool {
// 	return p.Val != 0
// }

// func (p Float32Param) getString() string {
// 	return strconv.FormatFloat(float64(p.Val), 'f', 4, 32)
// }

// func (p Float32Param) getFloat32() float32 {
// 	return p.Val
// }
