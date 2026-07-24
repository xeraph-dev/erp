package api

import "net/http"

type Handler interface {
	http.Handler
	Pattern() string
}

type Middleware func(next http.Handler) http.Handler

type Router interface {
	__internal()
	Use(middlewares ...Middleware)
	Handle(handlers ...Handler)
}
