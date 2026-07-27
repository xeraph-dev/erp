package handlers

import (
	"erp/internal/auth/helpers"
	"erp/internal/auth/services"
	"erp/internal/shared/middlewares"
	"errors"
	"net/http"
)

func Refresh(auth services.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		codec := middlewares.GetCodec(ctx)

		refreshToken, ok := helpers.ExtractRefreshToken(r, codec)
		if !ok {
			http.Error(w, "missing refresh token", http.StatusUnauthorized)
			return
		}

		pair, err := auth.Refresh(ctx, refreshToken)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrInvalidCredentials):
				helpers.ClearAuthCookies(w)
				http.Error(w, "invalid refresh token", http.StatusUnauthorized)
				return
			default:
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
