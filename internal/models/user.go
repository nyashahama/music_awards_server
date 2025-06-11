package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username       string    `gorm:"unique;not null"`
	PasswordHash   string    `gorm:"not null"`
	Email          string    `gorm:"unique;not null"`
	Role           string    `gorm:"not null"`
	AvailableVotes int       `gorm:"not null;default:5"` //new
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	Votes          []Vote    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
