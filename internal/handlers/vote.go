package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/dtos"
	"github.com/nyashahama/music-awards/internal/middleware"
	"github.com/nyashahama/music-awards/internal/repositories"
	"github.com/nyashahama/music-awards/internal/services"
)

type VoteHandler struct {
	votingService services.VotingService
}

func NewVoteHandler(votingService services.VotingService) *VoteHandler {
	return &VoteHandler{votingService: votingService}
}

func (h *VoteHandler) RegisterRoutes(r *gin.Engine) {
	votes := r.Group("/votes")
	votes.Use(middleware.AuthMiddleware())
	{
		votes.POST("", h.CastVote)
		votes.GET("/me", h.GetMyVotes)
		votes.GET("/me/summary", h.GetMyVoteSummary)
		votes.GET("/me/available", h.GetAvailableVotes)
		votes.PUT("/:id", h.ChangeVote)
		votes.DELETE("/:id", h.DeleteVote)
	}

	// Admin endpoints
	adminVotes := votes.Group("")
	adminVotes.Use(middleware.AdminMiddleware())
	{
		adminVotes.GET("/all", h.GetAllVotes)
		adminVotes.GET("/category/:category_id/stats", h.GetCategoryStats)
		adminVotes.GET("/nominee/:nominee_id/stats", h.GetNomineeStats)
	}
}

// CastVote allows a user to vote for a nominee in a category
func (h *VoteHandler) CastVote(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req dtos.CastVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vote, err := h.votingService.CastVote(c.Request.Context(), userID, req.NomineeID, req.CategoryID, req.UsePaidVote)
	if err != nil {
		handleVotingError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dtos.NewUserVoteResponse(vote))
}

// GetMyVotes returns all votes for the authenticated user
func (h *VoteHandler) GetMyVotes(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	votes, err := h.votingService.GetUserVotes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get votes"})
		return
	}

	response := make([]dtos.UserVoteResponse, len(votes))
	for i, vote := range votes {
		response[i] = dtos.NewUserVoteResponse(&vote)
	}

	c.JSON(http.StatusOK, response)
}

// GetMyVoteSummary returns voting summary by category for the authenticated user
func (h *VoteHandler) GetMyVoteSummary(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	summary, err := h.votingService.GetUserVoteSummary(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get vote summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetAvailableVotes returns the number of available votes for the authenticated user
func (h *VoteHandler) GetAvailableVotes(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	freeVotes, paidVotes, err := h.votingService.GetAvailableVotes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get available votes"})
		return
	}

	response := dtos.AvailableVotesResponse{
		FreeVotes: freeVotes,
		PaidVotes: paidVotes,
		Total:     freeVotes + paidVotes,
	}

	c.JSON(http.StatusOK, response)
}

// ChangeVote allows a user to change their vote to a different nominee in the same category
func (h *VoteHandler) ChangeVote(c *gin.Context) {
	voteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vote ID"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	currentUserRole := c.MustGet("user_role").(string)

	// Get vote to check ownership
	vote, err := h.votingService.GetVote(c.Request.Context(), voteID)
	if err != nil {
		handleVotingError(c, err)
		return
	}

	// Authorization: Only vote owner or admin can change
	if currentUserRole != "admin" && vote.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req dtos.ChangeVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedVote, err := h.votingService.ChangeVote(c.Request.Context(), voteID, req.NomineeID)
	if err != nil {
		handleVotingError(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.NewUserVoteResponse(updatedVote))
}

// DeleteVote allows a user to delete their vote and get their vote back
func (h *VoteHandler) DeleteVote(c *gin.Context) {
	voteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vote ID"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	currentUserRole := c.MustGet("user_role").(string)

	// Get vote to check ownership
	vote, err := h.votingService.GetVote(c.Request.Context(), voteID)
	if err != nil {
		handleVotingError(c, err)
		return
	}

	// Authorization: Only vote owner or admin can delete
	if currentUserRole != "admin" && vote.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.votingService.DeleteVote(c.Request.Context(), voteID); err != nil {
		handleVotingError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllVotes returns all votes (admin only)
func (h *VoteHandler) GetAllVotes(c *gin.Context) {
	votes, err := h.votingService.GetAllVotes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get votes"})
		return
	}

	response := make([]dtos.VoteResponse, len(votes))
	for i, vote := range votes {
		response[i] = dtos.NewVoteResponse(&vote)
	}

	c.JSON(http.StatusOK, response)
}

// GetCategoryStats returns vote statistics for a category (admin only)
func (h *VoteHandler) GetCategoryStats(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	stats, err := h.votingService.GetCategoryVoteStats(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get category stats"})
		return
	}

	response := make([]dtos.VoteStatsResponse, len(stats))
	for i, stat := range stats {
		response[i] = dtos.NewVoteStatsResponse(&stat)
	}

	c.JSON(http.StatusOK, response)
}

// GetNomineeStats returns vote statistics for a nominee across all categories (admin only)
func (h *VoteHandler) GetNomineeStats(c *gin.Context) {
	nomineeID, err := uuid.Parse(c.Param("nominee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nominee ID"})
		return
	}

	stats, err := h.votingService.GetNomineeVoteStats(c.Request.Context(), nomineeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get nominee stats"})
		return
	}

	response := make([]dtos.VoteStatsResponse, len(stats))
	for i, stat := range stats {
		response[i] = dtos.NewVoteStatsResponse(&stat)
	}

	c.JSON(http.StatusOK, response)
}

// handleVotingError maps service errors to appropriate HTTP responses
func handleVotingError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNoFreeVotesAvailable):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No free votes available",
			"code":  "NO_FREE_VOTES",
		})
	case errors.Is(err, services.ErrNoPaidVotesAvailable):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No paid votes available",
			"code":  "NO_PAID_VOTES",
		})
	case errors.Is(err, services.ErrAlreadyVotedWithFreeVote):
		c.JSON(http.StatusConflict, gin.H{
			"error": "Already voted in this category with a free vote. Use a paid vote to vote again.",
			"code":  "ALREADY_VOTED_FREE",
		})
	case errors.Is(err, services.ErrVotingPeriodClosed):
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Voting period is closed for this category",
			"code":  "VOTING_CLOSED",
		})
	case errors.Is(err, services.ErrCategoryNotActive):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Category is not active",
			"code":  "CATEGORY_INACTIVE",
		})
	case errors.Is(err, services.ErrNomineeNotActive):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nominee is not active",
			"code":  "NOMINEE_INACTIVE",
		})
	case errors.Is(err, services.ErrNomineeNotInCategory):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nominee is not in this category",
			"code":  "NOMINEE_NOT_IN_CATEGORY",
		})
	case errors.Is(err, repositories.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Vote not found",
			"code":  "NOT_FOUND",
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "INTERNAL_ERROR",
		})
	}
}
