package services

import (
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
)

type UserService interface {
	RegisterUser(user *models.User) error
	Authenticate(email,password string) (*models.User, error)

}

type userService struct {
	userRepo repositories.UserRepository
}