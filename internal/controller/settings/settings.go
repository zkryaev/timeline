package settings

import "net/http"

const (
	POST = iota
	GET
	PUT
	DELETE
	
	ALL
)

type Settings struct {
	SupportedMethodsMap map[string]uint8
}

func NewDefaultSettings() *Settings {
	return &Settings{
		SupportedMethodsMap: defaultSupportedMethodsHTTP(),
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
