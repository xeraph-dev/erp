package handlers

import (
	"erp/internal/middlewares"
	"erp/internal/services"
	"errors"
	"net/http"
)

type AuthLogoutHandler struct {
	auth services.AuthService
}

var _ Handler = (*AuthLogoutHandler)(nil)

func NewAuthLogoutHandler(auth services.AuthService) Handler {
	return AuthLogoutHandler{auth: auth}
}

func (AuthLogoutHandler) __internal()     {}
func (AuthLogoutHandler) Pattern() string { return "POST /api/auth/logout" }
func (handler AuthLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middlewares.GetLogger(ctx)
	codec := middlewares.GetCodec(ctx)

	refreshToken, ok := extractRefreshToken(r, codec)
	if !ok {
		http.Error(w, "missing refresh token", http.StatusUnauthorized)
		return
	}

	if err := handler.auth.Logout(ctx, refreshToken); err != nil {
		logger.ErrorContext(ctx, "refreshing tokens", "error", err)
		switch {
		case errors.Is(err, services.ErrRefreshTokenNotFound):
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	clearAuthCookies(w)
	w.WriteHeader(http.StatusOK)
}
