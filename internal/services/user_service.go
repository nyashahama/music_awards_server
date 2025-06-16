package services

import (
	"context"
	"errors"
	"fmt"
	//	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
	"github.com/nyashahama/music-awards/internal/security"
	"github.com/nyashahama/music-awards/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordValidation = errors.New("password validation failed")
	ErrInvalidID          = errors.New("invalid id")
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
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(ctx context.Context, username, email, password string) (*models.User, error) {

	email = strings.ToLower(email)

	if !validation.ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	if err := validation.ValidatePassword(password); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrPasswordValidation, err)
	}

	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailExists
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		UserID:         uuid.New(),
		Username:       username,
		Email:          email,
		PasswordHash:   hashedPassword,
		Role:           "user",
		AvailableVotes: 5,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(email)
	/* log.Printf("Login attempt for: %s", email) */

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		/* log.Printf("Login error: %v", err) */
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		/* log.Printf("User not found: %s", email) */
		return "", ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		/* log.Printf("Invalid password for: %s", email) */
		return "", ErrInvalidCredentials
	}

	/* log.Printf("Login successful for user: %s (%s)", user.Username, user.UserID) */

	token, err := security.GenerateJWT(user.UserID, user.Username, user.Role, user.Email)
	if err != nil {
		/* log.Printf("Token generation failed: %v", err) */
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	/* log.Printf("Generated token: %s", token) */
	return token, nil
}
func (s *userService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidID
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, updateData map[string]interface{}) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidID
	}

	if username, ok := updateData["username"].(string); ok {
		user.Username = username
	}

	if email, ok := updateData["email"].(string); ok {
		// Add validation
		email = strings.ToLower(email)
		if !validation.ValidateEmail(email) {
			return nil, fmt.Errorf("invalid email format")
		}
		existing, err := s.userRepo.GetByEmail(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if existing != nil && existing.UserID != user.UserID {
			return nil, ErrEmailExists
		}
		user.Email = email
	}

	if password, ok := updateData["password"].(string); ok {
		// Add validation
		if err := validation.ValidatePassword(password); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrPasswordValidation, err)
		}
		hashedPassword, err := hashPassword(password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = hashedPassword
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Let repository handle not-found vs. internal errors; you may choose to wrap ErrInvalidID as needed.
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *userService) PromoteToAdmin(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return ErrInvalidID
	}

	user.Role = "admin"
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to promote user: %w", err)
	}
	return nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAll(ctx)
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}
