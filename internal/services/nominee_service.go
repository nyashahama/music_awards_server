package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/dtos"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
)

var (
	ErrNomineeNotFound = errors.New("nominee not found")
	ErrInvalidJSON     = errors.New("invalid JSON data")
)

type NomineeService interface {
	CreateNominee(ctx context.Context, req dtos.CreateNomineeRequest) (*models.Nominee, error)
	UpdateNominee(ctx context.Context, nomineeID uuid.UUID, req dtos.UpdateNomineeRequest) (*models.Nominee, error)
	DeleteNominee(ctx context.Context, nomineeID uuid.UUID) error
	GetNomineeDetails(ctx context.Context, nomineeID uuid.UUID) (*models.Nominee, error)
	GetAllNominees(ctx context.Context) ([]models.Nominee, error)
}

type nomineeService struct {
	repo                repositories.NomineeRepository
	categoryRepo        repositories.CategoryRepository
	nomineeCategoryRepo repositories.NomineeCategoryRepository
}

func NewNomineeService(
	repo repositories.NomineeRepository,
	categoryRepo repositories.CategoryRepository,
	nomineeCategoryRepo repositories.NomineeCategoryRepository,
) NomineeService {
	return &nomineeService{
		repo:                repo,
		categoryRepo:        categoryRepo,
		nomineeCategoryRepo: nomineeCategoryRepo,
	}
}

func (s *nomineeService) CreateNominee(ctx context.Context, req dtos.CreateNomineeRequest) (*models.Nominee, error) {
	// Validate sample works JSON
	if req.SampleWorks != nil && !json.Valid(req.SampleWorks) {
		return nil, ErrInvalidJSON
	}

	// Validate categories exist
	for _, categoryID := range req.CategoryIDs {
		category, err := s.categoryRepo.GetByID(ctx, categoryID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate category: %w", err)
		}
		if category == nil {
			return nil, ErrCategoryNotFound
		}
	}

	nominee := &models.Nominee{
		NomineeID:   uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		SampleWorks: req.SampleWorks,
		ImageURL:    req.ImageURL,
	}

	// Create nominee
	if err := s.repo.Create(ctx, nominee); err != nil {
		return nil, fmt.Errorf("failed to create nominee: %w", err)
	}

	// Set categories if provided
	if len(req.CategoryIDs) > 0 {
		if err := s.nomineeCategoryRepo.SetCategories(ctx, nominee.NomineeID, req.CategoryIDs); err != nil {
			return nil, fmt.Errorf("failed to set categories: %w", err)
		}
	}

	// Reload nominee to get associations
	return s.repo.GetByID(ctx, nominee.NomineeID)
}

func (s *nomineeService) UpdateNominee(ctx context.Context, nomineeID uuid.UUID, req dtos.UpdateNomineeRequest) (*models.Nominee, error) {
	nominee, err := s.repo.GetByID(ctx, nomineeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nominee: %w", err)
	}
	if nominee == nil {
		return nil, ErrNomineeNotFound
	}

	// Update fields
	if req.Name != nil {
		nominee.Name = *req.Name
	}
	if req.Description != nil {
		nominee.Description = *req.Description
	}
	if req.ImageURL != nil {
		nominee.ImageURL = *req.ImageURL
	}
	if req.SampleWorks != nil {
		if !json.Valid(*req.SampleWorks) {
			return nil, ErrInvalidJSON
		}
		nominee.SampleWorks = *req.SampleWorks
	}

	// Update categories if provided
	// Validate new categories if provided
	if req.CategoryIDs != nil {
		for _, categoryID := range *req.CategoryIDs {
			category, err := s.categoryRepo.GetByID(ctx, categoryID)
			if err != nil {
				return nil, fmt.Errorf("failed to validate category: %w", err)
			}
			if category == nil {
				return nil, ErrCategoryNotFound
			}
		}
	}
	// Save nominee
	if err := s.repo.Update(ctx, nominee); err != nil {
		return nil, fmt.Errorf("failed to update nominee: %w", err)
	}

	// Reload to get updated associations
	return s.repo.GetByID(ctx, nomineeID)
}

func (s *nomineeService) DeleteNominee(ctx context.Context, nomineeID uuid.UUID) error {
	if err := s.repo.Delete(ctx, nomineeID); err != nil {
		return fmt.Errorf("failed to delete nominee: %w", err)
	}
	return nil
}

func (s *nomineeService) GetNomineeDetails(ctx context.Context, nomineeID uuid.UUID) (*models.Nominee, error) {
	nominee, err := s.repo.GetByID(ctx, nomineeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nominee: %w", err)
	}
	if nominee == nil {
		return nil, ErrNomineeNotFound
	}
	return nominee, nil
}

func (s *nomineeService) GetAllNominees(ctx context.Context) ([]models.Nominee, error) {
	return s.repo.GetAll(ctx)
}
