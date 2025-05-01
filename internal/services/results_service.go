package services

import (
	"context"

	"github.com/google/uuid"
)

// ResultsService handles voting results and reporting
type ResultsService interface {
	GetCategoryResults(ctx context.Context, categoryID uuid.UUID) ([]models.NomineeResult, error)
	GenerateVotingReport(ctx context.Context, filters map[string]interface{}) (*models.VotingReport, error)
	ExportResults(ctx context.Context, format string) ([]byte, error)
	GetHistoricalResults(ctx context.Context, year int) (map[uuid.UUID]models.NomineeResult, error)
	GetRealTimeTallies(ctx context.Context) (map[uuid.UUID]models.CategoryTally, error)
}

type resultsService struct {
	// Dependencies
}