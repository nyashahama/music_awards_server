package repositories

import (
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/gorm"
)

// Nominee Repository
type NomineeRepository interface {
	Create(nominee *models.Nominee) error
	GetByID(id uuid.UUID) (*models.Nominee, error)
	GetAll() ([]models.Nominee, error)
	Update(nominee *models.Nominee) error
	Delete(id uuid.UUID) error
	GetByCategory(categoryID uuid.UUID) ([]models.Nominee, error)
}

type nomineeRepository struct {
	db *gorm.DB
}

func NewNomineeRepository(db *gorm.DB) NomineeRepository {
	return &nomineeRepository{db: db}
}

func (r *nomineeRepository) Create(nominee *models.Nominee) error {
	return r.db.Create(nominee).Error
}

func (r *nomineeRepository) GetByID(id uuid.UUID) (*models.Nominee, error) {
	var nominee models.Nominee
	err := r.db.Preload("Categories").First(&nominee, "nominee_id = ?", id).Error
	return &nominee, err
}

func (r *nomineeRepository) GetAll() ([]models.Nominee, error) {
	var nominees []models.Nominee
	err := r.db.Preload("Categories").Find(&nominees).Error
	return nominees, err
}

func (r *nomineeRepository) Update(nominee *models.Nominee) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(nominee).Error
}

func (r *nomineeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Nominee{}, "nominee_id = ?", id).Error
}

func (r *nomineeRepository) GetByCategory(categoryID uuid.UUID) ([]models.Nominee, error) {
	var nominees []models.Nominee
	err := r.db.Joins("JOIN nominee_categories ON nominees.nominee_id = nominee_categories.nominee_id").
		Where("nominee_categories.category_id = ?", categoryID).
		Preload("Categories").
		Find(&nominees).Error
	return nominees, err
}
