package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	VoteID     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null"`
	NomineeID  uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	User     User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Category Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	Nominee  Nominee  `gorm:"foreignKey:NomineeID;constraint:OnDelete:CASCADE"`
}
