package main

import (
	"context"
	"erp/internal/api"
	"erp/internal/config"
	"erp/internal/handlers"
	"erp/internal/middlewares"
	"erp/internal/repositories"
	"erp/internal/services"
	"fmt"
	"log"

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

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	addr := fmt.Sprintf(":%d", config.Port)
	server := api.NewServer(addr)

	server.Use(
		middlewares.Recoverer,
		middlewares.Logger(),
		middlewares.RequestID,
		middlewares.HTTPLogger,
		middlewares.Codec,
	)

	userRepo := repositories.NewUserRepository()
	roleRepo := repositories.NewRoleRepository()
	refreshTokenRepo := repositories.NewRefreshTokenRepository()

	authService := services.NewAuthService(config.JWTSecret, pool, userRepo, roleRepo, refreshTokenRepo)

	server.Handle(handlers.NewAuthRegisterHandler(authService))
	server.Handle(handlers.NewAuthLoginHandler(authService))

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
	log.Println("shutdown complete")
}
