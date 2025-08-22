package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid password", "securePassword123", false},
		{"Empty password", "", false},
		{"Long password", "veryLongPasswordWithSpecialChars!@#$%^&*()1234567890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, hash)
			assert.NotEqual(t, tt.password, hash)
		})
	}
}

func TestComparePassword(t *testing.T) {
	password := "testPassword123"
	hash, err := HashPassword(password)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		hash     string
		password string
		wantErr  bool
	}{
		{"Correct password", hash, password, false},
		{"Wrong password", hash, "wrongPassword", true},
		{"Empty hash", "", password, true},
		{"Empty password", hash, "", true},
		{"Invalid hash", "invalidHash", password, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ComparePassword(tt.hash, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHashPassword_Consistency(t *testing.T) {
	password := "samePassword"
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	// Hashes should be different due to different salts
	assert.NotEqual(t, hash1, hash2)

	// But both should validate correctly
	err1 = ComparePassword(hash1, password)
	err2 = ComparePassword(hash2, password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
}
