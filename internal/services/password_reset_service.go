// Package services
package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
	"github.com/nyashahama/music-awards/internal/validation"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidToken    = errors.New("invalid or expired reset token")
	ErrTokenGeneration = errors.New("failed to generate reset token")
)

type PasswordResetService interface {
	RequestPasswordReset(ctx context.Context, email string) (string, error) // Returns reset token
	ResetPassword(ctx context.Context, token, newPassword string) error
	ValidateResetToken(ctx context.Context, token string) (*models.User, error)
}

type passwordResetService struct {
	userRepo     repositories.UserRepository
	emailService EmailService
}

func NewPasswordResetService(userRepo repositories.UserRepository, emailService EmailService) PasswordResetService {
	return &passwordResetService{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (s *passwordResetService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		// Don't reveal whether email exists or not for security
		return "", nil
	}

	// Generate secure random token
	token, err := generateSecureToken(32)
	if err != nil {
		return "", ErrTokenGeneration
	}

	// Set token expiration (1 hour from now)
	expiresAt := time.Now().Add(1 * time.Hour)

	if err := s.userRepo.SetPasswordResetToken(ctx, user.UserID, token, expiresAt); err != nil {
		return "", fmt.Errorf("failed to set reset token: %w", err)
	}

	// Send password reset email (async)
	go func() {
		if err := s.emailService.SendPasswordResetEmail(user.Email, user.FirstName, token); err != nil {
			log.Printf("Failed to send password reset email to %s: %v", user.Email, err)
		}
	}()

	return token, nil
}

func (s *passwordResetService) ValidateResetToken(ctx context.Context, token string) (*models.User, error) {
	if strings.TrimSpace(token) == "" {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByResetToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidToken
	}

	return user, nil
}

func (s *passwordResetService) ResetPassword(ctx context.Context, token, newPassword string) error {
	user, err := s.ValidateResetToken(ctx, token)
	if err != nil {
		return err
	}

	// Validate new password
	if err := validation.ValidatePassword(newPassword); err != nil {
		return fmt.Errorf("%w: %s", ErrPasswordValidation, err)
	}

	// Hash new password
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password and clear reset token in a transaction
	if err := s.userRepo.Update(ctx, &models.User{
		UserID:              user.UserID,
		PasswordHash:        hashedPassword,
		ResetToken:          nil,
		ResetTokenExpiresAt: nil,
	}); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
