package handlers

import "net/http"

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