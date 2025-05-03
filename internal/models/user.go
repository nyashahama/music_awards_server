package models

import (
	"time"

	"github.com/google/uuid"
)


type User struct {
    UserID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
    Username     string    `gorm:"unique;not null"`  // Changed from uniqueIndex to unique
    PasswordHash string    `gorm:"not null"`
    Email        string    `gorm:"unique;not null"`   // Changed from uniqueIndex to unique
    Role         string    `gorm:"not null"`
    CreatedAt    time.Time `gorm:"autoCreateTime"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime"`
    Votes []Vote `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}