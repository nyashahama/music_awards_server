package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
)

// ViewService handles nominee presentation and filtering
type ViewService interface {
	ListAllNominees(ctx context.Context, filters map[string]interface{}) ([]models.Nominee, error)
	GetNomineesByCategory(ctx context.Context, categoryID uuid.UUID) ([]models.Nominee, error)
	SearchNominees(ctx context.Context, query string) ([]models.Nominee, error)
	GetPopularNominees(ctx context.Context, limit int) ([]models.Nominee, error)
	GetNomineeDetails(ctx context.Context, nomineeID uuid.UUID) (*models.Nominee, error)
	TrackNomineeView(ctx context.Context, nomineeID uuid.UUID) error
}

type viewService struct {
	// Dependencies
}
