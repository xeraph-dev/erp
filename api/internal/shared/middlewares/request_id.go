package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type requestIDCtxKet int

const requestIDKey requestIDCtxKet = iota

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestIDString := r.Header.Get("X-Request-ID")
		if requestIDString == "" {
			requestIDString = uuid.NewString()
		}
		requestID, err := uuid.Parse(requestIDString)
		if err != nil {
			http.Error(w, "malformed X-Request-ID header", http.StatusBadRequest)
			return
		}

		w.Header().Set("X-Request-ID", requestID.String())

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		logger := GetLogger(ctx).With("request_id", requestID)
		ctx = context.WithValue(ctx, loggerKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(requestIDKey).(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}
