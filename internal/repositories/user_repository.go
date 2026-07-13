package repositories

import (
	"context"
	"erp/internal/models"
)

type UserRepository interface {
	Repository
	Create(ctx context.Context, db Querier, user models.User) (err error)
}

type UserRepositoryImpl struct{}

var _ UserRepository = UserRepositoryImpl{}

func (repo UserRepositoryImpl) __internal() {}

func (repo UserRepositoryImpl) Create(ctx context.Context, db Querier, user models.User) (err error) {
	return
}
