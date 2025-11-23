// Package services - Improved voting service with business logic
package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/repositories"
)

var (
	ErrNoVotesAvailable         = errors.New("no votes available")
	ErrNoFreeVotesAvailable     = errors.New("no free votes available")
	ErrNoPaidVotesAvailable     = errors.New("no paid votes available")
	ErrAlreadyVotedWithFreeVote = errors.New("already voted in this category with free vote")
	ErrVotingPeriodClosed       = errors.New("voting period is closed")
	ErrCategoryNotActive        = errors.New("category is not active")
	ErrNomineeNotActive         = errors.New("nominee is not active")
	ErrNomineeNotInCategory     = errors.New("nominee not in category")
	ErrInvalidVoteType          = errors.New("invalid vote type")
)

type VotingService interface {
	// Core voting
	CastVote(ctx context.Context, userID, nomineeID, categoryID uuid.UUID, usePaidVote bool) (*models.Vote, error)
	ChangeVote(ctx context.Context, voteID uuid.UUID, newNomineeID uuid.UUID) (*models.Vote, error)
	DeleteVote(ctx context.Context, voteID uuid.UUID) error

	// Queries
	GetVote(ctx context.Context, voteID uuid.UUID) (*models.Vote, error)
	GetUserVotes(ctx context.Context, userID uuid.UUID) ([]models.Vote, error)
	GetAllVotes(ctx context.Context) ([]models.Vote, error)
	GetAvailableVotes(ctx context.Context, userID uuid.UUID) (free int, paid int, err error)

	// Analytics
	GetCategoryVoteStats(ctx context.Context, categoryID uuid.UUID) ([]models.VoteCount, error)
	GetNomineeVoteStats(ctx context.Context, nomineeID uuid.UUID) ([]models.VoteCount, error)
	GetUserVoteSummary(ctx context.Context, userID uuid.UUID) ([]models.UserVoteSummary, error)

	// Validation
	CanVoteInCategory(ctx context.Context, userID, categoryID uuid.UUID, usePaidVote bool) error
	ValidateVotingPeriod(ctx context.Context, categoryID uuid.UUID) (bool, error)
}

type votingService struct {
	voteRepo            repositories.VoteRepository
	userRepo            repositories.UserRepository
	categoryRepo        repositories.CategoryRepository
	nomineeRepo         repositories.NomineeRepository
	nomineeCategoryRepo repositories.NomineeCategoryRepository
}

func NewVotingService(
	voteRepo repositories.VoteRepository,
	userRepo repositories.UserRepository,
	categoryRepo repositories.CategoryRepository,
	nomineeRepo repositories.NomineeRepository,
	nomineeCategoryRepo repositories.NomineeCategoryRepository,
) VotingService {
	return &votingService{
		voteRepo:            voteRepo,
		userRepo:            userRepo,
		categoryRepo:        categoryRepo,
		nomineeRepo:         nomineeRepo,
		nomineeCategoryRepo: nomineeCategoryRepo,
	}
}

func (s *votingService) CastVote(ctx context.Context, userID, nomineeID, categoryID uuid.UUID, usePaidVote bool) (*models.Vote, error) {
	// 1. Validate user exists and has votes
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 2. Validate category is active
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	if !category.IsActive {
		return nil, ErrCategoryNotActive
	}

	// 3. Validate nominee is active and in category
	nominee, err := s.nomineeRepo.GetByID(ctx, nomineeID)
	if err != nil {
		return nil, fmt.Errorf("nominee not found: %w", err)
	}
	if !nominee.IsActive {
		return nil, ErrNomineeNotActive
	}

	// Verify nominee is in the category
	nomineesInCategory, err := s.nomineeCategoryRepo.GetNomineesForCategory(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify nominee in category: %w", err)
	}

	nomineeInCategory := false
	for _, n := range nomineesInCategory {
		if n.NomineeID == nomineeID {
			nomineeInCategory = true
			break
		}
	}
	if !nomineeInCategory {
		return nil, ErrNomineeNotInCategory
	}

	// 4. Validate voting period
	isOpen, err := s.ValidateVotingPeriod(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error validating voting period: %w", err)
	}
	if !isOpen {
		return nil, ErrVotingPeriodClosed
	}

	// 5. Check if can vote in category
	if err := s.CanVoteInCategory(ctx, userID, categoryID, usePaidVote); err != nil {
		return nil, err
	}

	// 6. Determine vote type and decrement appropriate vote count
	var voteType models.VoteType
	if usePaidVote {
		if user.PaidVotes <= 0 {
			return nil, ErrNoPaidVotesAvailable
		}
		voteType = models.VoteTypePaid
		if err := s.userRepo.DecrementPaidVotes(ctx, userID); err != nil {
			return nil, fmt.Errorf("failed to decrement paid votes: %w", err)
		}
	} else {
		if user.FreeVotes <= 0 {
			return nil, ErrNoFreeVotesAvailable
		}
		voteType = models.VoteTypeFree
		if err := s.userRepo.DecrementFreeVotes(ctx, userID); err != nil {
			return nil, fmt.Errorf("failed to decrement free votes: %w", err)
		}
	}

	// 7. Create the vote
	vote := &models.Vote{
		VoteID:     uuid.New(),
		UserID:     userID,
		NomineeID:  nomineeID,
		CategoryID: categoryID,
		VoteType:   voteType,
	}

	if err := s.voteRepo.Create(ctx, vote); err != nil {
		// Rollback vote count if vote creation fails
		if voteType == models.VoteTypePaid {
			s.userRepo.IncrementPaidVotes(ctx, userID)
		} else {
			s.userRepo.IncrementFreeVotes(ctx, userID)
		}
		return nil, fmt.Errorf("failed to cast vote: %w", err)
	}

	// Reload vote with relationships
	return s.voteRepo.GetByID(ctx, vote.VoteID)
}

func (s *votingService) CanVoteInCategory(ctx context.Context, userID, categoryID uuid.UUID, usePaidVote bool) error {
	// If using paid vote, can always vote (no restriction per category)
	if usePaidVote {
		return nil
	}

	// For free votes, check if already voted in this category
	freeVoteCount, err := s.voteRepo.CountFreeVotesByUserAndCategory(ctx, userID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to check existing votes: %w", err)
	}

	if freeVoteCount > 0 {
		return ErrAlreadyVotedWithFreeVote
	}

	return nil
}

func (s *votingService) ChangeVote(ctx context.Context, voteID uuid.UUID, newNomineeID uuid.UUID) (*models.Vote, error) {
	// 1. Get existing vote
	vote, err := s.voteRepo.GetByID(ctx, voteID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vote: %w", err)
	}

	// 2. Validate new nominee is active and in same category
	nominee, err := s.nomineeRepo.GetByID(ctx, newNomineeID)
	if err != nil {
		return nil, fmt.Errorf("nominee not found: %w", err)
	}
	if !nominee.IsActive {
		return nil, ErrNomineeNotActive
	}

	// Verify nominee is in the category
	nomineesInCategory, err := s.nomineeCategoryRepo.GetNomineesForCategory(ctx, vote.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify nominee in category: %w", err)
	}

	nomineeInCategory := false
	for _, n := range nomineesInCategory {
		if n.NomineeID == newNomineeID {
			nomineeInCategory = true
			break
		}
	}
	if !nomineeInCategory {
		return nil, ErrNomineeNotInCategory
	}

	// 3. Update vote
	vote.NomineeID = newNomineeID
	if err := s.voteRepo.Update(ctx, vote); err != nil {
		return nil, fmt.Errorf("failed to update vote: %w", err)
	}

	return s.voteRepo.GetByID(ctx, voteID)
}

func (s *votingService) DeleteVote(ctx context.Context, voteID uuid.UUID) error {
	vote, err := s.voteRepo.GetByID(ctx, voteID)
	if err != nil {
		return err
	}

	if err := s.voteRepo.Delete(ctx, voteID); err != nil {
		return err
	}

	// Return vote to user based on type
	if vote.VoteType == models.VoteTypePaid {
		if err := s.userRepo.IncrementPaidVotes(ctx, vote.UserID); err != nil {
			return fmt.Errorf("failed to return paid vote: %w", err)
		}
	} else {
		if err := s.userRepo.IncrementFreeVotes(ctx, vote.UserID); err != nil {
			return fmt.Errorf("failed to return free vote: %w", err)
		}
	}

	return nil
}

func (s *votingService) GetVote(ctx context.Context, voteID uuid.UUID) (*models.Vote, error) {
	return s.voteRepo.GetByID(ctx, voteID)
}

func (s *votingService) GetUserVotes(ctx context.Context, userID uuid.UUID) ([]models.Vote, error) {
	return s.voteRepo.GetByUser(ctx, userID)
}

func (s *votingService) GetAllVotes(ctx context.Context) ([]models.Vote, error) {
	return s.voteRepo.GetAll(ctx)
}

func (s *votingService) GetAvailableVotes(ctx context.Context, userID uuid.UUID) (free int, paid int, err error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 0, 0, err
	}
	return user.FreeVotes, user.PaidVotes, nil
}

func (s *votingService) GetCategoryVoteStats(ctx context.Context, categoryID uuid.UUID) ([]models.VoteCount, error) {
	return s.voteRepo.GetVoteCountsByCategory(ctx, categoryID)
}

func (s *votingService) GetNomineeVoteStats(ctx context.Context, nomineeID uuid.UUID) ([]models.VoteCount, error) {
	return s.voteRepo.GetVoteCountsByNominee(ctx, nomineeID)
}

func (s *votingService) GetUserVoteSummary(ctx context.Context, userID uuid.UUID) ([]models.UserVoteSummary, error) {
	return s.voteRepo.GetUserVoteSummary(ctx, userID)
}

func (s *votingService) ValidateVotingPeriod(ctx context.Context, categoryID uuid.UUID) (bool, error) {
	// TODO: Implement actual voting period validation
	// For now, return true if category is active
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return false, err
	}
	return category.IsActive, nil
}
