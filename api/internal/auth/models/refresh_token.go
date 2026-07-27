package models

import (
	"database/sql"
	"erp/internal/auth/vos"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	TokenHash vos.TokenHash `db:"token_hash"`
	UserID    uuid.UUID     `db:"user_id"`
	FamilyID  uuid.UUID     `db:"family_id"`
	ExpiresAt time.Time     `db:"expires_at"`
	RevokedAt sql.NullTime  `db:"revoked_at"`
}

func NewRefreshTokenInFamily(token string, userID, familyID uuid.UUID, expiresAt time.Time) RefreshToken {
	return RefreshToken{
		TokenHash: vos.NewTokenHash(token),
		UserID:    userID,
		FamilyID:  familyID,
		ExpiresAt: expiresAt,
	}
}

func NewRefreshToken(token string, userID uuid.UUID, expiresAt time.Time) RefreshToken {
	return NewRefreshTokenInFamily(token, userID, uuid.New(), expiresAt)
}
