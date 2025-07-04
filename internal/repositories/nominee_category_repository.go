package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/gorm"
)

//NomineeCategory Repository
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
		Select("categories.category_id", "categories.name").
		Joins("JOIN nominee_categories ON categories.category_id = nominee_categories.category_id").
		Where("nominee_categories.nominee_id = ?", nomineeId).
		Find(&categories).Error
	return categories, err
}

func (r *nomineeCategoryRepository) GetNomineesForCategory(ctx context.Context, categoryID uuid.UUID) ([]models.Nominee, error) {
	var nominees []models.Nominee
	err := r.db.WithContext(ctx).
		Select("nominees.nominee_id", "nominees.name", "nominees.image_url").
		Joins("JOIN nominee_categories ON nominees.nominee_id = nominee_categories.nominee_id").
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

			if err := tx.Model(&models.Nominee{NomineeID: nomineeId}).
				Association("Categories").Replace(categories); err != nil {
				return err
			}
		}
		return nil
	})
}
