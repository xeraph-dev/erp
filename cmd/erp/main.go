package main

import (
	"context"
	"erp/internal/config"
	"erp/internal/handlers"
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
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	_ "github.com/redis/go-redis/v9"
	_ "github.com/spf13/viper"
	_ "github.com/stretchr/testify"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.ErrorContext(ctx, "dotenv error", "error", err)
		os.Exit(1)
	}

	config, err := config.New()
	if err != nil {
		logger.ErrorContext(ctx, "config error", "error", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		logger.ErrorContext(ctx, "database error", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	server := server.NewServer()

	server.Use(
		middlewares.Logger(logger),
		middlewares.Recoverer,
		middlewares.RequestID,
		middlewares.HTTPLogger,
		middlewares.Codec,
	)

	userRepo := repositories.NewUserRepository()

	userService := services.NewUserService(pool, userRepo)
	jwtService := services.NewJWTService(config.JWTSecret)

	server.Handle(
		handlers.NewAuthRegisterHandler(userService, jwtService),
	)

	addr := fmt.Sprintf(":%d", config.Port)
	logger.Info("starting server", slog.String("addr", addr))
	if err := server.Serve(addr); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
	logger.Info("shutdown complete")
}
