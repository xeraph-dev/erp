package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TxBeginner interface {
	Querier
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Tx interface {
	TxBeginner
}

func WithTx(ctx context.Context, b TxBeginner, txFunc func(tx Tx) (err error)) (err error) {
	tx, err := b.Begin(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	err = txFunc(tx)
	return
}

func SetCurrentUserID(ctx context.Context, q Querier, id uuid.UUID) (err error) {
	_, err = q.Exec(ctx, "SET LOCAL app.current_user_id TO $1", id)
	return
}

func QueryExactlyOneRow[T any](op, query string, args pgx.StrictNamedArgs, ctx context.Context, q Querier, translateErrFunc func(error) error) (out T, err error) {
	rows, err := q.Query(ctx, query, args)
	if err != nil {
		return out, fmt.Errorf("%s: %w", op, err)
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		return out, fmt.Errorf("%s: %w", op, translateErrFunc(err))
	}

	return
}

func QueryRows[T any](op, query string, args pgx.StrictNamedArgs, ctx context.Context, q Querier, translateErrFunc func(error) error) (out []T, err error) {
	rows, err := q.Query(ctx, query, args)
	if err != nil {
		return out, fmt.Errorf("%s: %w", op, err)
	}

	out, err = pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return out, fmt.Errorf("%s: %w", op, translateErrFunc(err))
	}

	return
}
