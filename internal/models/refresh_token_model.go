package models

import (
	"database/sql"
	"erp/internal/vos"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID    uuid.UUID     `db:"user_id"`
	FamilyID  uuid.UUID     `db:"family_id"`
	TokenHash vos.TokenHash `db:"token_hash"`
	ExpiresAt time.Time     `db:"expires_at"`
	RevokedAt sql.NullTime  `db:"revoked_at"`
}

func NewRefreshToken(userID uuid.UUID, familyID uuid.UUID, token string, expiresAt time.Time) RefreshToken {
	return RefreshToken{
		UserID:    userID,
		FamilyID:  familyID,
		TokenHash: vos.NewTokenHash(token),
		ExpiresAt: expiresAt,
	}
}
