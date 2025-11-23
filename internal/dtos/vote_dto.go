package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
)

type CastVoteRequest struct {
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	NomineeID   uuid.UUID `json:"nominee_id" binding:"required"`
	UsePaidVote bool      `json:"use_paid_vote"` // false = use free vote, true = use paid vote
}

type ChangeVoteRequest struct {
	NomineeID uuid.UUID `json:"nominee_id" binding:"required"`
}

type VoteResponse struct {
	VoteID     uuid.UUID `json:"vote_id"`
	UserID     uuid.UUID `json:"user_id"`
	CategoryID uuid.UUID `json:"category_id"`
	NomineeID  uuid.UUID `json:"nominee_id"`
	VoteType   string    `json:"vote_type"` // "free" or "paid"
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserVoteResponse struct {
	VoteID    uuid.UUID     `json:"vote_id"`
	VoteType  string        `json:"vote_type"`
	Category  CategoryBrief `json:"category"`
	Nominee   NomineeBrief  `json:"nominee"`
	CreatedAt time.Time     `json:"created_at"`
}

type AvailableVotesResponse struct {
	FreeVotes int `json:"free_votes"`
	PaidVotes int `json:"paid_votes"`
	Total     int `json:"total"`
}

type VoteStatsResponse struct {
	NomineeID   uuid.UUID `json:"nominee_id"`
	NomineeName string    `json:"nominee_name"`
	CategoryID  uuid.UUID `json:"category_id"`
	TotalVotes  int64     `json:"total_votes"`
	FreeVotes   int64     `json:"free_votes"`
	PaidVotes   int64     `json:"paid_votes"`
}

type UserVoteSummaryResponse struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	FreeVotes    int64     `json:"free_votes_used"`
	PaidVotes    int64     `json:"paid_votes_used"`
	TotalVotes   int64     `json:"total_votes"`
}

// Response builders
func NewVoteResponse(vote *models.Vote) VoteResponse {
	return VoteResponse{
		VoteID:     vote.VoteID,
		UserID:     vote.UserID,
		CategoryID: vote.CategoryID,
		NomineeID:  vote.NomineeID,
		VoteType:   string(vote.VoteType),
		CreatedAt:  vote.CreatedAt,
		UpdatedAt:  vote.UpdatedAt,
	}
}

func NewUserVoteResponse(vote *models.Vote) UserVoteResponse {
	return UserVoteResponse{
		VoteID:   vote.VoteID,
		VoteType: string(vote.VoteType),
		Category: CategoryBrief{
			CategoryID: vote.Category.CategoryID,
			Name:       vote.Category.Name,
		},
		Nominee: NomineeBrief{
			NomineeID: vote.Nominee.NomineeID,
			Name:      vote.Nominee.Name,
			ImageURL:  vote.Nominee.ImageURL,
		},
		CreatedAt: vote.CreatedAt,
	}
}

func NewVoteStatsResponse(voteCount *models.VoteCount) VoteStatsResponse {
	return VoteStatsResponse{
		NomineeID:   voteCount.NomineeID,
		NomineeName: voteCount.NomineeName,
		CategoryID:  voteCount.CategoryID,
		TotalVotes:  voteCount.TotalVotes,
		FreeVotes:   voteCount.FreeVotes,
		PaidVotes:   voteCount.PaidVotes,
	}
}
