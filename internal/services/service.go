package services

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	__internal()
}

func withTx(ctx context.Context, db *pgxpool.Pool, txFunc func(tx pgx.Tx) (err error)) (err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	return txFunc(tx)
}
