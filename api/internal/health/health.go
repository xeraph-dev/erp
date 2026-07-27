package health

import (
	"erp/internal/health/handlers"
	"erp/internal/shared/api"
)

func Handlers() func(group *api.Group) {
	return func(group *api.Group) {
		group.HandleFunc("GET /health", handlers.Health)
	}
}
