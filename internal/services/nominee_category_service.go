// services/nominee_category_service.go
package services

import (
	"context"
	//"errors"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
)

var (
//ErrInvalidID = errors.New("invalid ID")
)

type NomineeCategoryService interface {
	AddCategory(ctx context.Context, nomineeID, categoryID uuid.UUID) error
	RemoveCategory(ctx context.Context, nomineeID, categoryID uuid.UUID) error
	SetCategories(ctx context.Context, nomineeID uuid.UUID, categoryIDs []uuid.UUID) error
	GetCategories(ctx context.Context, nomineeID uuid.UUID) ([]models.Category, error)
	GetNominees(ctx context.Context, categoryID uuid.UUID) ([]models.Nominee, error)
}

type nomineeCategoryService struct {
	repo repositories.NomineeCategoryRepository
}

func NewNomineeCategoryService(repo repositories.NomineeCategoryRepository) NomineeCategoryService {
	return &nomineeCategoryService{repo: repo}
}

func (s *nomineeCategoryService) AddCategory(ctx context.Context, nomineeID, categoryID uuid.UUID) error {
	if nomineeID == uuid.Nil || categoryID == uuid.Nil {
		return ErrInvalidID
	}
	return s.repo.AddCategory(ctx, nomineeID, categoryID)
}

func (s *nomineeCategoryService) RemoveCategory(ctx context.Context, nomineeID, categoryID uuid.UUID) error {
	if nomineeID == uuid.Nil || categoryID == uuid.Nil {
		return ErrInvalidID
	}
	return s.repo.RemoveCategory(ctx, nomineeID, categoryID)
}

func (s *nomineeCategoryService) SetCategories(ctx context.Context, nomineeID uuid.UUID, categoryIDs []uuid.UUID) error {
	if nomineeID == uuid.Nil {
		return ErrInvalidID
	}
	return s.repo.SetCategories(ctx, nomineeID, categoryIDs)
}

func (s *nomineeCategoryService) GetCategories(ctx context.Context, nomineeID uuid.UUID) ([]models.Category, error) {
	if nomineeID == uuid.Nil {
		return nil, ErrInvalidID
	}
	return s.repo.GetCategoriesForNominee(ctx, nomineeID)
}

func (s *nomineeCategoryService) GetNominees(ctx context.Context, categoryID uuid.UUID) ([]models.Nominee, error) {
	if categoryID == uuid.Nil {
		return nil, ErrInvalidID
	}
	return s.repo.GetNomineesForCategory(ctx, categoryID)
}
