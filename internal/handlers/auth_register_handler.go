package handlers

import (
	"erp/internal/dtos"
	"erp/internal/middlewares"
	"erp/internal/services"
	"errors"
	"net/http"
)

type AuthRegisterHandler struct {
	auth services.AuthService
}

var _ Handler = (*AuthRegisterHandler)(nil)

func NewAuthRegisterHandler(auth services.AuthService) Handler {
	return AuthRegisterHandler{auth: auth}
}

func (AuthRegisterHandler) __internal()     {}
func (AuthRegisterHandler) Pattern() string { return "POST /api/auth/register" }
func (handler AuthRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middlewares.GetLogger(ctx)
	codec := middlewares.GetCodec(ctx)

	var dto dtos.UserRegister
	if err := codec.Decode(r.Body, &dto); err != nil {
		logger.ErrorContext(ctx, "decoding request body", "error", err)
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	pair, err := handler.auth.Register(ctx, dto)
	if err != nil {
		logger.ErrorContext(ctx, "registering user", "error", err)
		switch err.(type) {
		case services.ErrCreatingUserModel:
			http.Error(w, errors.Unwrap(err).Error(), http.StatusUnprocessableEntity)
		case services.ErrUserExists:
			http.Error(w, errors.Unwrap(err).Error(), http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	setAuthCookies(w, pair)

	w.WriteHeader(http.StatusCreated)
	if err := codec.Encode(w, pair); err != nil {
		logger.ErrorContext(ctx, "encoding response body", "error", err)
	}
}
