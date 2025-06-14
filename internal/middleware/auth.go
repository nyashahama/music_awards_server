package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/security"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		log.Printf("Received token: %s", tokenString) // Debug

		claims, err := security.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("Token validation failed: %v", err) // Debug
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			return
		}

		log.Printf("Token claims: %+v", claims) // Debug

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			log.Printf("UUID parsing failed: %v", err) // Debug
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
			return
		}

		log.Printf("Authenticated userID: %s", userID) // Debug

		c.Set("user_id", userID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)
		c.Set("email", claims.Email)
		c.Next()
	}
}
