package stores

import (
	_ "embed"
	"errors"

	"context"
	"erp/internal/auth/models"
	"erp/internal/auth/vos"
	"erp/internal/shared/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const roleLayer = "stores/Role"

var (
	ErrRoleNotFound   = errors.New("role not found")
	ErrRolenameExists = errors.New("role name already exists")
)

type Role interface {
	GetByID(ctx context.Context, q db.Querier, id uuid.UUID) (out models.Role, err error)
	GetByName(ctx context.Context, q db.Querier, name vos.RoleName) (out models.Role, err error)
	Create(ctx context.Context, q db.Querier, in models.Role) (out models.Role, err error)
	DeleteByID(ctx context.Context, q db.Querier, id uuid.UUID) (out models.Role, err error)
	Assign(ctx context.Context, q db.Querier, roleID uuid.UUID, userID uuid.UUID) (err error)
	AssignUser(ctx context.Context, q db.Querier, userID uuid.UUID) (err error)
	AssignAdmin(ctx context.Context, q db.Querier, userID uuid.UUID) (err error)
	ListByUserID(ctx context.Context, q db.Querier, userID uuid.UUID) (out []models.Role, err error)
}

type roleImpl struct{}

var _ Role = (*roleImpl)(nil)

func NewRole() Role {
	return roleImpl{}
}

//go:embed queries/get_role_by_id.sql
var getRoleByID string

func (store roleImpl) GetByID(ctx context.Context, q db.Querier, id uuid.UUID) (out models.Role, err error) {
	return db.QueryExactlyOneRow[models.Role](roleLayer+".GetByID", getRoleByID, pgx.StrictNamedArgs{
		"id": id,
	}, ctx, q, store.translate)
}

//go:embed queries/get_role_by_name.sql
var getRoleByName string

func (store roleImpl) GetByName(ctx context.Context, q db.Querier, name vos.RoleName) (out models.Role, err error) {
	return db.QueryExactlyOneRow[models.Role](roleLayer+".GetByName", getRoleByName, pgx.StrictNamedArgs{
		"role_name": name,
	}, ctx, q, store.translate)
}

//go:embed queries/create_role.sql
var createRole string

func (store roleImpl) Create(ctx context.Context, q db.Querier, in models.Role) (out models.Role, err error) {
	return db.QueryExactlyOneRow[models.Role](roleLayer+".Create", createRole, pgx.StrictNamedArgs{
		"role_name":  in.Name,
		"role_level": in.Level,
	}, ctx, q, store.translate)
}

//go:embed queries/delete_role_by_id.sql
var deleteRoleByID string

func (store roleImpl) DeleteByID(ctx context.Context, q db.Querier, id uuid.UUID) (out models.Role, err error) {
	return db.QueryExactlyOneRow[models.Role](roleLayer+".DeleteByID", deleteRoleByID, pgx.StrictNamedArgs{
		"id": id,
	}, ctx, q, store.translate)
}

//go:embed queries/assign_role.sql
var assignRole string

func (store roleImpl) Assign(ctx context.Context, q db.Querier, roleID uuid.UUID, userID uuid.UUID) (err error) {
	return db.Exec(roleLayer+".Assign", assignRole, pgx.StrictNamedArgs{
		"role_id": roleID,
		"user_id": userID,
	}, ctx, q, store.translate)
}

//go:embed queries/assign_role_user.sql
var assignRoleUser string

func (store roleImpl) AssignUser(ctx context.Context, q db.Querier, userID uuid.UUID) (err error) {
	return db.Exec(roleLayer+".AssignUser", assignRoleUser, pgx.StrictNamedArgs{
		"user_id": userID,
	}, ctx, q, store.translate)
}

//go:embed queries/assign_role_admin.sql
var assignRoleAdmin string

func (store roleImpl) AssignAdmin(ctx context.Context, q db.Querier, userID uuid.UUID) (err error) {
	return db.Exec(roleLayer+".AssignAdmin", assignRoleUser, pgx.StrictNamedArgs{
		"user_id": userID,
	}, ctx, q, store.translate)
}

//go:embed queries/get_roles_by_user_id.sql
var getRolesByUserID string

func (store roleImpl) ListByUserID(ctx context.Context, q db.Querier, userID uuid.UUID) (out []models.Role, err error) {
	return db.QueryRows[models.Role](roleLayer+".ListByUserID", getRolesByUserID, pgx.StrictNamedArgs{
		"user_id": userID,
	}, ctx, q, store.translate)
}

func (store roleImpl) translate(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrRoleNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "active_roles_role_name_idx":
			return ErrRolenameExists
		}
	}

	return err
}
