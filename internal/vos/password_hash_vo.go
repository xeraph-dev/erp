package vos

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHash string

var (
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong    = errors.New("password must not exceed 72 characters")
	ErrPasswordWeak       = errors.New("password must has at least one lowercase character, one uppercase character, one digit, and one special symbol")
)

func NewPasswordHash(raw string) (hash PasswordHash, err error) {
	var lower bool
	var upper bool
	var digit bool
	var special bool

	for _, ch := range raw {
		switch {
		case ch >= 'a' && ch <= 'z':
			lower = true
		case ch >= 'A' && ch <= 'Z':
			upper = true
		case ch >= '0' && ch <= '9':
			digit = true
		case ch == '_' || ch == '.' || ch == '-':
			special = true
		}
	}

	switch {
	case len(raw) < 8:
		err = ErrPasswordTooShort
	case len(raw) > 72:
		err = ErrPasswordTooLong
	case !(lower && upper && digit && special):
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
