package handlers

import (
	"erp/internal/shared/api"
	"net/http"
)

type Health struct{}

var _ api.Handler = (*Health)(nil)

func NewHealthHandler() api.Handler {
	return Health{}
}

func (Health) Pattern() string { return "/health" }
func (handler Health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
