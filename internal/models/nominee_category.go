package models

import (
	"time"

	"github.com/google/uuid"
)

type NomineeCategory struct {
	NomineeID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	CategoryID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	Nominee    Nominee   `gorm:"foreignKey:NomineeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Category   Category  `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
