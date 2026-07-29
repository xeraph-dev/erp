package middlewares

import (
	"erp/internal/shared/db"
	"net/http"
)

func RequirePermissions(permissions ...db.Permission) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !GetUserPermissions(r.Context()).Satisfies(permissions...) {
				http.Error(w, "insufficent permissions", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
