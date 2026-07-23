package handlers

import (
	"erp/internal/dtos"
	"erp/internal/middlewares"
	"erp/internal/services"
	"errors"
	"net/http"
)

type AuthLoginHandler struct {
	auth services.AuthService
}

var _ Handler = (*AuthLoginHandler)(nil)

func NewAuthLoginHandler(auth services.AuthService) Handler {
	return AuthLoginHandler{auth: auth}
}

func (AuthLoginHandler) __internal()     {}
func (AuthLoginHandler) Pattern() string { return "POST /api/auth/login" }
func (handler AuthLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middlewares.GetLogger(ctx)
	codec := middlewares.GetCodec(ctx)

	var dto dtos.UserLogin
	if err := codec.Decode(r.Body, &dto); err != nil {
		logger.ErrorContext(ctx, "decoding request body", "error", err)
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	pair, err := handler.auth.Login(ctx, dto)
	if err != nil {
		logger.ErrorContext(ctx, "logging in user", "error", err)
		switch err.(type) {
		case services.ErrCreatingUserModel:
			http.Error(w, errors.Unwrap(err).Error(), http.StatusUnprocessableEntity)
		case services.ErrUserNotExists:
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
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
