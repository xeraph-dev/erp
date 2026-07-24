package handlers

import (
	"erp/internal/shared/api"
	"net/http"
)

type HealthHandler struct{}

var _ api.Handler = (*HealthHandler)(nil)

func NewHealthHandler() api.Handler {
	return HealthHandler{}
}

func (HealthHandler) Pattern() string { return "/health" }
func (handler HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
