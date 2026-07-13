package middlewares

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

const DatabaseKey = "database"

func Database(pool *pgxpool.Pool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, DatabaseKey, pool)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetDatabase(ctx context.Context) *pgxpool.Pool {
	return ctx.Value(DatabaseKey).(*pgxpool.Pool)
}
