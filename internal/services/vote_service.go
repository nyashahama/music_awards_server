package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
)

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
	voteRepo repositories.VoteRepository
}

func NewVotingMechanismService(voteRepo repositories.VoteRepository) VotingMechanismService {
	return &votingMechanismService{voteRepo: voteRepo}
}

func (s *votingMechanismService) CastVote(ctx context.Context, userID, nomineeID, categoryID uuid.UUID) (*models.Vote, error) {
	vote := &models.Vote{
		VoteID:     uuid.New(),
		UserID:     userID,
		NomineeID:  nomineeID,
		CategoryID: categoryID,
	}

	if err := s.voteRepo.Create(ctx, vote); err != nil {
		return nil, fmt.Errorf("failed to cast vote: %w", err)
	}
	return vote, nil
}

func (s *votingMechanismService) GetVote(ctx context.Context, voteID uuid.UUID) (*models.Vote, error) {
	vote, err := s.voteRepo.GetByID(ctx, voteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vote: %w", err)
	}
	return vote, nil
}

func (s *votingMechanismService) ChangeVote(ctx context.Context, voteID uuid.UUID, newNomineeID uuid.UUID) (*models.Vote, error) {
	vote, err := s.voteRepo.GetByID(ctx, voteID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vote: %w", err)
	}
	// Update nominee
	vote.NomineeID = newNomineeID
	if err := s.voteRepo.Update(ctx, vote); err != nil {
		return nil, fmt.Errorf("failed to update vote: %w", err)
	}
	return vote, nil
}

func (s *votingMechanismService) GetUserVotes(ctx context.Context, userID uuid.UUID) ([]models.Vote, error) {
	votes, err := s.voteRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user votes: %w", err)
	}
	return votes, nil
}

func (s *votingMechanismService) HasVotedInCategory(ctx context.Context, userID, categoryID uuid.UUID) (bool, error) {
	vote, err := s.voteRepo.GetByUserAndCategory(ctx, userID, categoryID)
	if err != nil {
		return false, fmt.Errorf("error checking existing vote: %w", err)
	}
	return vote != nil, nil
}

func (s *votingMechanismService) GetCategoryVotes(ctx context.Context, categoryID uuid.UUID) ([]models.Vote, error) {
	allVotes, err := s.voteRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve votes: %w", err)
	}
	var filtered []models.Vote
	for _, v := range allVotes {
		if v.CategoryID == categoryID {
			filtered = append(filtered, v)
		}
	}
	return filtered, nil
}

func (s *votingMechanismService) ValidateVotingPeriod(ctx context.Context, categoryID uuid.UUID) (bool, error) {
	// TODO: implement this
	return true, nil
}

