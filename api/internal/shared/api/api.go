package api

import "net/http"

type Middleware func(next http.Handler) http.Handler

type Router interface {
	__internal()
	Use(middlewares ...Middleware)
	HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...Middleware)
}
