package scope

import (
	"net/http"
	"timeline/internal/config"
)

// supported methods
const (
	POST   = http.MethodPost
	GET    = http.MethodGet
	PUT    = http.MethodPut
	DELETE = http.MethodDelete

	ALL = "ALL"
)

// supported query param types
const (
	INT     = "int"
	BOOL    = "bool"
	STRING  = "string"
	FLOAT32 = "float32"
)

// supported entity types
const (
	GALLERY = "gallery"
	BANNER  = "banner"
	ORG     = "org"
	USER    = "user"
	WORKER  = "worker"
)

type Settings struct {
	SupportedMethodsMap map[string]struct{}
	SupportedParams     SupportedParams
	EnableAuthorization bool
	EnableMedia         bool
	EnableMail          bool
	EnableMetrics       bool
}

func NewDefaultSettings(appCfg config.Application) *Settings {
	return &Settings{
		EnableAuthorization: appCfg.Settings.EnableAuthorization,
		EnableMedia:         appCfg.Settings.EnableMedia,
		EnableMail:          appCfg.Settings.EnableMail,
		EnableMetrics:       appCfg.Settings.EnableMetrics,
		SupportedMethodsMap: defaultSupportedMethodsHTTP(),
		SupportedParams:     defaultSupportedParams(),
	}
}

func defaultSupportedMethodsHTTP() map[string]struct{} {
	return map[string]struct{}{
		http.MethodGet:    {},
		http.MethodPost:   {},
		http.MethodPut:    {},
		http.MethodDelete: {},
	}
}
