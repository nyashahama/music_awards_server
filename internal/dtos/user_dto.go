// Package dtos
package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
)

// RegisterRequest is the request payload for user registration.
type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Location  string `json:"location" binding:"required,min=2,max=100"`
}

// LoginRequest is the request payload for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type PurchaseVotesRequest struct {
	Amount int `json:"amount" binding:"required,min=1,max=1000"`
}

// UpdateProfileRequest is the request payload for updating a user's profile.
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
	Location  *string `json:"location"`
}

// UserResponse is the response payload for user details.
type UserResponse struct {
	UserID     uuid.UUID `json:"user_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	Location   string    `json:"location"`
	FreeVotes  int       `json:"free_votes"`
	PaidVotes  int       `json:"paid_votes"`
	TotalVotes int       `json:"total_votes"`
	CreatedAt  time.Time `json:"created_at"`
}

// LoginResponse is the response payload after a successful login.
type LoginResponse struct {
	Token string `json:"token"`
}

// NewUserResponse converts a models.User to a UserResponse DTO.
func NewUserResponse(user *models.User) UserResponse {
	return UserResponse{
		UserID:     user.UserID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Role:       user.Role,
		Location:   user.Location,
		FreeVotes:  user.FreeVotes,
		PaidVotes:  user.PaidVotes,
		TotalVotes: user.FreeVotes + user.PaidVotes,
		CreatedAt:  user.CreatedAt,
	}
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ValidateResetTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type PasswordResetResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"` // Only for development, remove in production
}
