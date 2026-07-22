package vos

import (
	"errors"
)

type RoleLevel int16

var (
	ErrRoleLevelTooSmall = errors.New("role level should be greater than 0")
)

func NewRoleLevel(raw int16) (level RoleLevel, err error) {
	switch {
	case raw < 0:
		err = ErrRoleLevelTooSmall
		return
	default:
		level = RoleLevel(raw)
	}
	return
}
