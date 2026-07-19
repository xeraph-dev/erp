package server

import (
	"erp/internal/handlers"
	"erp/internal/middlewares"
)

type Router interface {
	__internal()
	Use(middlewares ...middlewares.Middleware)
	Handle(handlers ...handlers.Handler)
}
