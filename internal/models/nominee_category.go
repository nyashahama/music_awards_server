package models

import (
	"time"

	"github.com/google/uuid"
)

type NomineeCategory struct {
	NomineeID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	CategoryID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	// Associations
	Nominee  Nominee  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Category Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Nominee    Nominee   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Category   Category  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
