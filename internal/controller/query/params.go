package query

import (
	"fmt"
	"net/url"
	"timeline/internal/controller/scope"
)

const (
	ErrParamNotProvided = "required param \"%s\" not provided"
)

type Params []Param

func NewParams(settings *scope.Settings, params ...Param) Params {
	supported := settings.SupportedParams
	list := make(Params, 0, len(params))
	for _, spec := range params {
		ptype := supported.GetParam(spec.getName())
		switch {
		case scope.INT == ptype:
			list = append(list, spec)
		case scope.BOOL == ptype:
			list = append(list, spec)
		case scope.STRING == ptype:
			list = append(list, spec)
		case scope.FLOAT32 == ptype:
			list = append(list, spec)
		default:
			continue
		}
	}
	return list
}

func (p Params) Parse(vals url.Values) error {
	if err := validate(p, vals); err != nil {
		return err
	}
	if err := parse(p, vals); err != nil {
		return err
	}
	return nil
}

func validate(p Params, vals url.Values) error {
	for i := range p {
		if p[i].isRequired() {
			if !vals.Has(p[i].getName()) {
				return fmt.Errorf(ErrParamNotProvided, p[i].getName())
			}
		}
	}
	return nil
}

func parse(p Params, vals url.Values) error {
	for _, param := range p {
		if vals.Has(param.getName()) {
			if err := param.load(vals.Get(param.getName())); err != nil {
				return fmt.Errorf("param.Load: %w", err)
			}
		}
	}
	return nil
}
