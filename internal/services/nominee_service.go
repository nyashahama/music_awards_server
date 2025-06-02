package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
	"gorm.io/datatypes"
)

// NomineeService handles nominee-related operations
type NomineeService interface {
	CreateNominee(ctx context.Context, nominee models.Nominee, categoryIDs []uuid.UUID) (*models.Nominee, error)
	UpdateNominee(ctx context.Context, nomineeID uuid.UUID, updateData map[string]interface{}) (*models.Nominee, error)
	DeleteNominee(ctx context.Context, nomineeID uuid.UUID) error
	AddNomineeCategory(ctx context.Context, nomineeID, categoryID uuid.UUID) error
	RemoveNomineeCategory(ctx context.Context, nomineeID, categoryID uuid.UUID) error
	GetNomineeDetails(ctx context.Context, nomineeID uuid.UUID) (*models.Nominee, error)
	SetNominationPeriod(ctx context.Context, categoryID uuid.UUID, start, end time.Time) error
	GetAllNominees(ctx context.Context) ([]models.Nominee, error)
}

type nomineeService struct {
	// Dependencies
	repo repositories.NomineeRepository
}

func NewNomineeService(repo repositories.NomineeRepository) NomineeService {
	return &nomineeService{repo: repo}
}

func (s *nomineeService) CreateNominee(ctx context.Context, nominee models.Nominee, categoryIDs []uuid.UUID) (*models.Nominee, error) {
	// create nominee
	if err := s.repo.Create(&nominee); err != nil {
		return nil, err
	}

	for _, catId := range categoryIDs {
		if err := s.AddNomineeCategory(ctx, nominee.NomineeID, catId); err != nil {
			return nil, err
		}
	}
	return s.repo.GetByID(nominee.NomineeID)
}
func (s *nomineeService) UpdateNominee(ctx context.Context, nomineeID uuid.UUID, updateData map[string]interface{}) (*models.Nominee, error) {
	nominee, err := s.repo.GetByID(nomineeID)
	if err != nil {
		return nil, err
	}

	if name, ok := updateData["name"].(string); ok {
		nominee.Name = name
	}

	if desc, ok := updateData["description"].(string); ok {
		nominee.Description = desc
	}

	if works, ok := updateData["sample_works"].(datatypes.JSON); ok {
		nominee.SampleWorks = works
	}

	if img, ok := updateData["image_url"].(string); ok {
		nominee.ImageURL = img
	}

	if err := s.repo.Update(nominee); err != nil {
		return nil, err
	}

	return nominee, nil

}
func (s *nomineeService) DeleteNominee(ctx context.Context, nomineeID uuid.UUID) error {
	return s.repo.Delete(nomineeID)
}
func (s *nomineeService) AddNomineeCategory(ctx context.Context, nomineeID, categoryID uuid.UUID) error {
	return s.repo.AddCategory(nomineeID, categoryID)
}
func (s *nomineeService) RemoveNomineeCategory(ctx context.Context, nomineeID, categoryID uuid.UUID) error {
	return s.repo.RemoveCategory(nomineeID, categoryID)
}
func (s *nomineeService) GetNomineeDetails(ctx context.Context, nomineeID uuid.UUID) (*models.Nominee, error) {
	return s.repo.GetByID(nomineeID)
}
func (s *nomineeService) SetNominationPeriod(ctx context.Context, categoryID uuid.UUID, start, end time.Time) error {
	return errors.New("cool for now")
}

func (s *nomineeService) GetAllNominees(ctx context.Context) ([]models.Nominee, error) {
	return s.repo.GetAll()
}
