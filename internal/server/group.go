package server

import (
	"erp/internal/controllers"
	"erp/internal/middlewares"
	"net/http"
)

type Group struct {
	mux         *http.ServeMux
	middlewares []middlewares.Middleware
	controllers []controllers.Controller
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

func (group *Group) Add(controllers ...controllers.Controller) {
	group.controllers = append(group.controllers, controllers...)
}

func (group *Group) Group(groupFunc func(g *Group)) {
	ng := NewGroup(group.mux)
	ng.middlewares = make([]middlewares.Middleware, len(group.middlewares))
	copy(ng.middlewares, group.middlewares)
	groupFunc(&ng)
	ng.Chain()
}

func (group *Group) Chain() {
	for _, controller := range group.controllers {
		pattern := controller.Pattern()
		handler := http.Handler(controller)
		for i := len(group.middlewares) - 1; i >= 0; i-- {
			handler = group.middlewares[i](handler)
		}
		group.mux.Handle(pattern, handler)
	}
}
