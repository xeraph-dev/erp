package models

import (
	"database/sql"
	"erp/internal/auth/vos"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID        `db:"id"`
	Username     vos.Username     `db:"username"`
	Email        vos.Email        `db:"email"`
	PasswordHash vos.PasswordHash `db:"password_hash"`
	FirstName    sql.NullString   `db:"first_name"`
	LastName     sql.NullString   `db:"last_name"`
}

func (user User) PasswordMatches(raw string) bool {
	return user.PasswordHash.Matches(raw)
}
