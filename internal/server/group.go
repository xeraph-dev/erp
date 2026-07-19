package server

import (
	"erp/internal/handlers"
	"erp/internal/middlewares"
	"net/http"
)

type Group struct {
	mux         *http.ServeMux
	middlewares []middlewares.Middleware
	handlers    []handlers.Handler
}

var _ Router = (*Group)(nil)

func NewGroup(mux *http.ServeMux) (group Group) {
	group.mux = mux
	return
}

func (*Group) __internal() {}

func (group *Group) Use(middlewares ...middlewares.Middleware) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *Group) Handle(handlers ...handlers.Handler) {
	group.handlers = append(group.handlers, handlers...)
}

func (group *Group) Group(groupFunc func(g *Group)) {
	ng := NewGroup(group.mux)
	ng.middlewares = make([]middlewares.Middleware, len(group.middlewares))
	copy(ng.middlewares, group.middlewares)
	groupFunc(&ng)
	ng.Chain()
}

func (group *Group) Chain() {
	for _, handler := range group.handlers {
		pattern := handler.Pattern()
		chainHandler := http.Handler(handler)
		for i := len(group.middlewares) - 1; i >= 0; i-- {
			chainHandler = group.middlewares[i](chainHandler)
		}
		group.mux.Handle(pattern, chainHandler)
	}
}
