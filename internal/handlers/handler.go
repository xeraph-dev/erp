package handlers

import (
	"erp/internal/codecs"
	"erp/internal/dtos"
	"errors"
	"io"
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

func extractRefreshToken(r *http.Request, codec codecs.Codec) (token string, ok bool) {
	var dto dtos.RefreshToken
	if err := codec.Decode(r.Body, &dto); (err == nil || errors.Is(err, io.EOF)) && dto.RefreshToken != "" {
		return dto.RefreshToken, true
	}
	if cookie, err := r.Cookie("refresh_token"); err == nil {
		return cookie.Value, true
	}

	return "", false
}

func setAuthCookies(w http.ResponseWriter, pair dtos.TokenPair) {
	http.SetCookie(w, cookie("access_token", pair.AccessToken, pair.AccessTokenExpiresAt))
	http.SetCookie(w, cookie("refresh_token", pair.RefreshToken, pair.RefreshTokenExpiresAt))
}
