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
		return email, ErrEmailTooLong
	case !emailPattern.MatchString(raw):
		return email, ErrEmailInvalid
	default:
		return Email(strings.ToLower(raw)), nil
	}
}

func (email Email) String() string {
	return string(email)
}
