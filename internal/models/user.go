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
	Email               string    `gorm:"unique;not null;index"`
	Role                string    `gorm:"not null;default:'user'"`
	Location            string    `gorm:"type:varchar(255)"`
	FreeVotes           int       `gorm:"not null;default:3"` // 3 free votes at registration
	PaidVotes           int       `gorm:"not null;default:0"` // Purchased votes
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
	Votes               []Vote    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ResetToken          *string   `gorm:"type:varchar(255);index"`
	ResetTokenExpiresAt *time.Time
}
