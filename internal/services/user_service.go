package services

import (
	"context"
	"erp/internal/dtos"
	"erp/internal/middlewares"
	"erp/internal/models"
	"erp/internal/repositories"
	"errors"

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

type UserServiceImpl struct {
	DB    *pgxpool.Pool
	Repos struct {
		User repositories.UserRepository
	}
}

var _ UserService = UserServiceImpl{}

func (service UserServiceImpl) __internal() {}

func (service UserServiceImpl) Register(ctx context.Context, in dtos.UserRegister) (out dtos.User, err error) {
	logger := middlewares.GetLogger(ctx)

	model, err := models.NewUserFromRegisterDTO(ctx, in)
	if err != nil {
		logger.ErrorContext(ctx, "creating model", "error", err)
		err = ErrRegisteringUser
		return
	}

	user, err := service.Repos.User.Create(ctx, service.DB, model)
	if err != nil {
		logger.ErrorContext(ctx, "creating entry", "error", err)
		err = ErrRegisteringUser
		return
	}

	out = user.DTO()
	return
}

func (service UserServiceImpl) Login(ctx context.Context, in dtos.UserLogin) (out dtos.User, err error) {
	logger := middlewares.GetLogger(ctx)

	model := models.NewUserFromLoginDTO(in)

	user, err := service.Repos.User.GetByName(ctx, service.DB, model)
	if err != nil {
		logger.ErrorContext(ctx, "getting user by name", "error", err)
		err = ErrUserNotExists
		return
	}

	if !user.PasswordHash.Matches(in.Password) {
		err = ErrIncorrectUserPassword
		return
	}

	out = user.DTO()
	return
}
