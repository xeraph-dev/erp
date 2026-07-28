package vos

import (
	"errors"
	"regexp"
)

type RoleName string

var roleNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

var (
	ErrRoleNameTooShort = errors.New("role name must be at least 3 characters")
	ErrRoleNameTooLong  = errors.New("role name must not exceed 32 characters")
	ErrRoleNameInvalid  = errors.New("role name must contain only alphanumeric characters, underscore, period and dash")
)

func NewRoleName(raw string) (name RoleName, err error) {
	switch {
	case len(raw) < 3:
		return name, ErrRoleNameTooShort
	case len(raw) > 32:
		return name, ErrRoleNameTooLong
	case !roleNamePattern.MatchString(raw):
		return name, ErrRoleNameInvalid
	default:
		return RoleName(raw), nil
	}
}

func (name RoleName) String() string {
	return string(name)
}
