// Package models
package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	FirstName           string    `gorm:"not null"`
	LastName            string    `gorm:"not null"`
	PasswordHash        string    `gorm:"not null"`
	Email               string    `gorm:"unique;not null"`
	Role                string    `gorm:"not null"`
	Location            string    `gorm:"type:varchar(255)"` // Store user location/country
	AvailableVotes      int       `gorm:"not null;default:5"`
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
	Votes               []Vote    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ResetToken          *string   `gorm:"type:varchar(255)"`
	ResetTokenExpiresAt *time.Time
}
