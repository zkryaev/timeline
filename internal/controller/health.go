package controller

import "net/http"

// @Summary Health
// @Description Server health check
// @Tags Server
// @Success 200
// @Failure 500
// @Router /health [get]
func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
