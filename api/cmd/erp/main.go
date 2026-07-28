package main

import (
	"context"
	"erp/internal/auth"
	"erp/internal/health"
	"erp/internal/shared/api"
	"erp/internal/shared/config"
	"erp/internal/shared/middlewares"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

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

	logger := slog.New(config.LogHandler())

	pool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		logger.ErrorContext(ctx, "creating database pooler", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.ErrorContext(ctx, "checking database connection", "error", err)
		os.Exit(1)
	}

	server := api.NewServer()

	server.Use(
		middlewares.Recoverer,
		middlewares.Logger(logger),
		middlewares.RequestID,
		middlewares.HTTPLogger,
		http.NewCrossOriginProtection().Handler,
		middlewares.Codec,
	)

	server.Group(health.Handlers())
	server.Group(auth.Handlers(pool, config.JWTSecret))

	server.Group(func(group *api.Group) {
		group.Use(middlewares.Auth(config.JWTSecret))
		group.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	addr := fmt.Sprintf(":%d", config.Port)
	logger.InfoContext(ctx, "starting server", slog.String("addr", addr))
	defer logger.InfoContext(ctx, "shutdown complete")

	if err := server.Serve(addr); err != nil {
		logger.ErrorContext(ctx, "starting server", "error", err)
	}
}
