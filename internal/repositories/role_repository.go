package repositories

import (
	"context"
	"erp/db/queries"
	"erp/internal/middlewares"
	"erp/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RoleRepository interface {
	Repository
	GetUser(ctx context.Context, db Querier) (out models.Role, err error)
	Assign(ctx context.Context, db Querier, roleID uuid.UUID, userID uuid.UUID) (err error)
}

type roleRepositoryImpl struct{}

var _ RoleRepository = (*roleRepositoryImpl)(nil)

func NewRoleRepository() RoleRepository {
	return roleRepositoryImpl{}
}

func (roleRepositoryImpl) __internal() {}

func (roleRepositoryImpl) GetUser(ctx context.Context, db Querier) (out models.Role, err error) {
	logger := middlewares.GetLogger(ctx)

	rows, err := db.Query(ctx, queries.GetRoleUser)
	if err != nil {
		logger.ErrorContext(ctx, "quering role user", "error", err)
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Role])
	if err != nil {
		logger.ErrorContext(ctx, "collecting role user", "error", err)
		return
	}

	return
}

func (roleRepositoryImpl) Assign(ctx context.Context, db Querier, roleID uuid.UUID, userID uuid.UUID) (err error) {
	logger := middlewares.GetLogger(ctx)

	rows, err := db.Query(ctx, queries.AssignRoleToUser, roleID, userID)
	if err != nil {
		logger.ErrorContext(ctx, "creating roles users entry", "error", err)
		return
	}

	_ , err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.RoleUser])
	if err != nil {
		logger.ErrorContext(ctx, "collecging role user row", "error", err)
		return
	}

	return
}
