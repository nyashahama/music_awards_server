package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
	"github.com/nyashahama/music-awards/internal/security"
	"github.com/nyashahama/music-awards/internal/validation"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService handles user-related business logic
type UserService interface {
	Register(ctx context.Context, username, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error) // Returns JWT token
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, updateData map[string]interface{}) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	PromoteToAdmin(ctx context.Context, userID uuid.UUID) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	 // 1) Validate email format
	 if !validation.ValidateEmail(email) {
        return nil, errors.New("invalid email format")
    }

    // 2) Validate password strength
    if err := validation.ValidatePassword(password); err != nil {
        return nil, fmt.Errorf("password validation failed: %w", err)
    }

	// Check for existing email
	_, err := s.userRepo.GetByEmail(email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	// if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, fmt.Errorf("failed to check email: %w", err)
	// }

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		UserID:       uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid credentials")
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}
	token, err := security.GenerateJWT(user.UserID,user.Username, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return token, nil
}

func (s *userService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, updateData map[string]interface{}) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	if username, ok := updateData["username"].(string); ok {
		user.Username = username
	} else if updateData["username"] != nil {
		return nil, errors.New("invalid type for username")
	}

	if email, ok := updateData["email"].(string); ok {
		existingUser, err := s.userRepo.GetByEmail(email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check email: %v", err)
		}
		if existingUser != nil && existingUser.UserID != user.UserID {
			return nil, errors.New("email already in use")
		}
		user.Email = email
	} else if updateData["email"] != nil {
		return nil, errors.New("invalid type for email")
	}

	if password, ok := updateData["password"].(string); ok {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
		user.PasswordHash = hashedPassword
	} else if updateData["password"] != nil {
		return nil, errors.New("invalid type for password")
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepo.Delete(userID); err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

func (s *userService) PromoteToAdmin(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	user.Role = "admin"
	return s.userRepo.Update(user)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.User, error){
	users, err := s.userRepo.GetAll()
	if err != nil{
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}
