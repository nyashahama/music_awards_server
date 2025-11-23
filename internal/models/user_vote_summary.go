package models

import "github.com/google/uuid"

// UserVoteSummary tracks how many votes a user has cast per category
type UserVoteSummary struct {
	UserID     uuid.UUID `json:"user_id"`
	CategoryID uuid.UUID `json:"category_id"`
	FreeVotes  int64     `json:"free_votes_used"`
	PaidVotes  int64     `json:"paid_votes_used"`
}
