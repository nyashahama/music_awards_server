package repositories

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/gorm"
)

// Vote Repository
type VoteRepository interface {
	Create(vote *models.Vote) error
	GetByID(id uuid.UUID) (*models.Vote, error)
	GetAll() ([]models.Vote, error)
	GetByUser(userID uuid.UUID) ([]models.Vote, error)
	GetByUserAndCategory(userID, categoryID uuid.UUID) (*models.Vote, error)
	Update(vote *models.Vote) error
	Delete(id uuid.UUID) error
}

type voteRepository struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) VoteRepository {
	return &voteRepository{db: db}
}

func (r *voteRepository) Create(vote *models.Vote) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if user already voted in this category
		var existing models.Vote
		if err := tx.Where("user_id = ? AND category_id = ?", vote.UserID, vote.CategoryID).
			First(&existing).Error; err == nil {
			return fmt.Errorf("user already voted in this category")
		}

		return tx.Create(vote).Error
	})
}
func (r *voteRepository) GetByID(id uuid.UUID) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.First(&vote, "vote_id = ?", id).Error
	return &vote, err
}

func (r *voteRepository) GetAll() ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Find(&votes).Error
	return votes, err
}

func (r *voteRepository) GetByUser(userID uuid.UUID) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("category_id", "name")
	}).
		Preload("Nominee", func(db *gorm.DB) *gorm.DB {
			return db.Select("nominee_id", "name")
		}).
		Where("user_id = ?", userID).
		Find(&votes).Error
	return votes, err
}
func (r *voteRepository) GetByUserAndCategory(userID, categoryID uuid.UUID) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.Where("user_id = ? AND category_id = ?", userID, categoryID).First(&vote).Error
	return &vote, err
}

func (r *voteRepository) Update(vote *models.Vote) error {
	return r.db.Save(vote).Error
}

func (r *voteRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Vote{}, "vote_id = ?", id).Error
}
