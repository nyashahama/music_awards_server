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
	return r.db.Exec(`
		INSERT INTO nominee_categories (nominee_id, category_id)
		VALUES (?, ?)
		ON CONFLICT (nominee_id, category_id) DO NOTHING
	`, nomineeId, categoryId).Error
}

func (r *nomineeCategoryRepository) RemoveCategory(nomineeId, categoryId uuid.UUID) error {
	return r.db.Exec(`
		DELETE FROM nominee_categories
		WHERE nominee_id = ? AND category_id = ?
		`, nomineeId, categoryId).Error
}

func (r *nomineeCategoryRepository) GetCategoriesForNominee(nomineeId uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Raw(`
		SELECT n.*
		FROM nominees n 
		JOIN nominee_categories nc ON n.nominee_id = nc.nominee_id
		WHERE nc.category_id = ?
		`, nomineeId).Scan(&categories).Error
	return categories, err
}

func (r *nomineeCategoryRepository) GetNomineesForCategory(categoryID uuid.UUID) ([]models.Nominee, error) {
	var nominees []models.Nominee
	err := r.db.Raw(`
		SELECT n.* 
		FROM nominees n
		JOIN nominee_categories nc ON n.nominee_id = nc.nominee_id
		WHERE nc.category_id = ?
	`, categoryID).Scan(&nominees).Error
	return nominees, err
}

func (r *nomineeCategoryRepository) SetCategories(nomineeId uuid.UUID, categoryIds []uuid.UUID) error {
	// Start transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Clear existing categories
	if err := tx.Exec(`
		DELETE FROM nominee_categories 
		WHERE nominee_id = ?
	`, nomineeId).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add new categories
	for _, categoryId := range categoryIds {
		if err := tx.Exec(`
			INSERT INTO nominee_categories (nominee_id, category_id)
			VALUES (?, ?)
			ON CONFLICT (nominee_id, category_id) DO NOTHING
		`, nomineeId, categoryId).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
