package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
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
}

type nomineeService struct {
	// Dependencies
}
