package validation

import (
	"strconv"
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
		name      string
		password  string
		wantErr   bool
		wantSteps int
	}{
		{"Valid password", "Pass123!", false, 0},
		{"Too short", "A1!", true, 3},
		{"Missing types", "password", true, 2},
		{"Repeating chars", "aaaBBB111", true, 3},
		{"Long password with repeats", "aaaaaaaaaaaaaaaaaaaaa", true, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				var steps int
				if err != nil {
					steps, _ = strconv.Atoi(err.Error()[len(err.Error())-1:])
				}
				if steps != tt.wantSteps {
					t.Errorf("Expected %d steps, got %d", tt.wantSteps, steps)
				}
			}
		})
	}
}
