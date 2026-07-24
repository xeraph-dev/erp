package api

import (
	"net/http"
)

type Group struct {
	mux         *http.ServeMux
	middlewares []Middleware
	handlers    []Handler
}

var _ Router = (*Group)(nil)

func NewGroup(mux *http.ServeMux) (group *Group) {
	group = new(Group)
	group.mux = mux
	return
}

func NewGroupWithMiddlewares(mux *http.ServeMux, baseMiddlewares []Middleware) (group *Group) {
	group = NewGroup(mux)
	group.middlewares = make([]Middleware, len(baseMiddlewares))
	copy(group.middlewares, baseMiddlewares)
	return
}

func (*Group) __internal() {}

func (group *Group) Use(middlewares ...Middleware) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *Group) Handle(handlers ...Handler) {
	group.handlers = append(group.handlers, handlers...)
}

func (group *Group) Group(groupFunc func(group *Group)) {
	ng := NewGroupWithMiddlewares(group.mux, group.middlewares)
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
