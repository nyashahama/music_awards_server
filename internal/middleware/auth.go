package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyashahama/music-awards/internal/security"
)

// Middleware interface
type AuthMiddleware interface {
	Authenticate(next http.Handler) http.Handler
	AdminOnly(next http.Handler) http.Handler
	RequestLogger(next http.Handler) http.Handler
	ValidateVotingPeriod(next http.Handler) http.Handler
}

type authMiddleware struct {
	// Dependencies
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) < 7 || token[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid auth header"})
			return
		}
		claims, err := security.ValidateJWT(token[7:])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
