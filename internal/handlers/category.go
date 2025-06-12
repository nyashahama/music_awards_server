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

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) RegisterRoutes(r *gin.Engine) {
	// Public category endpoints
	categories := r.Group("/categories")
	categories.GET("", h.ListCategories)
	categories.GET("/active", h.ListActiveCategories)
	categories.GET("/:categoryId", h.GetCategory)

	//admin
	adminCategories := r.Group("/categories")
	adminCategories.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	adminCategories.POST("", h.CreateCategory)
	adminCategories.PUT("/:categoryId", h.UpdateCategory)
	adminCategories.DELETE("/:categoryId", h.DeleteCategory)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dtos.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.CreateCategory(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dtos.NewCategoryResponse(category))
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var req dtos.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle pointer fields
	name := ""
	if req.Name != nil {
		name = *req.Name
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	category, err := h.categoryService.UpdateCategory(c.Request.Context(), categoryID, name, description)
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.NewCategoryResponse(category))
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), categoryID); err != nil {
		handleCategoryError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	category, err := h.categoryService.GetCategoryDetails(c.Request.Context(), categoryID)
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.NewCategoryResponse(category))
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.categoryService.ListAllCategories(c.Request.Context())
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	response := make([]dtos.CategoryResponse, len(categories))
	for i, cat := range categories {
		response[i] = dtos.NewCategoryResponse(&cat)
	}
	c.JSON(http.StatusOK, response)
}

func (h *CategoryHandler) ListActiveCategories(c *gin.Context) {
	categories, err := h.categoryService.ListActiveCategories(c.Request.Context())
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, categories)
}

func handleCategoryError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrCategoryNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrCategoryExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
