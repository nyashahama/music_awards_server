package models

import "github.com/google/uuid"


type NomineeCategory struct {
	NomineeID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	CategoryID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Nominee    Nominee   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Category   Category  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
