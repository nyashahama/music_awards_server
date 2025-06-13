package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
)

// CastVoteRequest represents the payload for casting a vote
type CastVoteRequest struct {
	CategoryID uuid.UUID `json:"category_id" binding:"required"`
	NomineeID  uuid.UUID `json:"nominee_id" binding:"required"`
}

// VoteResponse represents a vote in API responses
type VoteResponse struct {
	VoteID     uuid.UUID `json:"vote_id"`
	UserID     uuid.UUID `json:"user_id"`
	CategoryID uuid.UUID `json:"category_id"`
	NomineeID  uuid.UUID `json:"nominee_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// UserVotesResponse represents a user's votes with category/nominee details
type UserVotesResponse struct {
	VoteID    uuid.UUID       `json:"vote_id"`
	Category  CategoryDetails `json:"category"`
	Nominee   NomineeDetails  `json:"nominee"`
	CreatedAt time.Time       `json:"created_at"`
}

type CategoryDetails struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type NomineeDetails struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// NewVoteResponse converts a models.Vote to VoteResponse
func NewVoteResponse(vote *models.Vote) VoteResponse {
	return VoteResponse{
		VoteID:     vote.VoteID,
		UserID:     vote.UserID,
		CategoryID: vote.CategoryID,
		NomineeID:  vote.NomineeID,
		CreatedAt:  vote.CreatedAt,
	}
}

// NewUserVotesResponse converts a models.Vote to UserVotesResponse
func NewUserVotesResponse(vote *models.Vote) UserVotesResponse {
	return UserVotesResponse{
		VoteID: vote.VoteID,
		Category: CategoryDetails{
			ID:   vote.Category.CategoryID,
			Name: vote.Category.Name,
		},
		Nominee: NomineeDetails{
			ID:   vote.Nominee.NomineeID,
			Name: vote.Nominee.Name,
		},
		CreatedAt: vote.CreatedAt,
	}
}
