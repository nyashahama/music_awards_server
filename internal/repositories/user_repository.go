package repositories

import (
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User,error)
}

type userRepository struct {
	db *gorm.DB
}