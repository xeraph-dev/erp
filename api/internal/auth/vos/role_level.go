package vos

import (
	"errors"
)

type RoleLevel int16

var (
	ErrRoleLevelTooSmall = errors.New("role level must be greater than 0")
)

func NewRoleLevel(raw int16) (level RoleLevel, err error) {
	switch {
	case raw < 0:
		return level, ErrRoleLevelTooSmall
	default:
		return RoleLevel(raw), nil
	}
}
