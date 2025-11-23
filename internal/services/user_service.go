// Package services
package services

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	Register(ctx context.Context, firstName, lastName, email, password, location string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error) // Returns JWT token
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, updateData map[string]any) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	PromoteToAdmin(ctx context.Context, userID uuid.UUID) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	RequestPasswordReset(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
	ValidateResetToken(ctx context.Context, token string) (*models.User, error)
}

type userService struct {
	userRepo             repositories.UserRepository
	passwordResetService PasswordResetService
	emailService         EmailService
}

func NewUserService(userRepo repositories.UserRepository, passwordResetService PasswordResetService, emailService EmailService) UserService {
	return &userService{
		userRepo:             userRepo,
		passwordResetService: passwordResetService,
		emailService:         emailService,
	}
}

func (s *userService) Register(ctx context.Context, firstName, lastName, email, password, location string) (*models.User, error) {
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
		return nil, fmt.Errorf("failed to hash the password: %w", err)
	}

	user := &models.User{
		UserID:       uuid.New(),
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: hashedPassword,
		FreeVotes:    3,
		Role:         "user",
		Location:     location,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Send welcome email (async - don't block registration)
	go func() {
		if err := s.emailService.SendWelcomeEmail(user.Email, user.FirstName, user.FreeVotes); err != nil {
			log.Printf("Failed to send welcome email to %s: %v", user.Email, err)
		}
	}()
	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(email)

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}

	// Use first name + last name for display name in JWT
	displayName := user.FirstName + " " + user.LastName
	token, err := security.GenerateJWT(user.UserID, displayName, user.Role, user.Email)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	userAgent := getStringFromContext(ctx, "User-Agent")
	ipAddress := getStringFromContext(ctx, "IP-Address")

	// Send login notification (async)
	go func() {
		if err := s.emailService.SendLoginNotificationEmail(user.Email, user.FirstName, userAgent, ipAddress); err != nil {
			log.Printf("Failed to send login notification to %s: %v", user.Email, err)
		}
	}()

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

func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, updateData map[string]any) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidID
	}

	if firstName, ok := updateData["first_name"].(string); ok {
		user.FirstName = firstName
	}

	if lastName, ok := updateData["last_name"].(string); ok {
		user.LastName = lastName
	}

	if location, ok := updateData["location"].(string); ok {
		user.Location = location
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

func (s *userService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	return s.passwordResetService.RequestPasswordReset(ctx, email)
}

func (s *userService) ResetPassword(ctx context.Context, token, newPassword string) error {
	return s.passwordResetService.ResetPassword(ctx, token, newPassword)
}

func (s *userService) ValidateResetToken(ctx context.Context, token string) (*models.User, error) {
	return s.passwordResetService.ValidateResetToken(ctx, token)
}

func getStringFromContext(ctx context.Context, key string) string {
	if value, ok := ctx.Value(key).(string); ok {
		return value
	}
	return "Unknown"
}
