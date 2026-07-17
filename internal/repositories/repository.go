package repositories

import (
	"context"
	"erp/db/queries"

	"github.com/google/uuid"
)

type Repository interface {
	__internal()
}

func SetCurrentUserID(ctx context.Context, db Querier, id uuid.UUID) (err error) {
	_, err = db.Exec(ctx, queries.SetCurrentUserID, id)
	return
}
