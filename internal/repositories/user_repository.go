package repositories

import (
	"api/internal/database"
	"api/internal/models"
	"context"
)

type UserRepository interface {
	Repository
	Create(ctx context.Context, db database.Database, user models.User) (err error)
}

type UserRepositoryImpl struct{}

var _ UserRepository = UserRepositoryImpl{}

func (repo UserRepositoryImpl) __internal() {}

func (repo UserRepositoryImpl) Create(ctx context.Context, db database.Database, user models.User) (err error) {
	return
}
