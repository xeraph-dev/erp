package vos

import (
	"errors"
	"regexp"
)

type RoleName string

var roleNamePattern = regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]+$`)

var (
	ErrRoleNameTooShort = errors.New("role name must be at least 3 characters")
	ErrRoleNameTooLong  = errors.New("role name must not exceed 32 characters")
	ErrRoleNameInvalid  = errors.New("role name must contain only alphanumeric characteres, underscore, period and dash")
)

func NewRoleName(raw string) (name RoleName, err error) {
	switch {
	case len(raw) < 3:
		err = ErrRoleNameTooShort
	case len(raw) > 32:
		err = ErrRoleNameTooLong
	case !roleNamePattern.MatchString(raw):
		err = ErrRoleNameInvalid
	default:
		name = RoleName(raw)
	}
	return
}
