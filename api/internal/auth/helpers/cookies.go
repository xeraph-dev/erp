package helpers

import (
	"erp/internal/shared/codecs"
	"erp/internal/shared/tokens"
	"errors"
	"io"
	"net/http"
	"time"
)

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

func ExtractRefreshToken(r *http.Request, codec codecs.Codec) (token string, ok bool) {
	if err := codec.Decode(r.Body, &token); (err == nil || errors.Is(err, io.EOF)) && token != "" {
		return token, true
	}
	if cookie, err := r.Cookie("refresh_token"); err == nil && cookie.Value != "" {
		return cookie.Value, true
	}

	return "", false
}

func SetAuthCookies(w http.ResponseWriter, pair tokens.Pair) {
	http.SetCookie(w, cookie("access_token", pair.AccessToken, pair.AccessTokenExpiresAt))
	http.SetCookie(w, cookie("refresh_token", pair.RefreshToken, pair.RefreshTokenExpiresAt))
}

func ClearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, cookie("access_token", "", time.Unix(0, 0)))
	http.SetCookie(w, cookie("refresh_token", "", time.Unix(0, 0)))
}
