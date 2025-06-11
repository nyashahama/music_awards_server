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
	DeleteVote(ctx context.Context, voteID uuid.UUID) error
	GetAvailableVotes(ctx context.Context, userID uuid.UUID) (int, error)
}

type votingMechanismService struct {
	voteRepo repositories.VoteRepository
	userRepo repositories.UserRepository
}

func NewVotingMechanismService(voteRepo repositories.VoteRepository, userRepo repositories.UserRepository) VotingMechanismService {
	return &votingMechanismService{
		voteRepo: voteRepo,
		userRepo: userRepo,
	}
}

func (s *votingMechanismService) CastVote(ctx context.Context, userID, nomineeID, categoryID uuid.UUID) (*models.Vote, error) {
	// Check available votes
	if err := s.userRepo.DecrementAvailableVotes(ctx, userID); err != nil {
		return nil, fmt.Errorf("insufficient votes: %w", err)
	}

	vote := &models.Vote{
		VoteID:     uuid.New(),
		UserID:     userID,
		NomineeID:  nomineeID,
		CategoryID: categoryID,
	}

	if err := s.voteRepo.Create(ctx, vote); err != nil {
		// Rollback vote count if vote creation fails
		s.userRepo.IncrementAvailableVotes(ctx, userID)
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

func (s *votingMechanismService) DeleteVote(ctx context.Context, voteID uuid.UUID) error {
	vote, err := s.voteRepo.GetByID(ctx, voteID)
	if err != nil {
		return err
	}

	if err := s.voteRepo.Delete(ctx, voteID); err != nil {
		return err
	}

	// Return vote to user
	if err := s.userRepo.IncrementAvailableVotes(ctx, vote.UserID); err != nil {
		return fmt.Errorf("failed to return vote: %w", err)
	}
	return nil
}

func (s *votingMechanismService) GetAvailableVotes(ctx context.Context, userID uuid.UUID) (int, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return user.AvailableVotes, nil
}

