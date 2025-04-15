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
	SupportedMethodsMap map[string]uint8
	SupportedParams     SupportedParams
	EnableAuthorization bool
	EnableRepoS3        bool
	EnableRepoMail      bool
}

func NewDefaultSettings(appCfg config.Application) *Settings {
	return &Settings{
		EnableAuthorization: appCfg.Settings.EnableAuthorization,
		EnableRepoS3:        appCfg.Settings.EnableRepoS3,
		EnableRepoMail:      appCfg.Settings.EnableRepoMail,
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
