package vos_test

import (
	"erp/internal/auth/vos"
	"testing"
)

func TestNewTokenHash(t *testing.T) {
	hash1 := vos.NewTokenHash("some-refresh-token")
	hash2 := vos.NewTokenHash("some-refresh-token")
	hash3 := vos.NewTokenHash("different-token")

	if hash1 != hash2 {
		t.Fatal("same input must produce the same hash")
	}
	if hash1 == hash3 {
		t.Fatal("different input must produce different hashes")
	}
	if len(hash1.String()) != 64 {
		t.Fatalf("expected 64-char hex hash, got %d chars", len(hash1.String()))
	}
}
