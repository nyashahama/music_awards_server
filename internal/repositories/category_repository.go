package repositories

import (
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/gorm"
)

// Category Repository
type CategoryRepository interface {
	Create(category *models.Category) error
	GetByID(id uuid.UUID) (*models.Category, error)
	GetByName(name string) (*models.Category, error)
	GetAll() ([]models.Category, error)
	Update(category *models.Category) error
	Delete(id uuid.UUID) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, "category_id = ?", id).Error
	return &category, err
}

func (r *categoryRepository) GetByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, "name = ?", name).Error
	return &category, err
}

func (r *categoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Category{}, "category_id = ?", id).Error
}