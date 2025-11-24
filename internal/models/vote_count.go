package models

import "github.com/google/uuid"

type VoteCount struct {
	NomineeID   uuid.UUID `json:"nominee_id"`
	NomineeName string    `json:"nominee_name"`
	CategoryID  uuid.UUID `json:"category_id"`
	TotalVotes  int64     `json:"total_votes"`
	FreeVotes   int64     `json:"free_votes"`
	PaidVotes   int64     `json:"paid_votes"`
}
