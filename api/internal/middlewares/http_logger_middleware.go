package middlewares

import (
	"log/slog"
	"net/http"
)

func HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := GetLogger(ctx)
		logger.InfoContext(ctx, "entering request",
			slog.String("addr", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("uri", r.RequestURI))
		// TODO: log response's status code
		defer logger.InfoContext(ctx, "exiting request")

		next.ServeHTTP(w, r)
	})
}
