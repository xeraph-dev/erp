package stores

import (
	_ "embed"

	"context"
	"erp/internal/auth/models"
	"erp/internal/auth/vos"
	"erp/internal/shared/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type User interface {
	GetByID(ctx context.Context, q db.Querier, id uuid.UUID) (out models.User, err error)
	GetByUsername(ctx context.Context, q db.Querier, username vos.Username) (out models.User, err error)
	GetByEmail(ctx context.Context, q db.Querier, email vos.Email) (out models.User, err error)
	Create(ctx context.Context, q db.Querier, in models.User) (out models.User, err error)
	DeleteByID(ctx context.Context, q db.Querier, id uuid.UUID) (out models.User, err error)
}

type userImpl struct{}

var _ User = (*userImpl)(nil)

//go:embed queries/get_user_by_id.sql
var getUserByID string

func (store userImpl) GetByID(ctx context.Context, q db.Querier, id uuid.UUID) (out models.User, err error) {
	rows, err := q.Query(ctx, getUserByID, id)
	if err != nil {
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return
	}

	return
}

//go:embed queries/get_user_by_username.sql
var getUserByUsername string

func (store userImpl) GetByUsername(ctx context.Context, q db.Querier, username vos.Username) (out models.User, err error) {
	rows, err := q.Query(ctx, getUserByUsername, username)
	if err != nil {
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return
	}

	return
}

//go:embed queries/get_user_by_email.sql
var getUserByEmail string

func (store userImpl) GetByEmail(ctx context.Context, q db.Querier, email vos.Email) (out models.User, err error) {
	rows, err := q.Query(ctx, getUserByEmail, email)
	if err != nil {
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return
	}

	return
}

//go:embed queries/create_user.sql
var createUser string

func (store userImpl) Create(ctx context.Context, q db.Querier, in models.User) (out models.User, err error) {
	rows, err := q.Query(ctx, createUser, in.Username, in.Email, in.PasswordHash, in.FirstName, in .LastName)
	if err != nil {
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return
	}

	return
}


//go:embed queries/delete_user_by_id.sql
var deleteUserByID string

func (store userImpl) DeleteByID(ctx context.Context, q db.Querier, id uuid.UUID) (out models.User, err error) {
	rows, err := q.Query(ctx, deleteUserByID, id)
	if err != nil {
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return
	}

	return
}
