package handlers

import (
	"net/http"
)

const AuthRegisterPattern = "POST /api/auth/register"

func AuthRegister(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
