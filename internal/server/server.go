package server

import (
	"erp/internal/middlewares"
	"net/http"
)

type Group struct {
	mux         *http.ServeMux
	middlewares []middlewares.Middleware
	handlers    map[string]http.Handler
	groups      []Group
}

func NewGroup(mux *http.ServeMux) (group Group) {
	group.mux = mux
	group.handlers = make(map[string]http.Handler)
	return
}

func (g *Group) Use(middlewares ...middlewares.Middleware) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) HandleFunc(pattern string, handler http.HandlerFunc) {
	g.handlers[pattern] = handler
}

func (g *Group) Group(group func(g *Group)) {
	ng := NewGroup(g.mux)
	ng.middlewares = make([]middlewares.Middleware, len(g.middlewares))
	copy(ng.middlewares, g.middlewares)
	group(&ng)
	ng.Chain()
}

func (g *Group) Chain() {
	for pattern, handler := range g.handlers {
		for i := len(g.middlewares) - 1; i >= 0; i-- {
			handler = g.middlewares[i](handler)
		}
		g.mux.Handle(pattern, handler)
	}
}

type Server struct {
	mux         *http.ServeMux
	middlewares []middlewares.Middleware
}

func NewServer() (server Server) {
	server.mux = http.NewServeMux()
	return
}

func (s *Server) Use(middlewares ...middlewares.Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) HandleFunc(pattern string, handler http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *Server) Group(group func(g *Group)) {
	g := NewGroup(s.mux)
	group(&g)
	g.Chain()
}

func (s *Server) Serve(addr string) error {
	return http.ListenAndServe(addr, s.chain())
}

func (s *Server) chain() http.Handler {
	handler := http.Handler(s.mux)
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		handler = s.middlewares[i](handler)
	}
	return handler
}
