package repositories

import (
	"context"
	"erp/db/queries"
	"erp/internal/middlewares"
	"erp/internal/models"
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrCreatingUser = errors.New("failed creating user")
	ErrGettingUser  = errors.New("failed getting user")
)

type UserRepository interface {
	Repository
	Create(ctx context.Context, db Querier, in models.User) (out models.User, err error)
	GetByName(ctx context.Context, db Querier, in models.User) (out models.User, err error)
}

type UserRepositoryImpl struct{}

var _ UserRepository = UserRepositoryImpl{}

func (UserRepositoryImpl) __internal() {}

func (UserRepositoryImpl) Create(ctx context.Context, db Querier, in models.User) (out models.User, err error) {
	logger := middlewares.GetLogger(ctx)

	rows, err := db.Query(ctx, queries.CreateUser, in.Username, in.PasswordHash, in.Email)
	if err != nil {
		logger.ErrorContext(ctx, "quering user", "error", err)
		err = ErrCreatingUser
		return
	}
	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		logger.ErrorContext(ctx, "collecting user", "error", err)
		err = ErrCreatingUser
		return
	}
	return
}

func (UserRepositoryImpl) GetByName(ctx context.Context, db Querier, in models.User) (out models.User, err error) {
	logger := middlewares.GetLogger(ctx)

	rows, err := db.Query(ctx, queries.GetUserByName, in.Username)
	if err != nil {
		logger.ErrorContext(ctx, "quering user", "error", err)
		err = ErrGettingUser
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		logger.ErrorContext(ctx, "colleting user", "error", err)
		err = ErrGettingUser
		return
	}

	return
}
