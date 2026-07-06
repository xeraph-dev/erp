package main

import (
	"net/http"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/golang-jwt/jwt/v5"
	_ "github.com/go-playground/validator/v10"
	_ "github.com/spf13/viper"
	_ "github.com/stretchr/testify"
)

func main() {
	http.HandleFunc("GET /api/users", func(w http.ResponseWriter, r *http.Request) {
		println("GET /api/users")
	})
	http.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		println("POST /api/users")
	})
	http.HandleFunc("GET /api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		println("GET /api/users/{id}")
	})
	http.HandleFunc("PATCH /api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		println("PATCH /api/users/{id}")
	})
	http.HandleFunc("DELETE /api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		println("DELETE /api/users/{id}")
	})

	http.ListenAndServe(":8080", nil)
}
