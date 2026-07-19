package handlers

import (
	"net/http"
)

type Handler interface {
	http.Handler
	__internal()
	Pattern() string
}
