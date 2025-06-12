// handlers/nominee.go
package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/dtos"
	"github.com/nyashahama/music-awards/internal/middleware"
	"github.com/nyashahama/music-awards/internal/services"
	"gorm.io/gorm"
)

type NomineeHandler struct {
	nomineeService services.NomineeService
}

func NewNomineeHandler(nomineeService services.NomineeService) *NomineeHandler {
	return &NomineeHandler{nomineeService: nomineeService}
}

func (h *NomineeHandler) RegisterRoutes(r *gin.Engine) {
	public := r.Group("/nominees")
	{
		public.GET("", h.GetAllNominees)
		public.GET("/:id", h.GetNomineeDetails)
	}

	admin := r.Group("/nominees")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.POST("", h.CreateNominee)
		admin.PUT("/:id", h.UpdateNominee)
		admin.DELETE("/:id", h.DeleteNominee)
	}
}

func (h *NomineeHandler) CreateNominee(c *gin.Context) {
	var req dtos.CreateNomineeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nominee, err := h.nomineeService.CreateNominee(c.Request.Context(), req)
	if err != nil {
		handleNomineeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dtos.NewNomineeResponse(nominee))
}

func (h *NomineeHandler) UpdateNominee(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req dtos.UpdateNomineeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nominee, err := h.nomineeService.UpdateNominee(c.Request.Context(), id, req)
	if err != nil {
		handleNomineeError(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.NewNomineeResponse(nominee))
}

func (h *NomineeHandler) DeleteNominee(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.nomineeService.DeleteNominee(c.Request.Context(), id); err != nil {
		handleNomineeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *NomineeHandler) GetNomineeDetails(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	nominee, err := h.nomineeService.GetNomineeDetails(c.Request.Context(), id)
	if err != nil {
		handleNomineeError(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.NewNomineeResponse(nominee))
}

func (h *NomineeHandler) GetAllNominees(c *gin.Context) {
	nominees, err := h.nomineeService.GetAllNominees(c.Request.Context())
	if err != nil {
		handleNomineeError(c, err)
		return
	}

	response := make([]dtos.NomineeResponse, len(nominees))
	for i, nominee := range nominees {
		response[i] = dtos.NewNomineeResponse(&nominee)
	}
	c.JSON(http.StatusOK, response)
}

func handleNomineeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNomineeNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrInvalidJSON):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
	case errors.Is(err, services.ErrCategoryNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"error": "one or more categories not found"})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
