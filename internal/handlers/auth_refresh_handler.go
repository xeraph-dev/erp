package handlers

import (
	"erp/internal/middlewares"
	"erp/internal/services"
	"errors"
	"net/http"
)

type AuthRefreshHandler struct {
	auth services.AuthService
}

var _ Handler = (*AuthRefreshHandler)(nil)

func NewAuthRefreshHandler(auth services.AuthService) Handler {
	return AuthRefreshHandler{auth: auth}
}

func (AuthRefreshHandler) __internal()     {}
func (AuthRefreshHandler) Pattern() string { return "POST /api/auth/refresh" }
func (handler AuthRefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middlewares.GetLogger(ctx)
	codec := middlewares.GetCodec(ctx)

	refreshToken, ok := extractRefreshToken(r, codec)
	if !ok {
		http.Error(w, "missing refresh token", http.StatusUnauthorized)
		return
	}

	pair, err := handler.auth.Refresh(ctx, refreshToken)
	if err != nil {
		logger.ErrorContext(ctx, "refreshing tokens", "error", err)
		switch {
		case errors.Is(err, services.ErrRefreshTokenNotFound):
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	setAuthCookies(w, pair)

	w.WriteHeader(http.StatusOK)
	if err := codec.Encode(w, pair); err != nil {
		logger.ErrorContext(ctx, "encoding response body", "error", err)
	}
}
