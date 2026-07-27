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

const userLayer = "stores/User"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUsernameTaken  = errors.New("username already taken")
	ErrUserEmailTaken = errors.New("email already taken")
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

func NewUser() User {
	return userImpl{}
}

//go:embed queries/get_user_by_id.sql
var getUserByID string

func (store userImpl) GetByID(ctx context.Context, q db.Querier, id uuid.UUID) (out models.User, err error) {
	return db.QueryExactlyOneRow[models.User](userLayer+".GetByID", getUserByID, pgx.StrictNamedArgs{
		"id": id,
	}, ctx, q, store.translate)
}

//go:embed queries/get_user_by_username.sql
var getUserByUsername string

func (store userImpl) GetByUsername(ctx context.Context, q db.Querier, username vos.Username) (out models.User, err error) {
	return db.QueryExactlyOneRow[models.User](userLayer+".GetByUsername", getUserByUsername, pgx.StrictNamedArgs{
		"username": username,
	}, ctx, q, store.translate)
}

//go:embed queries/get_user_by_email.sql
var getUserByEmail string

func (store userImpl) GetByEmail(ctx context.Context, q db.Querier, email vos.Email) (out models.User, err error) {
	return db.QueryExactlyOneRow[models.User](userLayer+".GetByEmail", getUserByEmail, pgx.StrictNamedArgs{
		"email": email,
	}, ctx, q, store.translate)
}

//go:embed queries/create_user.sql
var createUser string

func (store userImpl) Create(ctx context.Context, q db.Querier, in models.User) (out models.User, err error) {
	return db.QueryExactlyOneRow[models.User](userLayer+".Create", createUser, pgx.StrictNamedArgs{
		"username":      in.Username,
		"email":         in.Email,
		"password_hash": in.PasswordHash,
		"first_name":    in.FirstName,
		"last_name":     in.LastName,
	}, ctx, q, store.translate)
}

//go:embed queries/delete_user_by_id.sql
var deleteUserByID string

func (store userImpl) DeleteByID(ctx context.Context, q db.Querier, id uuid.UUID) (out models.User, err error) {
	return db.QueryExactlyOneRow[models.User](userLayer+".DeleteByID", deleteUserByID, pgx.StrictNamedArgs{
		"id": id,
	}, ctx, q, store.translate)
}

func (store userImpl) translate(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "active_users_username_idx":
			return ErrUsernameTaken
		case "active_users_email_idx":
			return ErrUserEmailTaken
		}
	}

	return err
}
