package server

import (
	"erp/internal/controllers"
	"erp/internal/middlewares"
	"erp/internal/services"
	"net/http"
	"reflect"
)

type Group struct {
	mux         *http.ServeMux
	middlewares []middlewares.Middleware
	controllers []controllers.Controller
	services    []services.Service
}

var _ Router = &Group{}

func NewGroup(mux *http.ServeMux, services []services.Service) (group Group) {
	group.mux = mux
	group.services = services
	return
}

func (group *Group) Use(middlewares ...middlewares.Middleware) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *Group) Add(controllers ...controllers.Controller) {
	for _, controller := range controllers {
		v := reflect.ValueOf(controller).Elem()

		for st, sv := range v.FieldByName("Services").Fields() {
			for _, service := range group.services {
				if reflect.TypeOf(service).Implements(st.Type) {
					sv.Set(reflect.ValueOf(service))
				}
			}
		}

		group.controllers = append(group.controllers, controller)
	}
}

func (group *Group) Group(groupFunc func(g *Group)) {
	ng := NewGroup(group.mux, group.services)
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
