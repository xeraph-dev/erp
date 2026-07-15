package controllers

import (
	"encoding/json"
	"erp/internal/dtos"
	"erp/internal/middlewares"
	"erp/internal/services"
	"errors"
	"net/http"
)

type AuthRegisterController struct {
	Services struct {
		User services.UserService
	}
}

var _ Controller = AuthRegisterController{}

func (AuthRegisterController) __internal()     {}
func (AuthRegisterController) Pattern() string { return "POST /api/auth/register" }
func (controller AuthRegisterController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middlewares.GetLogger(ctx)

	var dto dtos.UserRegister

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		logger.ErrorContext(ctx, "decoding request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := controller.Services.User.Register(ctx, dto)
	if errors.Is(err, services.ErrRegisteringUser) {
		logger.ErrorContext(ctx, "registering user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		logger.ErrorContext(ctx, "encoding response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
