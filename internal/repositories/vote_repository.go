package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/gorm"
)

type VoteRepository interface {
	Create(ctx context.Context, vote *models.Vote) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Vote, error)
	GetAll(ctx context.Context) ([]models.Vote, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]models.Vote, error)
	GetByUserAndCategory(ctx context.Context, userID, categoryID uuid.UUID) ([]models.Vote, error)
	GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]models.Vote, error)
	Update(ctx context.Context, vote *models.Vote) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Analytics
	GetVoteCountsByCategory(ctx context.Context, categoryID uuid.UUID) ([]models.VoteCount, error)
	GetVoteCountsByNominee(ctx context.Context, nomineeID uuid.UUID) ([]models.VoteCount, error)
	GetUserVoteSummary(ctx context.Context, userID uuid.UUID) ([]models.UserVoteSummary, error)
	CountFreeVotesByUserAndCategory(ctx context.Context, userID, categoryID uuid.UUID) (int64, error)
}

type voteRepository struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) VoteRepository {
	return &voteRepository{db: db}
}

func (r *voteRepository) Create(ctx context.Context, vote *models.Vote) error {
	return r.db.WithContext(ctx).Create(vote).Error
}

func (r *voteRepository) GetByID(ctx context.Context, voteID uuid.UUID) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Nominee").
		Where("vote_id = ?", voteID).
		First(&vote).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &vote, err
}

func (r *voteRepository) GetAll(ctx context.Context) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Nominee").
		Find(&votes).Error
	return votes, err
}

func (r *voteRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Nominee").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&votes).Error
	return votes, err
}

func (r *voteRepository) GetByUserAndCategory(ctx context.Context, userID, categoryID uuid.UUID) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Nominee").
		Where("user_id = ? AND category_id = ?", userID, categoryID).
		Find(&votes).Error
	return votes, err
}

func (r *voteRepository) GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.WithContext(ctx).
		Preload("Nominee").
		Where("category_id = ?", categoryID).
		Find(&votes).Error
	return votes, err
}

func (r *voteRepository) Update(ctx context.Context, vote *models.Vote) error {
	return r.db.WithContext(ctx).Save(vote).Error
}

func (r *voteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Vote{}, "vote_id = ?", id).Error
}

// Analytics methods
func (r *voteRepository) GetVoteCountsByCategory(ctx context.Context, categoryID uuid.UUID) ([]models.VoteCount, error) {
	var results []models.VoteCount
	err := r.db.WithContext(ctx).
		Model(&models.Vote{}).
		Select(`
			nominee_id,
			nominees.name as nominee_name,
			category_id,
			COUNT(*) as total_votes,
			COUNT(CASE WHEN vote_type = 'free' THEN 1 END) as free_votes,
			COUNT(CASE WHEN vote_type = 'paid' THEN 1 END) as paid_votes
		`).
		Joins("JOIN nominees ON votes.nominee_id = nominees.nominee_id").
		Where("category_id = ?", categoryID).
		Group("nominee_id, nominees.name, category_id").
		Order("total_votes DESC").
		Scan(&results).Error
	return results, err
}

func (r *voteRepository) GetVoteCountsByNominee(ctx context.Context, nomineeID uuid.UUID) ([]models.VoteCount, error) {
	var results []models.VoteCount
	err := r.db.WithContext(ctx).
		Model(&models.Vote{}).
		Select(`
			nominee_id,
			nominees.name as nominee_name,
			category_id,
			COUNT(*) as total_votes,
			COUNT(CASE WHEN vote_type = 'free' THEN 1 END) as free_votes,
			COUNT(CASE WHEN vote_type = 'paid' THEN 1 END) as paid_votes
		`).
		Joins("JOIN nominees ON votes.nominee_id = nominees.nominee_id").
		Where("nominee_id = ?", nomineeID).
		Group("nominee_id, nominees.name, category_id").
		Scan(&results).Error
	return results, err
}

func (r *voteRepository) GetUserVoteSummary(ctx context.Context, userID uuid.UUID) ([]models.UserVoteSummary, error) {
	var results []models.UserVoteSummary
	err := r.db.WithContext(ctx).
		Model(&models.Vote{}).
		Select(`
			user_id,
			category_id,
			COUNT(CASE WHEN vote_type = 'free' THEN 1 END) as free_votes,
			COUNT(CASE WHEN vote_type = 'paid' THEN 1 END) as paid_votes
		`).
		Where("user_id = ?", userID).
		Group("user_id, category_id").
		Scan(&results).Error
	return results, err
}

func (r *voteRepository) CountFreeVotesByUserAndCategory(ctx context.Context, userID, categoryID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Vote{}).
		Where("user_id = ? AND category_id = ? AND vote_type = ?", userID, categoryID, models.VoteTypeFree).
		Count(&count).Error
	return count, err
}
