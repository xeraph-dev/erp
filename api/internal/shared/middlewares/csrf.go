package middlewares

import (
	"crypto/subtle"
	"net/http"
)

func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if isSafeMethod(r.Method) || GetAuthTransport(ctx) != authTransportCookie {
			next.ServeHTTP(w, r)
			return
		}

		switch r.Header.Get("Sec-Fetch-Site") {
		case "same-site", "same-origin":
			next.ServeHTTP(w, r)
			return
		case "cross-site":
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		headerToken := r.Header.Get("X-CSRF-Token")
		cookieToken, err := r.Cookie("csrf_token")
		if headerToken == "" || err != nil || cookieToken.Value == "" {
			http.Error(w, "missing CSRF token", http.StatusForbidden)
			return
		}
		if subtle.ConstantTimeCompare([]byte(headerToken), []byte(cookieToken.Value)) != 1 {
			http.Error(w, "invalid CSRF token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}
