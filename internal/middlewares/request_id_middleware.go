package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"
const requestIDKey = "request-id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New()
		r.Header.Set(requestIDHeader, requestID.String())
		w.Header().Set(requestIDHeader, requestID.String())

		ctx := r.Context()
		ctx = context.WithValue(ctx, requestIDKey, requestID)

		logger := ctx.Value(loggerKey).(*slog.Logger)
		logger = logger.With(requestIDKey, requestID)
		ctx = context.WithValue(ctx, loggerKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) uuid.UUID {
	return ctx.Value(requestIDKey).(uuid.UUID)
}
