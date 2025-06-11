package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/gorm"
)

// Nominee Repository
type NomineeRepository interface {
	Create(ctx context.Context, nominee *models.Nominee) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Nominee, error)
	GetAll(ctx context.Context) ([]models.Nominee, error)
	Update(ctx context.Context, nominee *models.Nominee) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type nomineeRepository struct {
	db *gorm.DB
}

func NewNomineeRepository(db *gorm.DB) NomineeRepository {
	return &nomineeRepository{db: db}
}

func (r *nomineeRepository) Create(ctx context.Context, nominee *models.Nominee) error {
	return r.db.WithContext(ctx).Create(nominee).Error
}

func (r *nomineeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Nominee, error) {
	var nominee models.Nominee
	err := r.db.WithContext(ctx).Preload("Categories", func(db *gorm.DB) *gorm.DB {
		return db.WithContext(ctx).Select("category_id", "name")
	}).First(&nominee, "nominee_id = ?", id).Error
	return &nominee, err
}

func (r *nomineeRepository) GetAll(ctx context.Context) ([]models.Nominee, error) {
	var nominees []models.Nominee
	err := r.db.WithContext(ctx).Preload("Categories", func(db *gorm.DB) *gorm.DB {
		return db.Select("category_id", "name")
	}).Find(&nominees).Error
	return nominees, err
}

func (r *nomineeRepository) Update(ctx context.Context, nominee *models.Nominee) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(nominee).Error
}

func (r *nomineeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Nominee{}, "nominee_id = ?", id).Error
}
