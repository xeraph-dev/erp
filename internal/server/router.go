package server

import (
	"erp/internal/controllers"
	"erp/internal/middlewares"
)

type Router interface {
	Use(middlewares ...middlewares.Middleware)
	Add(controller ...controllers.Controller)
}
