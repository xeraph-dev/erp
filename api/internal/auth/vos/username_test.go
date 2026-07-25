package vos_test

import (
	"erp/internal/auth/vos"
	"errors"
	"strings"
	"testing"
)

func TestNewUsername(t *testing.T) {
	for _, tt := range []struct {
		name    string
		raw     string
		wantErr error
	}{
		{"valid username", "admin", nil},
		{"too short", "ad", vos.ErrUsernameTooShort},
		{"too long", strings.Repeat("a", 33), vos.ErrUsernameTooLong},
		{"invalid characters", "admin!", vos.ErrUsernameInvalid},
	} {
		t.Run(tt.name, func(t *testing.T) {
			username, err := vos.NewUsername(tt.raw)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got err %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && username.String() != tt.raw {
				t.Fatalf("got username %q, want %q", username, tt.raw)
			}
		})
	}
}
