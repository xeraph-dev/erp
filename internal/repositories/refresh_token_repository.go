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

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type RefreshTokenRepository interface {
	Repository
	Create(ctx context.Context, db Querier, in models.RefreshToken) (out models.RefreshToken, err error)
	GetByTokenHash(ctx context.Context, db Querier, tokenHash vos.TokenHash) (out models.RefreshToken, err error)
	Revoke(ctx context.Context, db Querier, in models.RefreshToken) (out models.RefreshToken, err error)
}

type refreshTokenRepositoryImpl struct{}

var _ RefreshTokenRepository = (*refreshTokenRepositoryImpl)(nil)

func NewRefreshTokenRepository() RefreshTokenRepository {
	return refreshTokenRepositoryImpl{}
}

func (refreshTokenRepositoryImpl) __internal() {}

func (repo refreshTokenRepositoryImpl) Create(ctx context.Context, db Querier, in models.RefreshToken) (out models.RefreshToken, err error) {
	rows, err := db.Query(ctx, queries.CreateRefreshToken, in.UserID, in.FamilyID, in.TokenHash, in.ExpiresAt)
	if err != nil {
		err = fmt.Errorf("creating refresh token entry: %w", err)
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.RefreshToken])
	if err != nil {
		err = fmt.Errorf("collecting refresh token row: %w", err)
		return
	}

	return
}

func (repo refreshTokenRepositoryImpl) GetByTokenHash(ctx context.Context, db Querier, tokenHash vos.TokenHash) (out models.RefreshToken, err error) {
	rows, err := db.Query(ctx, queries.GetRefreshTokenByTokenHash, tokenHash)
	if err != nil {
		err = fmt.Errorf("querying refresh token family ID: %w", err)
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.RefreshToken])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = ErrRefreshTokenNotFound
			return
		}
		err = fmt.Errorf("collecting refresh token row: %w", err)
		return
	}

	return
}

func (repo refreshTokenRepositoryImpl) Revoke(ctx context.Context, db Querier, in models.RefreshToken) (out models.RefreshToken, err error) {
	rows, err := db.Query(ctx, queries.RevokeRefreshToken, in.TokenHash)
	if err != nil {
		err = fmt.Errorf("updating refresh token: %w", err)
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.RefreshToken])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = ErrRefreshTokenNotFound
			return
		}
		err = fmt.Errorf("collecting refresh token row: %w", err)
	}

	return
}
