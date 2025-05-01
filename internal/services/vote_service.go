package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
)

// VotingMechanismService handles core voting logic
type VotingMechanismService interface {
	CastVote(ctx context.Context, userID, nomineeID, categoryID uuid.UUID) (*models.Vote, error)
	GetVote(ctx context.Context, voteID uuid.UUID) (*models.Vote, error)
	ChangeVote(ctx context.Context, voteID uuid.UUID, newNomineeID uuid.UUID) (*models.Vote, error)
	GetUserVotes(ctx context.Context, userID uuid.UUID) ([]models.Vote, error)
	HasVotedInCategory(ctx context.Context, userID, categoryID uuid.UUID) (bool, error)
	GetCategoryVotes(ctx context.Context, categoryID uuid.UUID) ([]models.Vote, error)
	ValidateVotingPeriod(ctx context.Context, categoryID uuid.UUID) (bool, error)
}

type votingMechanismService struct {
	// Dependencies
}