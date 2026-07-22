package handlers

import (
	"net/http"
	"time"
)

type Handler interface {
	http.Handler
	__internal()
	Pattern() string
}

func cookie(name string, value string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expiresAt,
		Path:     "/api/auth",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}
}
