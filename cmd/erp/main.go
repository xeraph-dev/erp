package main

import (
	"context"
	"erp/internal/config"
	"erp/internal/controllers"
	"erp/internal/middlewares"
	"erp/internal/repositories"
	"erp/internal/server"
	"erp/internal/services"
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

	server := server.NewServer()

	server.Use(
		middlewares.Logger(logger),
		middlewares.RequestID,
		middlewares.HTTPLogger,
		middlewares.Recoverer,
	)

	server.Database(pool)

	server.RepoRegister(
		repositories.UserRepositoryImpl{},
	)

	server.ServiceRegister(
		&services.UserServiceImpl{},
	)

	server.Add(
		&controllers.AuthRegisterController{},
	)

	addr := fmt.Sprintf(":%d", cfg.Port)
	logger.Info("starting server", slog.String("addr", addr))
	if err := server.Serve(addr); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
