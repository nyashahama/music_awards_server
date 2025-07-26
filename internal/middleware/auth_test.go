package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nyashahama/music-awards/internal/security"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Save original function and restore after test
	originalValidate := security.ValidateJWT
	defer func() { security.ValidateJWT = originalValidate }()

	// Override with mock
	security.ValidateJWT = func(token string) (*security.JWTClaims, error) {
		if token == "valid" {
			return &security.JWTClaims{
				UserID:   "123e4567-e89b-12d3-a456-426614174000", // Valid UUID format
				Username: "testuser",
				Role:     "user",
				Email:    "test@example.com",
			}, nil
		}
		return nil, errors.New("invalid token")
	}

	tests := []struct {
		name     string
		token    string
		wantCode int
	}{
		{"Valid token", "Bearer valid", http.StatusOK},
		{"Invalid token", "Bearer invalid", http.StatusUnauthorized},
		{"Missing token", "", http.StatusUnauthorized},
		{"Malformed header", "InvalidFormat", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(AuthMiddleware())
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
