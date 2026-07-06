package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID      `db:"id"`
	Username  string         `db:"username"`
	Password  string         `db:"password"`
	Email     string         `db:"email"`
	FirstName sql.NullString `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
}
