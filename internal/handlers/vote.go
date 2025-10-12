package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/dtos"
	"github.com/nyashahama/music-awards/internal/middleware"
	"github.com/nyashahama/music-awards/internal/services"
	"gorm.io/gorm"
)

type VoteHandler struct {
	voteService services.VotingMechanismService
}

func NewVoteHandler(voteService services.VotingMechanismService) *VoteHandler {
	return &VoteHandler{voteService: voteService}
}

func (h *VoteHandler) RegisterRoutes(r *gin.Engine) {
	votes := r.Group("/votes")
	votes.Use(middleware.AuthMiddleware())
	{
		votes.POST("", h.CastVote)
		votes.GET("", h.GetUserVotes)
		votes.GET("/available", h.GetAvailableVotes)
		votes.PUT("/:id", h.ChangeVote)
		votes.DELETE("/:id", h.DeleteVote)
	}

	admin := votes.Group("")
	admin.Use(middleware.AdminMiddleware())
	{
		admin.GET("/category/:category_id", h.GetCategoryVotes)
		admin.GET("/all", h.GetAllVotes)
	}
}

func (h *VoteHandler) CastVote(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	fmt.Printf("Authenticated userId %s\n", userID)

	var req dtos.CastVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vote, err := h.voteService.CastVote(c.Request.Context(), userID, req.NomineeID, req.CategoryID)
	if err != nil {
		handleVoteServiceError(c, err)
		return
	}

	// Use NewUserVotesResponse instead of NewVoteResponse to get category/nominee details
	c.JSON(http.StatusCreated, dtos.NewUserVotesResponse(vote))
}

func (h *VoteHandler) GetUserVotes(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	votes, err := h.voteService.GetUserVotes(c.Request.Context(), userID)
	if err != nil {
		handleVoteServiceError(c, err)
		return
	}

	response := make([]dtos.UserVotesResponse, len(votes))
	for i, vote := range votes {
		response[i] = dtos.NewUserVotesResponse(&vote)
	}
	c.JSON(http.StatusOK, response)
}

func (h *VoteHandler) GetAvailableVotes(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	availableVotes, err := h.voteService.GetAvailableVotes(c.Request.Context(), userID)
	if err != nil {
		handleVoteServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"available_votes": availableVotes})
}

func (h *VoteHandler) ChangeVote(c *gin.Context) {
	voteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vote ID"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	currentUserRole := c.MustGet("user_role").(string)

	// Get vote to check ownership
	vote, err := h.voteService.GetVote(c.Request.Context(), voteID)
	if err != nil {
		handleVoteServiceError(c, err)
		return
	}

	if currentUserRole != "admin" && vote.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req struct {
		NomineeID uuid.UUID `json:"nominee_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedVote, err := h.voteService.ChangeVote(c.Request.Context(), voteID, req.NomineeID)
	if err != nil {
		handleVoteServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.NewVoteResponse(updatedVote))
}

func (h *VoteHandler) DeleteVote(c *gin.Context) {
	voteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vote ID"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	currentUserRole := c.MustGet("user_role").(string)

	// Get vote to check ownership
	vote, err := h.voteService.GetVote(c.Request.Context(), voteID)
	if err != nil {
		handleVoteServiceError(c, err)
		return
	}

	// Authorization: Only vote owner or admin can delete
	if currentUserRole != "admin" && vote.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.voteService.DeleteVote(c.Request.Context(), voteID); err != nil {
		handleVoteServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *VoteHandler) GetCategoryVotes(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	votes, err := h.voteService.GetCategoryVotes(c.Request.Context(), categoryID)
	if err != nil {
		handleVoteServiceError(c, err)
		return
	}

	response := make([]dtos.VoteResponse, len(votes))
	for i, vote := range votes {
		response[i] = dtos.NewVoteResponse(&vote)
	}
	c.JSON(http.StatusOK, response)
}

func (h *VoteHandler) GetAllVotes(c *gin.Context) {
	votes, err := h.voteService.GetAllVotes(c.Request.Context())
	if err != nil {
		handleVoteServiceError(c, err)
		return
	}

	response := make([]dtos.VoteResponse, len(votes))
	for i, vote := range votes {
		response[i] = dtos.NewVoteResponse(&vote)
	}
	c.JSON(http.StatusOK, response)
}

func handleVoteServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNoVotesAvailable):
		c.JSON(http.StatusBadRequest, gin.H{"error": "no votes available"})
	case errors.Is(err, services.ErrAlreadyVotedInCategory):
		c.JSON(http.StatusConflict, gin.H{"error": "already voted in this category"})
	case errors.Is(err, services.ErrVotingPeriodClosed):
		c.JSON(http.StatusForbidden, gin.H{"error": "voting period is closed"})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "vote not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
