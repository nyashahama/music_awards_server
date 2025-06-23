// Package dtos
package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
)

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=50"`
	Description string `json:"description" binding:"max=255"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=3,max=50"`
	Description *string `json:"description" binding:"omitempty,max=255"`
}

type CategoryResponse struct {
	CategoryID  uuid.UUID `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewCategoryResponse model response
func NewCategoryResponse(category *models.Category) CategoryResponse {
	return CategoryResponse{
		CategoryID:  category.CategoryID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}
