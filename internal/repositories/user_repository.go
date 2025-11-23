// Package repositories
package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateEntry = errors.New("duplicate entry")
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Vote management
	DecrementFreeVotes(ctx context.Context, userID uuid.UUID) error
	DecrementPaidVotes(ctx context.Context, userID uuid.UUID) error
	IncrementFreeVotes(ctx context.Context, userID uuid.UUID) error
	IncrementPaidVotes(ctx context.Context, userID uuid.UUID) error
	AddPaidVotes(ctx context.Context, userID uuid.UUID, amount int) error

	// Password reset
	SetPasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	GetByResetToken(ctx context.Context, token string) (*models.User, error)
	ClearPasswordResetToken(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrDuplicateEntry
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "user_id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &user, err
}

func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "user_id = ?", id).Error
}

func (r *userRepository) DecrementFreeVotes(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&user, "user_id = ?", userID).Error; err != nil {
			return err
		}

		if user.FreeVotes <= 0 {
			return errors.New("no free votes available")
		}

		return tx.Model(&user).Update("free_votes", gorm.Expr("free_votes - 1")).Error
	})
}

func (r *userRepository) DecrementPaidVotes(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&user, "user_id = ?", userID).Error; err != nil {
			return err
		}

		if user.PaidVotes <= 0 {
			return errors.New("no paid votes available")
		}

		return tx.Model(&user).Update("paid_votes", gorm.Expr("paid_votes - 1")).Error
	})
}

func (r *userRepository) IncrementFreeVotes(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("free_votes", gorm.Expr("free_votes + 1")).Error
}

func (r *userRepository) IncrementPaidVotes(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("paid_votes", gorm.Expr("paid_votes + 1")).Error
}

func (r *userRepository) AddPaidVotes(ctx context.Context, userID uuid.UUID, amount int) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("paid_votes", gorm.Expr("paid_votes + ?", amount)).Error
}

// Password reset methods remain the same...
func (r *userRepository) SetPasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"reset_token":            token,
			"reset_token_expires_at": expiresAt,
		}).Error
}

func (r *userRepository) GetByResetToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("reset_token = ? AND reset_token_expires_at > ?", token, time.Now()).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &user, err
}

func (r *userRepository) ClearPasswordResetToken(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"reset_token":            nil,
			"reset_token_expires_at": nil,
		}).Error
}
