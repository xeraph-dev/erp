package main

import (
	"context"
	"erp/internal/config"
	"erp/internal/handlers"
	"erp/internal/middlewares"
	"erp/internal/server"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/caarlos0/env/v10"
	_ "github.com/go-playground/validator/v10"
	_ "github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/spf13/viper"
	_ "github.com/stretchr/testify"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.ErrorContext(ctx, "config error", "error", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.ErrorContext(ctx, "database error", "error", err)
		os.Exit(1)
	}
	_ = pool

	s := server.NewServer()

	s.Use(middlewares.Logger(logger))
	s.Use(middlewares.RequestID)
	s.Use(middlewares.HTTPLogger)
	s.Use(middlewares.Recoverer)
	s.Use(middlewares.Database(pool))

	s.HandleFunc(handlers.AuthRegisterPattern, handlers.AuthRegister)

	addr := fmt.Sprintf(":%d", cfg.Port)
	logger.Info("starting server at", slog.String("addr", addr))
	if err := s.Serve(addr); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
