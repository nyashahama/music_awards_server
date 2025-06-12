package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/gorm"
)

type NomineeCategoryRepository interface {
	AddCategory(ctx context.Context, nomineeId, categoryId uuid.UUID) error
	RemoveCategory(ctx context.Context, nomineeId, categoryId uuid.UUID) error
	GetCategoriesForNominee(ctx context.Context, nomineeId uuid.UUID) ([]models.Category, error)
	GetNomineesForCategory(ctx context.Context, categoryID uuid.UUID) ([]models.Nominee, error)
	SetCategories(ctx context.Context, nomineeId uuid.UUID, categoryIds []uuid.UUID) error
}

type nomineeCategoryRepository struct {
	db *gorm.DB
}

func NewNomineeCategoryRepository(db *gorm.DB) NomineeCategoryRepository {
	return &nomineeCategoryRepository{db: db}
}

func (r *nomineeCategoryRepository) AddCategory(ctx context.Context, nomineeId, categoryId uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Nominee{NomineeID: nomineeId}).
		Association("Categories").
		Append(&models.Category{CategoryID: categoryId})
}

func (r *nomineeCategoryRepository) RemoveCategory(ctx context.Context, nomineeId, categoryId uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Nominee{NomineeID: nomineeId}).
		Association("Categories").
		Delete(&models.Category{CategoryID: categoryId})
}

func (r *nomineeCategoryRepository) GetCategoriesForNominee(ctx context.Context, nomineeId uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).
		Select("category_id", "name", "description").
		Joins("JOIN nominee_categories ON nominee_categories.category_id = categories.category_id").
		Where("nominee_categories.nominee_id = ?", nomineeId).
		Find(&categories).Error
	return categories, err
}

func (r *nomineeCategoryRepository) GetNomineesForCategory(ctx context.Context, categoryID uuid.UUID) ([]models.Nominee, error) {
	var nominees []models.Nominee
	err := r.db.WithContext(ctx).
		Select("nominee_id", "name").
		Joins("JOIN nominee_categories ON nominee_categories.nominee_id = nominees.nominee_id").
		Where("nominee_categories.category_id = ?", categoryID).
		Find(&nominees).Error
	return nominees, err
}

func (r *nomineeCategoryRepository) SetCategories(ctx context.Context, nomineeId uuid.UUID, categoryIds []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Clear existing associations
		if err := tx.Model(&models.Nominee{NomineeID: nomineeId}).
			Association("Categories").Clear(); err != nil {
			return err
		}

		// Only add new categories if there are any
		if len(categoryIds) > 0 {
			categories := make([]models.Category, len(categoryIds))
			for i, id := range categoryIds {
				categories[i] = models.Category{CategoryID: id}
			}

			// Use Replace instead of multiple Appends
			if err := tx.Model(&models.Nominee{NomineeID: nomineeId}).
				Association("Categories").Replace(categories); err != nil {
				return err
			}
		}
		return nil
	})
}
