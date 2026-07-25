package vos_test

import (
	"erp/internal/auth/vos"
	"errors"
	"testing"
)

func TestNewPasswordHash(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{"valid with punctuation", "Passw0rd!", nil},
		{"valid with symbol", "Passw0rd$", nil},
		{"too short", "Pw0!", vos.ErrPasswordTooShort},
		{"no uppercase", "passw0rd!", vos.ErrPasswordWeak},
		{"no lowercase", "PASSW0RD!", vos.ErrPasswordWeak},
		{"no digit", "Password!", vos.ErrPasswordWeak},
		{"no special char", "Passw0rd", vos.ErrPasswordWeak},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := vos.NewPasswordHash(tt.raw)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got err %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && !hash.Matches(tt.raw) {
				t.Fatalf("hash does not match original password %q", tt.raw)
			}
		})
	}
}

func TestPasswordHash_Matches(t *testing.T) {
	hash, err := vos.NewPasswordHash("Passw0rd!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !hash.Matches("Passw0rd!") {
		t.Fatal("expected hash to match original password")
	}
	if hash.Matches("WrongPassword1!") {
		t.Fatal("expected hash to not match a different password")
	}
}
