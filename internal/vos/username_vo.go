package vos

import (
	"errors"
	"regexp"
)

type Username string

var usernamePattern = regexp.MustCompile(`[a-zA-Z0-9\_\.\-]`)

var (
	ErrUsernameTooShort = errors.New("username must be at least 3 characters")
	ErrUsernameTooLong  = errors.New("username must not expceed 32 characters")
	ErrUsernameInvalid  = errors.New("username must contain only alphanumeric, underscore, period and dash characters")
)

func NewUsername(raw string) (username Username, err error) {
	switch {
	case len(raw) < 3:
		err = ErrUsernameTooShort
	case len(raw) > 32:
		err = ErrUsernameTooLong
	case !usernamePattern.MatchString(raw):
		err = ErrUsernameInvalid
	default:
		username = Username(raw)
	}
	return
}
