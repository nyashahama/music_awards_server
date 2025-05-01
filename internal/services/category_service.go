package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
)

// CategoryService handles category operations
type CategoryService interface {
	CreateCategory(ctx context.Context, name, description string) (*models.Category, error)
	UpdateCategory(ctx context.Context, categoryID uuid.UUID, name, description string) (*models.Category, error)
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
	GetCategoryDetails(ctx context.Context, categoryID uuid.UUID) (*models.Category, error)
	ListAllCategories(ctx context.Context) ([]models.Category, error)
}

type categoryService struct {
	// Dependencies
}