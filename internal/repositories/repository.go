package repositories

import (
	"context"
	"erp/db/queries"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	__internal()
}

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func SetCurrentUserID(ctx context.Context, db Querier, id uuid.UUID) (err error) {
	_, err = db.Exec(ctx, queries.SetCurrentUserID, id)
	return
}
