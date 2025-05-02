package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
)

// ResultsService handles voting results and reporting
type ResultsService interface {
	GetCategoryResults(ctx context.Context, categoryID uuid.UUID) ([]models.Nominee, error)
	GenerateVotingReport(ctx context.Context, filters map[string]interface{}) (*models.Vote, error)
	ExportResults(ctx context.Context, format string) ([]byte, error)
	GetHistoricalResults(ctx context.Context, year int) (map[uuid.UUID]models.Nominee, error)
	GetRealTimeTallies(ctx context.Context) (map[uuid.UUID]models.Category, error)
}

type resultsService struct {
	// Dependencies
}