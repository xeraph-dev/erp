package models

import "github.com/google/uuid"

type Role struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"role_name"`
	Level int16     `db:"role_level"`
}
