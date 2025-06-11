package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
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
	category := &models.Category{
		CategoryID:  uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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
		return nil, fmt.Errorf("category not found")
	}

	category.Name = name
	category.Description = description
	category.UpdatedAt = time.Now()

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
