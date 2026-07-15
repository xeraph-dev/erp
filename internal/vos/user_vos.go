package vos

import "golang.org/x/crypto/bcrypt"

type PasswordHash string

func NewPasswordHash(raw string) (hash PasswordHash, err error) {
	hash_bytes, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	hash = PasswordHash(hash_bytes)
	return
}

func (hash PasswordHash) Matches(raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)) != nil
}
