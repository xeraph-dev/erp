package api

import (
	"erp/internal/handlers"
	"erp/internal/middlewares"
	"net/http"
)

type Group struct {
	rootPattern string
	mux         *http.ServeMux
	middlewares []middlewares.Middleware
	handlers    []handlers.Handler
}

var _ Router = (*Group)(nil)

func NewGroup(mux *http.ServeMux) (group *Group) {
	group = new(Group)
	group.mux = mux
	return
}

func NewGroupWithMiddlewares(mux *http.ServeMux, baseMiddlewares []middlewares.Middleware) (group *Group) {
	group = NewGroup(mux)
	group.middlewares = make([]middlewares.Middleware, len(baseMiddlewares))
	copy(group.middlewares, baseMiddlewares)
	return
}

func (*Group) __internal() {}

func (group *Group) Use(middlewares ...middlewares.Middleware) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *Group) Handle(handlers ...handlers.Handler) {
	group.handlers = append(group.handlers, handlers...)
}

func (group *Group) Group(groupFunc func(group *Group)) {
	ng := NewGroupWithMiddlewares(group.mux, group.middlewares)
	groupFunc(ng)
	ng.Chain()
}

func (group *Group) Route(pattern string, groupFunc func(group *Group)) {
	ng := NewGroupWithMiddlewares(group.mux, group.middlewares)
	ng.rootPattern += pattern
	groupFunc(ng)
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
