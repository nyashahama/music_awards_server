package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/gorm"
)

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
	err := r.db.WithContext(ctx).
		Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN nominee_categories ON categories.category_id = nominee_categories.category_id").
				Where("nominee_categories.nominee_id = ?", id)
		}).
		First(&nominee, "nominee_id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &nominee, err
}

func (r *nomineeRepository) GetAll(ctx context.Context) ([]models.Nominee, error) {
	var nominees []models.Nominee
	err := r.db.WithContext(ctx).
		Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("category_id", "name")
		}).
		Joins("LEFT JOIN nominee_categories ON nominees.nominee_id = nominee_categories.nominee_id").
		Group("nominees.nominee_id").
		Find(&nominees).Error
	return nominees, err
}

func (r *nomineeRepository) Update(ctx context.Context, nominee *models.Nominee) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(nominee).Error
}

func (r *nomineeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Nominee{}, "nominee_id = ?", id).Error
}
