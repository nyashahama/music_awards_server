package dtos

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
)

type CreateNomineeRequest struct {
	Name        string          `json:"name" binding:"required,min=2"`
	Description string          `json:"description" binding:"max=500"`
	SampleWorks json.RawMessage `json:"sample_works"`
	ImageURL    string          `json:"image_url" binding:"omitempty,url"`
	CategoryIDs []uuid.UUID     `json:"category_ids" binding:"omitempty,dive,uuid"`
}

type UpdateNomineeRequest struct {
	Name        *string          `json:"name" binding:"omitempty,min=2"`
	Description *string          `json:"description" binding:"omitempty,max=500"`
	SampleWorks *json.RawMessage `json:"sample_works"`
	ImageURL    *string          `json:"image_url" binding:"omitempty,url"`
	CategoryIDs *[]uuid.UUID     `json:"category_ids" binding:"omitempty,dive,uuid"`
}

type SetCategoriesRequest struct {
	CategoryIDs []uuid.UUID `json:"category_ids" binding:"required"`
}

type NomineeResponse struct {
	NomineeID   uuid.UUID       `json:"nominee_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	SampleWorks json.RawMessage `json:"sample_works"`
	ImageURL    string          `json:"image_url"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Categories  []CategoryBrief `json:"categories,omitempty"` // Added omitempty
}

type CategoryBrief struct {
	CategoryID uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
}

type NomineeBrief struct {
	NomineeID uuid.UUID `json:"nominee_id"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
}

func NewNomineeResponse(nominee *models.Nominee) NomineeResponse {
	var categories []CategoryBrief
	if nominee.Categories != nil {
		categories = make([]CategoryBrief, len(nominee.Categories))
		for i, cat := range nominee.Categories {
			categories[i] = CategoryBrief{
				CategoryID: cat.CategoryID,
				Name:       cat.Name,
			}
		}
	}

	return NomineeResponse{
		NomineeID:   nominee.NomineeID,
		Name:        nominee.Name,
		Description: nominee.Description,
		SampleWorks: nominee.SampleWorks,
		ImageURL:    nominee.ImageURL,
		CreatedAt:   nominee.CreatedAt,
		UpdatedAt:   nominee.UpdatedAt,
		Categories:  categories,
	}
}

func NewNomineeBrief(nominee *models.Nominee) NomineeBrief {
	return NomineeBrief{
		NomineeID: nominee.NomineeID,
		Name:      nominee.Name,
		ImageURL:  nominee.ImageURL,
	}
}
