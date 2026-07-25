package vos_test

import (
	"erp/internal/auth/vos"
	"errors"
	"strings"
	"testing"
)

func TestNewEmail(t *testing.T) {
	for _, tt := range []struct {
		name    string
		raw     string
		want    string
		wantErr error
	}{
		{"valid email", "admin@erp.com", "admin@erp.com", nil},
		{"normilizes to lowercase", "Admin@ERP.COM", "admin@erp.com", nil},
		{"missing @", "adminerp.com", "", vos.ErrEmailInvalid},
		{"missing domain", "admin@", "", vos.ErrEmailInvalid},
		{"too long", strings.Repeat("a", 250) + "@erp.com", "", vos.ErrEmailTooLong},
	} {
		t.Run(tt.name, func(t *testing.T) {
			email, err := vos.NewEmail(tt.raw)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got err %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && email.String() != tt.want {
				t.Fatalf("got email %q, want %q", email, tt.want)
			}
		})
	}
}
