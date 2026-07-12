package main

import (
	"api/internal/handlers"
	"api/internal/helpers"
	"api/internal/middlewares"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/go-playground/validator/v10"
	_ "github.com/golang-jwt/jwt/v5"
	_ "github.com/spf13/viper"
	_ "github.com/stretchr/testify"
	_ "github.com/golang-migrate/migrate/v4"
)

const Host = ""
const Port = 8080
const DatabaseUrl = "postgresql://xeraph@localhost:5432/erp?sslmode=disable"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	s := helpers.NewServer()

	s.Use(middlewares.Logger(logger))
	s.Use(middlewares.RequestID)
	s.Use(middlewares.HTTPLogger)
	s.Use(middlewares.Recoverer)

	s.HandleFunc(handlers.AuthRegisterPattern, handlers.AuthRegister)

	addr := fmt.Sprintf("%s:%d", Host, Port)
	logger.Info("starting server at", slog.String("addr", addr))
	if err := s.Serve(addr); err != nil {
		logger.Error("failed to start server", "error", err)
	}
}
