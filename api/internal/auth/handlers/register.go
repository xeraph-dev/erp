package handlers

import (
	"erp/internal/auth/dtos"
	"erp/internal/auth/helpers"
	"erp/internal/auth/services"
	"erp/internal/shared/middlewares"
	"errors"
	"net/http"
)

func Register(auth services.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		codec := middlewares.GetCodec(ctx)
		logger := middlewares.GetLogger(ctx)

		var dto dtos.UserRegister
		if err := codec.Decode(r.Body, &dto); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		pair, err := auth.Register(ctx, dto)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrRecordConflict):
				http.Error(w, err.Error(), http.StatusConflict)
			case errors.Is(err, services.ErrValidationFailed):
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			default:
				logger.ErrorContext(ctx, "register failed", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}

		helpers.SetAuthCookies(w, pair)
		w.Header().Set("Content-Type", codec.ContentType())
		w.WriteHeader(http.StatusCreated)
		if err := codec.Encode(w, pair); err != nil {
			return
		}
	}
}
