package server

import (
	"erp/internal/controllers"
	"erp/internal/middlewares"
	"net/http"
)

type Server struct {
	mux         *http.ServeMux
	middlewares []middlewares.Middleware
}

var _ Router = (*Server)(nil)

func NewServer() (server Server) {
	server.mux = http.NewServeMux()
	return
}

func (Server) __internal() {}

func (server *Server) Use(middlewares ...middlewares.Middleware) {
	server.middlewares = append(server.middlewares, middlewares...)
}

func (server *Server) Add(controllers ...controllers.Controller) {
	for _, controller := range controllers {
		server.mux.Handle(controller.Pattern(), controller)
	}
}

func (server *Server) Group(groupFunc func(g *Group)) {
	group := NewGroup(server.mux)
	groupFunc(&group)
	group.Chain()
}
func (server *Server) Serve(addr string) error {
	return http.ListenAndServe(addr, server.chain())
}

func (server *Server) chain() http.Handler {
	handler := http.Handler(server.mux)
	for i := len(server.middlewares) - 1; i >= 0; i-- {
		handler = server.middlewares[i](handler)
	}
	return handler
}
