package stores

import (
	_ "embed"

	"context"
	"erp/internal/auth/models"
	"erp/internal/auth/vos"
	"erp/internal/shared/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RefreshToken interface {
	Create(ctx context.Context, q db.Querier, in models.RefreshToken) (out models.RefreshToken, err error)
	GetByHash(ctx context.Context, q db.Querier, hash vos.TokenHash) (out models.RefreshToken, err error)
	GetByFamilyID(ctx context.Context, q db.Querier, id uuid.UUID) (out []models.RefreshToken, err error)
	RevokeByHash(ctx context.Context, q db.Querier, hash vos.TokenHash) (out models.RefreshToken, err error)
	RevokeByFamilyID(ctx context.Context, q db.Querier, id uuid.UUID) (out []models.RefreshToken, err error)
}

type refreshTokenImpl struct{}

var _ RefreshToken = (*refreshTokenImpl)(nil)

func NewRefreshToken() RefreshToken {
	return refreshTokenImpl{}
}

//go:embed queries/create_refresh_token.sql
var createRefreshToken string

func (store refreshTokenImpl) Create(ctx context.Context, q db.Querier, in models.RefreshToken) (out models.RefreshToken, err error) {
	rows, err := q.Query(ctx, createRefreshToken, in.TokenHash, in.UserID, in.FamilyID, in.ExpiresAt)
	if err != nil {
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.RefreshToken])
	if err != nil {
		return
	}

	return
}

//go:embed queries/get_refresh_token_by_hash.sql
var getRefreshTokenByHash string

func (store refreshTokenImpl) GetByHash(ctx context.Context, q db.Querier, hash vos.TokenHash) (out models.RefreshToken, err error) {
	rows, err := q.Query(ctx, getRefreshTokenByHash, hash)
	if err != nil {
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.RefreshToken])
	if err != nil {
		return
	}

	return
}

//go:embed queries/get_refresh_tokens_family_id.sql
var getRefreshTokensByFamilyID string

func (store refreshTokenImpl) GetByFamilyID(ctx context.Context, q db.Querier, id uuid.UUID) (out []models.RefreshToken, err error) {
	rows, err := q.Query(ctx, getRefreshTokensByFamilyID, id)
	if err != nil {
		return
	}

	out, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.RefreshToken])
	if err != nil {
		return
	}

	return
}

//go:embed queries/revoke_refresh_token_by_hash.sql
var revokeRefreshTokenByHash string

func (store refreshTokenImpl) RevokeByHash(ctx context.Context, q db.Querier, hash vos.TokenHash) (out models.RefreshToken, err error) {
	rows, err := q.Query(ctx, revokeRefreshTokenByHash, hash)
	if err != nil {
		return
	}

	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.RefreshToken])
	if err != nil {
		return
	}

	return
}

//go:embed queries/revoke_refresh_tokens_by_family_id.sql
var revokeRefreshTokensByFamilyID string

func (store refreshTokenImpl) RevokeByFamilyID(ctx context.Context, q db.Querier, id uuid.UUID) (out []models.RefreshToken, err error) {
	rows, err := q.Query(ctx, revokeRefreshTokensByFamilyID, id)
	if err != nil {
		return
	}

	out, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.RefreshToken])
	if err != nil {
		return
	}

	return
}
