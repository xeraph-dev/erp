package repositories

import (
	"context"
	"erp/db/queries"
	"erp/internal/middlewares"
	"erp/internal/models"
	"erp/internal/vos"
	"errors"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Repository
	Create(ctx context.Context, db Querier, in models.User) (out models.User, err error)
	GetByUsername(ctx context.Context, db Querier, username vos.Username) (out models.User, err error)
	UsernameExists(ctx context.Context, db Querier, username vos.Username) (exists bool, err error)
	EmailExists(ctx context.Context, db Querier, email vos.Email) (exists bool, err error)
}

type userRepositoryImpl struct{}

var _ UserRepository = (*userRepositoryImpl)(nil)

func NewUserRepository() UserRepository { return userRepositoryImpl{} }

func (userRepositoryImpl) __internal() {}

func (userRepositoryImpl) Create(ctx context.Context, db Querier, in models.User) (out models.User, err error) {
	logger := middlewares.GetLogger(ctx)

	rows, err := db.Query(ctx, queries.CreateUser, in.Username, in.PasswordHash, in.Email)
	if err != nil {
		logger.ErrorContext(ctx, "creating user entry", "error", err)
		return
	}
	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		logger.ErrorContext(ctx, "collecting user row", "error", err)
		return
	}
	return
}

func (userRepositoryImpl) GetByUsername(ctx context.Context, db Querier, username vos.Username) (out models.User, err error) {
	logger := middlewares.GetLogger(ctx)

	rows, err := db.Query(ctx, queries.GetUserByUsername, username)
	if err != nil {
		logger.ErrorContext(ctx, "quering user by username", "error", err)
		return
	}
	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		logger.ErrorContext(ctx, "collecting user row", "error", err)
		return
	}
	return
}

func (userRepositoryImpl) UsernameExists(ctx context.Context, db Querier, username vos.Username) (exists bool, err error) {
	logger := middlewares.GetLogger(ctx)

	err = db.QueryRow(ctx, queries.UsernameExists, username).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
		return
	} else if err != nil {
		logger.ErrorContext(ctx, "quering username exists", "error", err)
		return
	}

	return
}

func (userRepositoryImpl) EmailExists(ctx context.Context, db Querier, email vos.Email) (exists bool, err error) {
	logger := middlewares.GetLogger(ctx)

	err = db.QueryRow(ctx, queries.UserEmailExists, email).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
		return
	} else if err != nil {
		logger.ErrorContext(ctx, "quering user email exists", "error", err)
		return
	}

	return
}
