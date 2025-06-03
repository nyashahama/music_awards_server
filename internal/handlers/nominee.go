package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/middleware"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/services"
	"gorm.io/datatypes"
)

type NomineeHandler struct {
	nomineeService services.NomineeService
}

func NewNomineeHandler(nomineeService services.NomineeService) *NomineeHandler {
	return &NomineeHandler{nomineeService: nomineeService}
}

//will create a DTO later

type createNomineeRequest struct {
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"`
	SampleWorks datatypes.JSON `json:"sample_works"`
	ImageURL    string         `json:"image_url"`
}

func (h *NomineeHandler) RegisterRoutes(r *gin.Engine) {
	//publiv endpoints for nominee
	nominees := r.Group("/nominees")
	nominees.GET("", h.GetAllNominees)
	nominees.GET("/:id", h.GetNomineeDetails)

	//admin/protected
	adminNominees := r.Group("/nominees")
	adminNominees.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		adminNominees.POST("", h.CreateNominee)
		adminNominees.PUT("/:id", h.UpdateNominee)
		adminNominees.DELETE("/:id", h.DeleteNominee)
	}
}

func (h *NomineeHandler) CreateNominee(c *gin.Context) {
	var req createNomineeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nominee := models.Nominee{
		Name:        req.Name,
		Description: req.Description,
		SampleWorks: datatypes.JSON(req.SampleWorks),
		ImageURL:    req.ImageURL,
	}

	result, err := h.nomineeService.CreateNominee(c.Request.Context(), nominee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)

}
func (h *NomineeHandler) UpdateNominee(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nominee, err := h.nomineeService.UpdateNominee(c.Request.Context(), id, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nominee)
}

func (h *NomineeHandler) DeleteNominee(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.nomineeService.DeleteNominee(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nominee)
}

func (h *NomineeHandler) GetAllNominees(c *gin.Context) {
	nominees, err := h.nomineeService.GetAllNominees(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nominees)
}

func (h *NomineeHandler) SetNominationPeriod(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
