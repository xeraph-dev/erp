package repositories

import (
	"context"
	"erp/db/queries"
	"erp/internal/models"
	"fmt"

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
	rows, err := db.Query(ctx, queries.GetRoleUser)
	if err != nil {
		err = fmt.Errorf("quering role user: %w", err)
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Role])
	if err != nil {
		err = fmt.Errorf("collecting role user: %w", err)
		return
	}

	return
}

func (roleRepositoryImpl) Assign(ctx context.Context, db Querier, roleID uuid.UUID, userID uuid.UUID) (err error) {
	rows, err := db.Query(ctx, queries.AssignRoleToUser, roleID, userID)
	if err != nil {
		err = fmt.Errorf("creating roles users entry: %w", err)
		return
	}

	_, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.RoleUser])
	if err != nil {
		err = fmt.Errorf("collecging role user row: %w", err)
		return
	}

	return
}
