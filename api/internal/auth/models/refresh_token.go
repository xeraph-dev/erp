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
