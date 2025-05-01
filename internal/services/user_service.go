package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
)

// UserService handles user-related business logic
type UserService interface {
	Register(ctx context.Context, username, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error) // Returns JWT token
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, updateData map[string]interface{}) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	PromoteToAdmin(ctx context.Context, userID uuid.UUID) error
}

type userService struct {
	userRepo repositories.UserRepository
}
