package server

import (
	"erp/internal/controllers"
	"erp/internal/middlewares"
	"erp/internal/repositories"
	"erp/internal/services"
	"fmt"
	"net/http"
	"reflect"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	mux          *http.ServeMux
	middlewares  []middlewares.Middleware
	pool         *pgxpool.Pool
	repositories []repositories.Repository
	services     []services.Service
}

var _ Router = &Server{}

func NewServer() (server Server) {
	server.mux = http.NewServeMux()
	return
}

func (server *Server) Use(middlewares ...middlewares.Middleware) {
	server.middlewares = append(server.middlewares, middlewares...)
}

func (server *Server) Add(controllers ...controllers.Controller) {
	for _, controller := range controllers {
		v := reflect.ValueOf(controller).Elem()

		for st, sv := range v.FieldByName("Services").Fields() {
			for _, service := range server.services {
				if reflect.TypeOf(service).Implements(st.Type) {
					sv.Set(reflect.ValueOf(service))
				}
			}
		}

		server.mux.Handle(controller.Pattern(), controller)
	}
}

func (server *Server) Database(pool *pgxpool.Pool) {
	server.pool = pool
}

func (server *Server) RepoRegister(repositories ...repositories.Repository) {
	server.repositories = append(server.repositories, repositories...)
}

func (server *Server) ServiceRegister(services ...services.Service) {
	for _, service := range services {
		v := reflect.ValueOf(service).Elem()

		dbv := v.FieldByName("DB")
		dbv.Set(reflect.ValueOf(server.pool))

		if dbv.IsNil() {
			panic("missing an implementation of a database driver")
		}

		for rt, rv := range v.FieldByName("Repos").Fields() {
			for _, repo := range server.repositories {
				if reflect.TypeOf(repo).Implements(rt.Type) {
					rv.Set(reflect.ValueOf(repo))
				}
			}
			if rv.IsNil() {
				panic(fmt.Sprintf("missing an implementation of %s", rt.Type))
			}
		}

		server.services = append(server.services, service)
	}
}

func (server *Server) Group(groupFunc func(g *Group)) {
	group := NewGroup(server.mux, server.services)
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
