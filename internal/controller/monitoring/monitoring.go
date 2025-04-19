package monitoring

import (
	"fmt"
	"net/http"
	"strings"
	"timeline/internal/controller/scope"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ServerMonitoring struct {
	Logger   *zap.Logger
	settings *scope.Settings
	Router   *mux.Router
}

func New(logger *zap.Logger, settings *scope.Settings) *ServerMonitoring {
	return &ServerMonitoring{
		Logger:   logger,
		settings: settings,
	}
}

// @Summary Health
// @Description docker health check
// @Tags server monitoring
// @Success 200
// @Failure 500
// @Router /health [get]
func (a *ServerMonitoring) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// @Summary Get registered routes
// @Tags server monitoring
// @Success 200 {object} string "list of routes"
// @Failure 500
// @Router /routes [get]
func (a *ServerMonitoring) GetRoutes(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	routesjson, err := getAllRoutes(a.Router)
	if err != nil {
		a.Logger.Error("getAllRoutes", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte(routesjson))
	}
}

/*
	{
		"paths": [
			{
				"path": /v1/users
				"methods": [
					"POST",
					"PUT",
					"GET"
				]
			},
			{
				"path": /v1/users/orgmap
				"methods": [
					"GET"
				]
			}
		]
	}
*/
func getAllRoutes(router *mux.Router) (string, error) {
	var body strings.Builder
	body.WriteString("{\"paths\":[")
	prevPath := ""
	pathCnt := 0
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil { // встретили subrouter - сброс prevPath
			prevPath = ""
			return nil
		}
		if prevPath == pathTemplate { // копим методы
			for i := range methods {
				body.WriteString(fmt.Sprintf(",\"%s\"", methods[i]))
			}
		} else { // добавляем новый элемент path в массив paths
			prevPath = pathTemplate
			if pathCnt > 0 {
				body.WriteString("]},")
			}
			body.WriteString(fmt.Sprintf("{\"path\":\"%s\"", pathTemplate))
			if len(methods) != 0 {
				body.WriteString(",")
			}
			body.WriteString("\"methods\":[")
			for i := range methods {
				if i > 0 {
					body.WriteString(",")
				}
				body.WriteString(fmt.Sprintf("\"%s\"", methods[i]))
			}
			pathCnt++
		}
		return nil
	})
	if err != nil {
		return "", err
	} else {
		body.WriteString("]}]}")
	}
	return body.String(), nil
}
