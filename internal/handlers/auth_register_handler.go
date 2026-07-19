package handlers

import (
	"erp/internal/dtos"
	"erp/internal/middlewares"
	"erp/internal/services"
	"net/http"
)

type AuthRegisterHandler struct {
	users services.UserService
	jwt   services.JWTService
}

var _ Handler = (*AuthRegisterHandler)(nil)

func NewAuthRegisterHandler(users services.UserService, jwt services.JWTService) Handler {
	return AuthRegisterHandler{users, jwt}
}

func (AuthRegisterHandler) __internal()     {}
func (AuthRegisterHandler) Pattern() string { return "POST /api/auth/register" }
func (handler AuthRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middlewares.GetLogger(ctx)
	codec := middlewares.GetCodec(ctx)

	var dto dtos.UserRegister
	if err := codec.Decode(r.Body, &dto); err != nil {
		logger.ErrorContext(ctx, "decoding request", "error", err)
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	user, err := handler.users.Register(ctx, dto)
	if err != nil {
		logger.ErrorContext(ctx, "registering user", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "access-token",
	})

	// w.Header().Add("Set-Cookie", "")

	w.WriteHeader(http.StatusCreated)
	if err := codec.Encode(w, user); err != nil {
		logger.ErrorContext(ctx, "encoding response", "error", err)
	}
}
