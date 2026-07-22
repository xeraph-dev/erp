package vos

import "errors"

type RoleLevel int16

var (
	ErrRoleLevelTooSmall = errors.New("role level should be greater than 0")
	ErrRoleLevelTooLarge = errors.New("role level should be lower than 32767")
)

func NewRoleLevel(raw int16) (level RoleLevel, err error) {
	switch {
		case raw < 0:
			err = ErrRoleLevelTooSmall
			return
		case raw > 32767:
			err = ErrRoleLevelTooLarge
			return
		default:
			level = RoleLevel(raw)
	}
	return
}
