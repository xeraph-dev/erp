package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	ctx    context.Context
	logger *slog.Logger
	pool   *pgxpool.Pool
}

func NewDB(ctx context.Context, logger *slog.Logger, connString string) (db *Database, err error) {
	db = new(Database)
	db.ctx = ctx
	db.logger = logger.With(slog.String("component", "DB"))

	if db.pool, err = pgxpool.New(db.ctx, connString); err != nil {
		db = nil
		return
	}

	return
}

func (db *Database) Close() {
	db.pool.Close()
}
