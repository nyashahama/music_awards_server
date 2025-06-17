package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category name already exists")
)

// CategoryService handles category operations
type CategoryService interface {
	CreateCategory(ctx context.Context, name, description string) (*models.Category, error)
	UpdateCategory(ctx context.Context, categoryID uuid.UUID, name, description string) (*models.Category, error)
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
	GetCategoryDetails(ctx context.Context, categoryID uuid.UUID) (*models.Category, error)
	ListAllCategories(ctx context.Context) ([]models.Category, error)
	ListActiveCategories(ctx context.Context) ([]models.Category, error)
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(ctx context.Context, name, description string) (*models.Category, error) {
	// Check for existing category

	existing, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check category name: %w", err)
	}
	if existing != nil {
		return nil, ErrCategoryExists
	}

	category := &models.Category{
		CategoryID:  uuid.New(),
		Name:        name,
		Description: description,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return category, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, categoryID uuid.UUID, name, description string) (*models.Category, error) {
	category, err := s.repo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	// Check for name conflict if name changed
	if name != "" && name != category.Name {
		existing, err := s.repo.GetByName(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("failed to check category name: %w", err)
		}
		if existing != nil {
			return nil, ErrCategoryExists
		}
		category.Name = name
	}

	if description != "" {
		category.Description = description
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}
	return category, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	return s.repo.Delete(ctx, categoryID)
}

func (s *categoryService) GetCategoryDetails(ctx context.Context, categoryID uuid.UUID) (*models.Category, error) {
	category, err := s.repo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}

func (s *categoryService) ListAllCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *categoryService) ListActiveCategories(ctx context.Context) ([]models.Category, error) {
	allCategories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	var activeCategories []models.Category
	for _, cat := range allCategories {
		if len(cat.Votes) > 0 {
			activeCategories = append(activeCategories, cat)
		}
	}
	return activeCategories, nil
}
