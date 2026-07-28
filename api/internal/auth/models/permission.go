package models

import "github.com/google/uuid"

type Permission struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"permission_name"`
}
