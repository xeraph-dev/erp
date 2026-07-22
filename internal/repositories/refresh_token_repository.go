package repositories

import (
	"context"
	"erp/db/queries"
	"erp/internal/middlewares"
	"erp/internal/models"

	"github.com/jackc/pgx/v5"
)

type RefreshTokenRepository interface {
	Repository
	Create(ctx context.Context, db Querier, in models.RefreshToken) (err error)
}

type refreshTokenRepositoryImpl struct{}

var _ RefreshTokenRepository = (*refreshTokenRepositoryImpl)(nil)

func NewRefreshTokenRepository() RefreshTokenRepository {
	return refreshTokenRepositoryImpl{}
}

func (refreshTokenRepositoryImpl) __internal() {}

func (repo refreshTokenRepositoryImpl) Create(ctx context.Context, db Querier, in models.RefreshToken) (err error) {
	logger := middlewares.GetLogger(ctx)

	rows, err := db.Query(ctx, queries.CreateRefreshToken, in.UserID, in.FamilyID, in.TokenHash, in.ExpiresAt)
	if err != nil {
		logger.ErrorContext(ctx, "creating refresh token entry", "error", err)
		return
	}

	_, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.RefreshToken])
	if err != nil {
		logger.ErrorContext(ctx, "collecting refresh token row", "error", err)
		return
	}

	return
}
