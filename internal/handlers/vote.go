package handlers

import "net/http"

// VotingHandler handles voting operations
type VotingHandler interface {
	CastVote(w http.ResponseWriter, r *http.Request)
	GetVote(w http.ResponseWriter, r *http.Request)
	ChangeVote(w http.ResponseWriter, r *http.Request)
	GetUserVotes(w http.ResponseWriter, r *http.Request)
	GetCategoryVotes(w http.ResponseWriter, r *http.Request)
	ValidateVotingPeriod(w http.ResponseWriter, r *http.Request)
}

type votingHandler struct {
	// Dependencies
}