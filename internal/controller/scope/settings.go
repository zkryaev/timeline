package scope

import (
	"net/http"
	"timeline/internal/config"
)

const (
	POST = iota
	GET
	PUT
	DELETE

	ALL
)

const (
	INT     = "int"
	BOOL    = "bool"
	STRING  = "string"
	FLOAT32 = "float32"
)

type Settings struct {
	SupportedMethodsMap map[string]uint8
	SupportedParams     SupportedParams
	EnableAuthorization bool
}

func NewDefaultSettings(appcfg config.Application) *Settings {
	return &Settings{
		EnableAuthorization: appcfg.EnableAuthorization,
		SupportedMethodsMap: defaultSupportedMethodsHTTP(),
		SupportedParams:     defaultSupportedParams(),
	}
}

func defaultSupportedMethodsHTTP() map[string]uint8 {
	return map[string]uint8{
		http.MethodGet:    GET,
		http.MethodPost:   POST,
		http.MethodPut:    PUT,
		http.MethodDelete: DELETE,
	}
}
