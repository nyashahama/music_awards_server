package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/middleware"
	"github.com/nyashahama/music-awards/internal/services"
	"gorm.io/gorm"
)

type NomineeCategoryHandler struct {
	service services.NomineeCategoryService
}

func NewNomineeCategoryHandler(service services.NomineeCategoryService) *NomineeCategoryHandler {
	return &NomineeCategoryHandler{service: service}
}

func (h *NomineeCategoryHandler) RegisterRoutes(r *gin.Engine) {
	nomineeCategoryGroup := r.Group("/nominees/:nominee_id/categories")
	nomineeCategoryGroup.Use(middleware.AuthMiddleware())
	{
		nomineeCategoryGroup.POST("", h.AddCategory)
		nomineeCategoryGroup.DELETE("/:category_id", h.RemoveCategory)
		nomineeCategoryGroup.PUT("", h.SetCategories)
		nomineeCategoryGroup.GET("", h.GetCategories)
	}

	categoryGroup := r.Group("/categories/:category_id/nominees")
	categoryGroup.Use(middleware.AuthMiddleware())
	{
		categoryGroup.GET("", h.GetNominees)
	}
}

func (h *NomineeCategoryHandler) AddCategory(c *gin.Context) {
	nomineeID, err := uuid.Parse(c.Param("nominee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nominee ID"})
		return
	}

	categoryID, err := uuid.Parse(c.Param("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	if err := h.service.AddCategory(c.Request.Context(), nomineeID, categoryID); err != nil {
		handleNomineeCategoryError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *NomineeCategoryHandler) RemoveCategory(c *gin.Context) {
	nomineeID, err := uuid.Parse(c.Param("nominee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nominee ID"})
		return
	}

	categoryID, err := uuid.Parse(c.Param("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	if err := h.service.RemoveCategory(c.Request.Context(), nomineeID, categoryID); err != nil {
		handleNomineeCategoryError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *NomineeCategoryHandler) SetCategories(c *gin.Context) {
	nomineeID, err := uuid.Parse(c.Param("nominee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nominee ID"})
		return
	}

	var req struct {
		CategoryIDs []uuid.UUID `json:"category_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetCategories(c.Request.Context(), nomineeID, req.CategoryIDs); err != nil {
		handleNomineeCategoryError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *NomineeCategoryHandler) GetCategories(c *gin.Context) {
	nomineeID, err := uuid.Parse(c.Param("nominee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nominee ID"})
		return
	}

	categories, err := h.service.GetCategories(c.Request.Context(), nomineeID)
	if err != nil {
		handleNomineeCategoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *NomineeCategoryHandler) GetNominees(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	nominees, err := h.service.GetNominees(c.Request.Context(), categoryID)
	if err != nil {
		handleNomineeCategoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, nominees)
}

func handleNomineeCategoryError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	switch {
	case errors.Is(err, services.ErrInvalidID):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

