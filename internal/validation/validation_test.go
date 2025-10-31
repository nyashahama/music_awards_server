package validation

import (
	"strings"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		want  bool
	}{
		{"test@example.com", true},
		{"invalid.email", false},
		{"name@domain", false},
		{"UPPERCASE@domain.com", false},
		{"test@sub.domain.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			got := ValidateEmail(tt.email)
			if got != tt.want {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string // expect specific error message content
	}{
		{"Valid password", "Pass123!", false, ""},
		{"Too short", "A1!", true, "at least 8 characters"},
		{"Common password", "password", true, "too common"},
		// ... etc
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing %q, got %v", tt.errMsg, err)
				}
			}
		})
	}
}
