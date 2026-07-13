package services

import (
	"erp/internal/models"
	"erp/internal/repositories"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	Login(ctx context.Context) (err error)
}

type UserServiceImpl struct {
	repos repositories.Repositories
}

var _ UserService = UserServiceImpl{}

func (svr UserServiceImpl) Login(ctx context.Context) (err error) {
	db := ctx.Value("database").(*pgxpool.Pool)

	svr.repos.User.Create(ctx, db, models.User{})
	return
}
