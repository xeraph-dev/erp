package models

import (
	"database/sql"
)

type User struct {
	base

	Name      string         `db:"name"`
	Password  string         `db:"password"`
	Email     string         `db:"email"`
	FirstName sql.NullString `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
}
