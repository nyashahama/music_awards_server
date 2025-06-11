package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	//"time"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
	// "gorm.io/datatypes"
)

// NomineeService handles nominee-related operations
type NomineeService interface {
	CreateNominee(ctx context.Context, nominee models.Nominee) (*models.Nominee, error)
	UpdateNominee(ctx context.Context, nomineeID uuid.UUID, updateData map[string]interface{}) (*models.Nominee, error)
	DeleteNominee(ctx context.Context, nomineeID uuid.UUID) error
	GetNomineeDetails(ctx context.Context, nomineeID uuid.UUID) (*models.Nominee, error)
	GetAllNominees(ctx context.Context) ([]models.Nominee, error)
}

type nomineeService struct {
	// Dependencies
	repo repositories.NomineeRepository
}

func NewNomineeService(repo repositories.NomineeRepository) NomineeService {
	return &nomineeService{repo: repo}

}

func (s *nomineeService) CreateNominee(ctx context.Context, nominee models.Nominee) (*models.Nominee, error) {
	if err := s.repo.Create(ctx, &nominee); err != nil {
		return nil, fmt.Errorf("failed to create nominee: %w", err)
	}
	return s.repo.GetByID(ctx, nominee.NomineeID)
}

func (s *nomineeService) UpdateNominee(ctx context.Context, nomineeID uuid.UUID, updateData map[string]interface{}) (*models.Nominee, error) {
	nominee, err := s.repo.GetByID(ctx, nomineeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nominee: %w", err)
	}
	if nominee == nil {
		return nil, errors.New("nominee not found")
	}

	if name, ok := updateData["name"].(string); ok {
		nominee.Name = name
	}
	if desc, ok := updateData["description"].(string); ok {
		nominee.Description = desc
	}
	if img, ok := updateData["image_url"].(string); ok {
		nominee.ImageURL = img
	}

	if works, ok := updateData["sample_works"]; ok {
		raw, err := json.Marshal(works)
		if err != nil {
			return nil, fmt.Errorf("could not marshal sample_works: %w", err)
		}
		nominee.SampleWorks = json.RawMessage(raw)
	}

	if err := s.repo.Update(ctx, nominee); err != nil {
		return nil, fmt.Errorf("failed to update nominee: %w", err)
	}

	return nominee, nil
}

func (s *nomineeService) DeleteNominee(ctx context.Context, nomineeID uuid.UUID) error {
	return s.repo.Delete(ctx, nomineeID)
}

func (s *nomineeService) GetNomineeDetails(ctx context.Context, nomineeID uuid.UUID) (*models.Nominee, error) {
	nominee, err := s.repo.GetByID(ctx, nomineeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nominee: %w", err)
	}
	if nominee == nil {
		return nil, errors.New("nominee not found")
	}
	return nominee, nil
}

func (s *nomineeService) GetAllNominees(ctx context.Context) ([]models.Nominee, error) {
	return s.repo.GetAll(ctx)
}
