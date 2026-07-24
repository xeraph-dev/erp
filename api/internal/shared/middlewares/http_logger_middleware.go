package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapper := &responseWriterWrapper{ResponseWriter: w}
		wrapper.statusCode = 200

		ctx := r.Context()
		logger := GetLogger(ctx)
		start := time.Now()

		defer func() {
			logger.InfoContext(ctx, "request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("status", fmt.Sprintf("%d %s", wrapper.statusCode, http.StatusText(wrapper.statusCode))),
				slog.Duration("time", time.Since(start)),
				slog.String("remove_addr", r.RemoteAddr),
			)
		}()

		next.ServeHTTP(wrapper, r)
	})
}
