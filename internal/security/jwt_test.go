package security

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestEnv() {
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-12345")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
}

func TestGenerateAndValidateJWT(t *testing.T) {
	setupTestEnv()

	userID := uuid.New()
	username := "Nyashaa"
	role := "admin"
	email := "nyashahama45@gmail.com"

	token, err := GenerateJWT(userID, username, role, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, userID.String(), claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, email, claims.Email)

	assert.WithinDuration(t, time.Now().Add(24*time.Hour), claims.ExpiresAt.Time, time.Minute)
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	setupTestEnv()

	tests := []struct {
		name      string
		token     string
		wantError string
	}{
		{"Empty token", "", "token parse error"},
		{"Malformed token", "invalid.token.here", "token parse error"},
		{"Invalid signature", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", "signature is invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateJWT(tt.token)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantError)
		})
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	setupTestEnv()

	userID := uuid.New()
	claims := JWTClaims{
		UserID:   userID.String(),
		Username: "testuser",
		Role:     "user",
		Email:    "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	assert.NoError(t, err)

	_, err = ValidateJWT(tokenString)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestValidateJWT_DifferentSecret(t *testing.T) {
	// Generate token with one secret
	os.Setenv("JWT_SECRET", "first-secret-12345")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	userID := uuid.New()
	token, err := GenerateJWT(userID, "testuser", "user", "test@example.com")
	assert.NoError(t, err)

	// Try to validate with different secret
	os.Setenv("JWT_SECRET", "different-secret-67890")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	_, err = ValidateJWT(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature is invalid")
}

func TestJWTClaims_Structure(t *testing.T) {
	setupTestEnv()

	userID := uuid.New()
	token, err := GenerateJWT(userID, "testuser", "admin", "test@example.com")
	assert.NoError(t, err)

	claims, err := ValidateJWT(token)
	assert.NoError(t, err)

	assert.Equal(t, userID.String(), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "admin", claims.Role)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)

	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func TestValidateJWT_EmptySecret(t *testing.T) {
	// Save original secret
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	// Set empty secret
	os.Unsetenv("JWT_SECRET")
	jwtSecret = []byte("")

	// Create a valid token with a non-empty secret first
	os.Setenv("JWT_SECRET", "valid-secret")
	jwtSecret = []byte("valid-secret")
	userID := uuid.New()
	token, err := GenerateJWT(userID, "testuser", "user", "test@example.com")
	assert.NoError(t, err)

	// Now set empty secret and try to validate
	os.Unsetenv("JWT_SECRET")
	jwtSecret = []byte("")

	_, err = ValidateJWT(token)
	assert.Error(t, err)
	// The error could be either "signature is invalid" or "malformed token"
	// depending on how the JWT library handles empty secrets
}

func TestGenerateJWT_WithEmptySecret(t *testing.T) {
	// Save original secret
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	// Set empty secret for generation
	os.Unsetenv("JWT_SECRET")
	jwtSecret = []byte("")

	userID := uuid.New()

	// This should not error even with empty secret
	token, err := GenerateJWT(userID, "testuser", "user", "test@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Restore a valid secret for validation
	os.Setenv("JWT_SECRET", "valid-secret-for-validation")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	// Now validation should fail because the token was signed with a different (empty) secret
	_, err = ValidateJWT(token)
	assert.Error(t, err)
	// The error could be about invalid signature
}

func TestJWT_TokenFormat(t *testing.T) {
	setupTestEnv()

	fixedUUID, _ := uuid.Parse("12470f7b-f5ae-431c-b2fc-81d7147614f6")

	token, err := GenerateJWT(fixedUUID, "Nyashaa", "admin", "nyashahama45@gmail.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token has the expected format (3 parts separated by dots)
	parts := 0
	for i := 0; i < len(token); i++ {
		if token[i] == '.' {
			parts++
		}
	}
	assert.Equal(t, 2, parts, "JWT should have exactly 2 dots separating 3 parts")

	claims, err := ValidateJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, fixedUUID.String(), claims.UserID)
	assert.Equal(t, "Nyashaa", claims.Username)
	assert.Equal(t, "admin", claims.Role)
	assert.Equal(t, "nyashahama45@gmail.com", claims.Email)
}
