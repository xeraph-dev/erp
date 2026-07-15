package controllers

import "net/http"

type Controller interface {
	http.Handler
	__internal()
	Pattern() string
}
