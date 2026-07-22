package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

const loggerKey = "logger"

func Logger() Middleware {
	// TODO: handler logger target via config
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, loggerKey, logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLogger(ctx context.Context) (logger *slog.Logger) {
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		logger = slog.Default()
	}
	return
}
