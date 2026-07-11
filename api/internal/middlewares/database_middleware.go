package middlewares

import (
	"api/internal/database"
	"context"
	"net/http"
)

const DatabaseKey = "database"

func Database(db *database.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, DatabaseKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetDatabase(ctx context.Context) *database.Database {
	return ctx.Value(DatabaseKey).(*database.Database)
}
