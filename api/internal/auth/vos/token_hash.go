package vos

import (
	"crypto/sha256"
	"encoding/hex"
)

type TokenHash string

func NewTokenHash(raw string) TokenHash {
	sum := sha256.Sum256([]byte(raw))
	return TokenHash(hex.EncodeToString(sum[:]))
}

func (hash TokenHash) String() string {
	return string(hash)
}
