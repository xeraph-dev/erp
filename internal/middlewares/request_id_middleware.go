package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"
const RequestIDKey = "request-id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New()
		r.Header.Set(RequestIDHeader, requestID.String())

		ctx := r.Context()
		ctx = context.WithValue(ctx, RequestIDKey, requestID)

		logger := ctx.Value(LoggerKey).(*slog.Logger)
		logger = logger.With(RequestIDKey, requestID)
		ctx = context.WithValue(ctx, LoggerKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) uuid.UUID {
	return ctx.Value(RequestIDKey).(uuid.UUID)
}
