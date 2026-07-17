package middlewares

import (
	"context"
	"log/slog"
	"net/http"
)

const loggerKey = "logger"

func Logger(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, loggerKey, logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLogger(ctx context.Context) *slog.Logger {
	return ctx.Value(loggerKey).(*slog.Logger)
}

func SetLogger(ctx context.Context, logger *slog.LogValuer) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
