package models

import (
	"time"

	"github.com/google/uuid"
)

type NomineeCategory struct {
	NomineeID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	CategoryID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
<<<<<<< Updated upstream
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	// Corrected associations
	Nominee  Nominee  `gorm:"foreignKey:NomineeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Category Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
=======
	Nominee    Nominee   `gorm:"foreignKey:NomineeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Category   Category  `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
>>>>>>> Stashed changes
}
