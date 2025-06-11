package repositories

import (
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/gorm"
)

type NomineeCategoryRepository interface {
	AddCategory(nomineeId, categoryId uuid.UUID) error
	RemoveCategory(nomineeId, categoryId uuid.UUID) error
	GetCategoriesForNominee(nomineeId uuid.UUID) ([]models.Category, error)
	GetNomineesForCategory(categoryID uuid.UUID) ([]models.Nominee, error)
	SetCategories(nomineeId uuid.UUID, categoryIds []uuid.UUID) error
}

type nomineeCategoryRepository struct {
	db *gorm.DB
}

func NewNomineeCategoryRepository(db *gorm.DB) NomineeCategoryRepository {
	return &nomineeCategoryRepository{db: db}
}

func (r *nomineeCategoryRepository) AddCategory(nomineeId, categoryId uuid.UUID) error {
	return r.db.Model(&models.Nominee{NomineeID: nomineeId}).
		Association("Categories").
		Append(&models.Category{CategoryID: categoryId})
}

func (r *nomineeCategoryRepository) RemoveCategory(nomineeId, categoryId uuid.UUID) error {
	return r.db.Model(&models.Nominee{NomineeID: nomineeId}).
		Association("Categories").
		Delete(&models.Category{CategoryID: categoryId})
}

func (r *nomineeCategoryRepository) GetCategoriesForNominee(nomineeId uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Joins("JOIN nominee_categories ON nominee_categories.category_id = categories.category_id").
		Where("nominee_categories.nominee_id = ?", nomineeId).
		Find(&categories).Error
	return categories, err
}

func (r *nomineeCategoryRepository) GetNomineesForCategory(categoryID uuid.UUID) ([]models.Nominee, error) {
	var nominees []models.Nominee
	err := r.db.Joins("JOIN nominee_categories ON nominee_categories.nominee_id = nominees.nominee_id").
		Where("nominee_categories.category_id = ?", categoryID).
		Find(&nominees).Error
	return nominees, err
}
func (r *nomineeCategoryRepository) SetCategories(nomineeId uuid.UUID, categoryIds []uuid.UUID) error {
	categories := make([]models.Category, len(categoryIds))
	for i, id := range categoryIds {
		categories[i] = models.Category{CategoryID: id}
	}

	return r.db.Model(&models.Nominee{NomineeID: nomineeId}).
		Association("Categories").
		Replace(categories)
}
