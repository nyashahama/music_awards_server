package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/gorm"
)

// Vote Repository
type VoteRepository interface {
	Create(ctx context.Context, vote *models.Vote) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Vote, error)
	GetAll(ctx context.Context) ([]models.Vote, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]models.Vote, error)
	GetByUserAndCategory(ctx context.Context, userID, categoryID uuid.UUID) (*models.Vote, error)
	Update(ctx context.Context, vote *models.Vote) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type voteRepository struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) VoteRepository {
	return &voteRepository{db: db}
}

func (r *voteRepository) Create(ctx context.Context, vote *models.Vote) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing models.Vote
		err := tx.WithContext(ctx).
			Where("user_id = ? AND category_id = ?", vote.UserID, vote.CategoryID).
			First(&existing).Error

		if err == nil {
			return fmt.Errorf("user already voted in this category")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		return tx.WithContext(ctx).Create(vote).Error
	})
}

func (r *voteRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.WithContext(ctx).First(&vote, "vote_id = ?", id).Error
	return &vote, err
}

func (r *voteRepository) GetAll(ctx context.Context) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.WithContext(ctx).Find(&votes).Error
	return votes, err
}

func (r *voteRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.WithContext(ctx).Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("category_id", "name")
	}).
		Preload("Nominee", func(db *gorm.DB) *gorm.DB {
			return db.Select("nominee_id", "name")
		}).
		Where("user_id = ?", userID).
		Find(&votes).Error
	return votes, err
}

func (r *voteRepository) GetByUserAndCategory(ctx context.Context, userID, categoryID uuid.UUID) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND category_id = ?", userID, categoryID).
		First(&vote).Error

	// Handle record not found case
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &vote, err
}

func (r *voteRepository) Update(ctx context.Context, vote *models.Vote) error {
	return r.db.WithContext(ctx).Save(vote).Error
}

func (r *voteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Vote{}, "vote_id = ?", id).Error
}
