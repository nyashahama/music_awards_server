package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
<<<<<<< Updated upstream
	UserID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username       string    `gorm:"unique;not null"`
	PasswordHash   string    `gorm:"not null"`
	Email          string    `gorm:"unique;not null"`
	Role           string    `gorm:"not null"`
	AvailableVotes int       `gorm:"not null;default:5"` //new
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	Votes          []Vote    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
=======
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
>>>>>>> Stashed changes
}
