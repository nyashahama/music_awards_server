package middleware

import (
	"errors"

	"github.com/nyashahama/music-awards/internal/security"
)

var mockValidateJWT = func(token string) (*security.JWTClaims, error) {
	if token == "valid" {
		return &security.JWTClaims{UserID: "123"}, nil

	}
	return nil, errors.New("invalid token")
}
