package custom

import (
	"net/url"
	"strconv"
)

// В зависимости от типа параметра переводит в нужные типы с обработкой ошибок
func QueryParamsConv(params map[string]string, vals url.Values) (map[string]any, error) {
	resp := make(map[string]any, 1)
	var err error
	for param, varType := range params {
		switch varType {
		case "int":
			resp[param], err = strconv.Atoi(vals.Get(param))
			if err != nil && vals.Get(param) != "" {
				return nil, err
			}
		case "float64":
			resp[param], err = strconv.ParseFloat(vals.Get(param), 64)
			if err != nil && vals.Get(param) != "" {
				return nil, err
			}
		case "bool":
			resp[param], err = strconv.ParseBool(vals.Get(param))
			if err != nil && vals.Get(param) != "" {
				return nil, err
			}
		case "string":
			resp[param] = vals.Get(param)
		}
	}
	return resp, nil
}
