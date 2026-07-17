package vos

import (
	"errors"
	"regexp"
	"strings"
)

type Email string

var (
	ErrEmailTooLong = errors.New("email must not exceed 254 characters")
	ErrEmailInvalid = errors.New("email is not a valid format")
	emailPattern    = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

func NewEmail(raw string) (email Email, err error) {
	switch {
	case len(raw) > 254:
		err = ErrEmailTooLong
	case !emailPattern.MatchString(raw):
		err = ErrEmailInvalid
	default:
		email = Email(strings.ToLower(raw))
	}
	return
}
