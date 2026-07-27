package handlers

import (
	"erp/internal/auth/helpers"
	"erp/internal/auth/services"
	"erp/internal/auth/stores"
	"erp/internal/shared/middlewares"
	"errors"
	"net/http"
)

func Logout(auth services.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		codec := middlewares.GetCodec(ctx)
		logger := middlewares.GetLogger(ctx)

		refreshToken, ok := helpers.ExtractRefreshToken(r, codec)
		if !ok {
			http.Error(w, "missing refresh token", http.StatusUnauthorized)
			return
		}

		if err := auth.Logout(ctx, refreshToken); err != nil && !errors.Is(err, stores.ErrRefreshTokenNotFound) {
			logger.ErrorContext(ctx, "logout failed", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return

		}

		helpers.ClearAuthCookies(w)
		w.WriteHeader(http.StatusOK)
	}
}
