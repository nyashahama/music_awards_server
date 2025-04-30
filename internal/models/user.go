package models

import (
	"time"

	"github.com/google/uuid"
)


type User struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username     string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Role         string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}