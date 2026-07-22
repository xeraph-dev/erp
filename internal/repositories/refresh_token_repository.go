package repositories

import (
	"context"
	"erp/internal/middlewares"
	"erp/internal/models"
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

	rows, err := db.Exec(ctx, "INSERT INTO erp.refresh_tokens (user_id, family_id, token_hash, expires_at) VALUES ($1, $2, $3, $4)", in.UserID, in.FamilyID, in.TokenHash, in.ExpiresAt)
	if err != nil {
		logger.ErrorContext(ctx, "creating refresh token entry", "error", err)
		return
	}

	if rows.RowsAffected() != 0 {
		// TODO: here should be an error
		return
	}

	return
}
