package main

import (
	"context"
	"erp/internal/health/handlers"
	"erp/internal/shared/api"
	"erp/internal/shared/config"
	"erp/internal/shared/middlewares"
	"fmt"
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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
	_ = pool

	server := api.NewServer(fmt.Sprintf(":%d", config.Port))

	logger := slog.New(config.LogHandler())

	server.Use(
		middlewares.Recoverer,
		middlewares.Logger(logger),
		middlewares.RequestID,
		middlewares.HTTPLogger,
		middlewares.Codec,
	)

	server.Handle(
		handlers.NewHealthHandler(),
	)

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
	log.Println("shutdown complete")
}
