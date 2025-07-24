package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	originalValidate := security.ValidateJWT
	defer func() { security.ValidateJWT = originalValidate }()

	security.ValidateJWT = func(token string) (*security.Claims, error) {
		if token == "valid" {
			return &security.Claims{UserID: "123"}, nil
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
