package vos

import (
	"crypto/sha256"
	"encoding/hex"
)

type TokenHash string

func NewTokenHash(raw string) TokenHash {
	sum := sha256.Sum256([]byte(raw))
	hash := hex.EncodeToString(sum[:])
	return TokenHash(hash)
}
