package repositories

import (
	"context"
	"erp/db/queries"
	"erp/internal/models"
	"erp/internal/vos"
	"errors"
	"fmt"

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
	rows, err := db.Query(ctx, queries.CreateUser, in.Username, in.PasswordHash, in.Email)
	if err != nil {
		err = fmt.Errorf("creating user entry: %w", err)
		return
	}
	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		err = fmt.Errorf("collecting user row: %w", err)
		return
	}
	return
}

func (userRepositoryImpl) GetByUsername(ctx context.Context, db Querier, username vos.Username) (out models.User, err error) {
	rows, err := db.Query(ctx, queries.GetUserByUsername, username)
	if err != nil {
		err = fmt.Errorf("quering user by username: %w", err)
		return
	}
	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		err = fmt.Errorf("collecting user row: %w", err)
		return
	}
	return
}

func (userRepositoryImpl) UsernameExists(ctx context.Context, db Querier, username vos.Username) (exists bool, err error) {
	err = db.QueryRow(ctx, queries.UsernameExists, username).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
		return
	} else if err != nil {
		err = fmt.Errorf("quering username exists: %w", err)
		return
	}

	return
}

func (userRepositoryImpl) EmailExists(ctx context.Context, db Querier, email vos.Email) (exists bool, err error) {
	err = db.QueryRow(ctx, queries.UserEmailExists, email).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
		return
	} else if err != nil {
		err = fmt.Errorf("quering user email exists: %w", err)
		return
	}

	return
}
