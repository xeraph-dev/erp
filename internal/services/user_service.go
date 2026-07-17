package services

import (
	"context"
	"erp/internal/dtos"
	"erp/internal/middlewares"
	"erp/internal/models"
	"erp/internal/repositories"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRegisteringUser       = errors.New("failed registering user")
	ErrLoginUser             = errors.New("failed login user")
	ErrUserNotExists         = errors.New("user does not exists")
	ErrIncorrectUserPassword = errors.New("incorrect user password")
)

type UserService interface {
	Service
	Register(ctx context.Context, in dtos.UserRegister) (out dtos.User, err error)
	Login(ctx context.Context, in dtos.UserLogin) (out dtos.User, err error)
}

type userServiceImpl struct {
	db    *pgxpool.Pool
	users repositories.UserRepository
}

var _ UserService = (*userServiceImpl)(nil)

func NewUserService(db *pgxpool.Pool, users repositories.UserRepository) UserService {
	return userServiceImpl{db, users}
}

func (service userServiceImpl) __internal() {}

func (service userServiceImpl) Register(ctx context.Context, in dtos.UserRegister) (out dtos.User, err error) {
	logger := middlewares.GetLogger(ctx)

	model, err := models.NewUserFromRegisterDTO(ctx, in)
	if err != nil {
		logger.ErrorContext(ctx, "creating model", "error", err)
		err = ErrRegisteringUser
		return
	}

	if err = withTx(ctx, service.db, func(tx pgx.Tx) (err error) {
		user, err := service.users.Create(ctx, tx, model)
		if err != nil {
			logger.ErrorContext(ctx, "creating entry", "error", err)
			err = ErrRegisteringUser
			return
		}

		out = user.DTO()
		return
	}); err != nil {
		return
	}

	return
}

func (service userServiceImpl) Login(ctx context.Context, in dtos.UserLogin) (out dtos.User, err error) {
	logger := middlewares.GetLogger(ctx)

	model, err := models.NewUserFromLoginDTO(in)
	if err != nil {
		return
	}

	var user models.User
	if err = withTx(ctx, service.db, func(tx pgx.Tx) (err error) {
		user, err = service.users.GetByName(ctx, tx, model)
		if errors.Is(err, repositories.ErrNotFound) {
			err = ErrUserNotExists
			return
		} else if err != nil {
			logger.ErrorContext(ctx, "getting user by name", "error", err)
			err = ErrLoginUser
			return
		}
		return
	}); err != nil {
		return
	}

	if !user.PasswordHash.Matches(in.Password) {
		err = ErrIncorrectUserPassword
		return
	}

	out = user.DTO()
	return
}
