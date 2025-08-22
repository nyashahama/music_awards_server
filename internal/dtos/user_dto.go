// Package dtos
package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
)

// RegisterRequest is the request payload for user registration.
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest is the request payload for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest is the request payload for updating a user's profile.
type UpdateProfileRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

// UserResponse is the response payload for user details.
type UserResponse struct {
	UserID         uuid.UUID `json:"user_id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	AvailableVotes int       `json:"available_votes"`
	CreatedAt      time.Time `json:"created_at"`
	// UpdatedAt      time.Time `json:"updated_at"`
}

// LoginResponse is the response payload after a successful login.
type LoginResponse struct {
	Token string `json:"token"`
}

// NewUserResponse converts a models.User to a UserResponse DTO.
func NewUserResponse(user *models.User) UserResponse {
	return UserResponse{
		UserID:         user.UserID,
		Username:       user.Username,
		Email:          user.Email,
		Role:           user.Role,
		AvailableVotes: user.AvailableVotes,
		CreatedAt:      user.CreatedAt,
		// UpdatedAt:      user.UpdatedAt,
	}
}
