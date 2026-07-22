package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{ResponseWriter: w}
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapper := newResponseWriterWrapper(w)

		ctx := r.Context()
		logger := GetLogger(ctx)
		logger.InfoContext(ctx, "entering request",
			slog.String("addr", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("uri", r.RequestURI))
		defer func() {
			logger.InfoContext(ctx, "exiting request", "status", fmt.Sprintf("%d %s", wrapper.statusCode, http.StatusText(wrapper.statusCode)))
		}()

		next.ServeHTTP(wrapper, r)
	})
}
