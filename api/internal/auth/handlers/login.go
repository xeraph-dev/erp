package handlers

import (
	"erp/internal/auth/dtos"
	"erp/internal/auth/helpers"
	"erp/internal/auth/services"
	"erp/internal/shared/middlewares"
	"errors"
	"net/http"
)

func Login(auth services.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		codec := middlewares.GetCodec(ctx)
		logger := middlewares.GetLogger(ctx)

		var dto dtos.UserLogin
		if err := codec.Decode(r.Body, &dto); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		pair, err := auth.Login(ctx, dto)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrInvalidCredentials):
				http.Error(w, "invalid credentials", http.StatusUnauthorized)
				return
			default:
				logger.ErrorContext(ctx, "login failed", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}

		helpers.SetAuthCookies(w, pair)
		w.Header().Set("Content-Type", codec.ContentType())
		w.WriteHeader(http.StatusOK)
		if err := codec.Encode(w, pair); err != nil {
			return
		}
	}
}
