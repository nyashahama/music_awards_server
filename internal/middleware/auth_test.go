package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nyashahama/music-awards/internal/security"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// TODO: complte this:
	originalValidate := security.ValidateJWT()
	security.ValidateJWT = mockValidateJWT
	defer func() { security.ValidateJWT = originalValidate }()

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
