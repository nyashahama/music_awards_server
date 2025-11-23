package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
)

// RegisterRequest is the request payload for user registration.
type RegisterRequest struct {
<<<<<<< Updated upstream
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
=======
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Location  string `json:"location" binding:"required,min=2,max=100"`
>>>>>>> Stashed changes
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
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

// UserResponse is the response payload for user details.
type UserResponse struct {
<<<<<<< Updated upstream
	UserID         uuid.UUID `json:"user_id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	AvailableVotes int       `json:"available_votes"`
	CreatedAt      time.Time `json:"created_at"`
	//UpdatedAt      time.Time `json:"updated_at"`
=======
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
>>>>>>> Stashed changes
}

// LoginResponse is the response payload after a successful login.
type LoginResponse struct {
	Token string `json:"token"`
}

// NewUserResponse converts a models.User to a UserResponse DTO.
func NewUserResponse(user *models.User) UserResponse {
	return UserResponse{
<<<<<<< Updated upstream
		UserID:         user.UserID,
		Username:       user.Username,
		Email:          user.Email,
		Role:           user.Role,
		AvailableVotes: user.AvailableVotes,
		CreatedAt:      user.CreatedAt,
		//UpdatedAt:      user.UpdatedAt,
=======
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
>>>>>>> Stashed changes
	}
}
