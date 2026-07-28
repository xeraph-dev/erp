package stores

import (
	"context"
	_ "embed"
	"erp/internal/auth/models"
	"erp/internal/shared/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const permissionLayer = "stores/Permission"

type Permission interface {
	ListByRoleID(ctx context.Context, q db.Querier, roleID uuid.UUID) (out []models.Permission, err error)
	ListByUserID(ctx context.Context, q db.Querier, userID uuid.UUID) (out []models.Permission, err error)
}

type permissionImpl struct{}

var _ Permission = (*permissionImpl)(nil)

func NewPermission() Permission {
	return permissionImpl{}
}

//go:embed queries/get_permissions_by_user_id.sql
var getPermissionsByUserID string

func (store permissionImpl) ListByRoleID(ctx context.Context, q db.Querier, roleID uuid.UUID) (out []models.Permission, err error) {
	return db.QueryRows[models.Permission](permissionLayer+".ListByRoleID", getPermissionsByUserID, pgx.StrictNamedArgs{
		"user_id": roleID,
	}, ctx, q, store.translate)
}

//go:embed queries/get_permissions_by_user_id.sql
var getPermissionsByRoleID string

func (store permissionImpl) ListByUserID(ctx context.Context, q db.Querier, userID uuid.UUID) (out []models.Permission, err error) {
	return db.QueryRows[models.Permission](permissionLayer+".ListByUserID", getPermissionsByRoleID, pgx.StrictNamedArgs{
		"user_id": userID,
	}, ctx, q, store.translate)
}

func (store permissionImpl) translate(err error) error {
	return err
}
