package controllers

import (
	"erp/internal/dtos"
	"erp/internal/middlewares"
	"erp/internal/services"
	"net/http"
)

type AuthRegisterController struct {
	users services.UserService
}

var _ Controller = (*AuthRegisterController)(nil)

func NewAuthRegisterController(users services.UserService) Controller {
	return AuthRegisterController{users}
}

func (AuthRegisterController) __internal()     {}
func (AuthRegisterController) Pattern() string { return "POST /api/auth/register" }
func (controller AuthRegisterController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middlewares.GetLogger(ctx)
	codec := middlewares.GetCodec(ctx)

	var dto dtos.UserRegister
	if err := codec.Decode(r.Body, &dto); err != nil {
		logger.ErrorContext(ctx, "decoding request", "error", err)
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	user, err := controller.users.Register(ctx, dto)
	if err != nil {
		logger.ErrorContext(ctx, "registering user", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// w.Header().Add("Set-Cookie", "")

	w.WriteHeader(http.StatusCreated)
	if err := codec.Encode(w, user); err != nil {
		logger.ErrorContext(ctx, "encoding response", "error", err)
	}
}
