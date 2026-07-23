package vos

import (
	"errors"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHash string

var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong  = errors.New("password must not exceed 72 characters")
	ErrPasswordWeak     = errors.New("password must has at least one lowercase character, one uppercase character, one digit, and one special symbol")
)

func NewPasswordHash(raw string) (hash PasswordHash, err error) {
	switch {
	case len(raw) < 8:
		err = ErrPasswordTooShort
	case len(raw) > 72:
		err = ErrPasswordTooLong
	case !strings.ContainsFunc(raw, unicode.IsLower) ||
		!strings.ContainsFunc(raw, unicode.IsUpper) ||
		!strings.ContainsFunc(raw, unicode.IsDigit) ||
		!strings.ContainsFunc(raw, unicode.IsSymbol):
		err = ErrPasswordWeak
	default:
		var password []byte
		password, err = bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
		if err != nil {
			return
		} else {
			hash = PasswordHash(password)
		}
	}

	return
}

func (hash PasswordHash) Matches(raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)) == nil
}
