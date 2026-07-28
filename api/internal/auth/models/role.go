package models

import (
	"erp/internal/auth/vos"

	"github.com/google/uuid"
)

type Role struct {
	ID    uuid.UUID     `db:"id"`
	Name  vos.RoleName  `db:"role_name"`
	Level vos.RoleLevel `db:"role_level"`
}
