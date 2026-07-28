package vos

import (
	"errors"
	"regexp"
)

type Username string

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

var (
	ErrUsernameTooShort = errors.New("username must be at least 3 characters")
	ErrUsernameTooLong  = errors.New("username must not exceed 32 characters")
	ErrUsernameInvalid  = errors.New("username must contain only alphanumeric characters, underscore, period and dash")
)

func NewUsername(raw string) (username Username, err error) {
	switch {
	case len(raw) < 3:
		return username, ErrUsernameTooShort
	case len(raw) > 32:
		return username, ErrUsernameTooLong
	case !usernamePattern.MatchString(raw):
		return username, ErrUsernameInvalid
	default:
		return Username(raw), nil
	}
}

func (username Username) String() string {
	return string(username)
}
