package server

import (
	"erp/internal/controllers"
	"erp/internal/middlewares"
)

type Router interface {
	__internal()
	Use(middlewares ...middlewares.Middleware)
	Add(controller ...controllers.Controller)
}
