package models

import (
	"time"

	"github.com/google/uuid"
)

type VoteType string

const (
	VoteTypeFree VoteType = "free"
	VoteTypePaid VoteType = "paid"
)

type Vote struct {
	VoteID     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index:idx_user_votes"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null;index:idx_category_votes"`
	NomineeID  uuid.UUID `gorm:"type:uuid;not null;index:idx_nominee_votes"`
	VoteType   VoteType  `gorm:"type:varchar(10);not null;default:'free'"` // Track if free or paid
	CreatedAt  time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	// Relationships
	User     User     `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE"`
	Category Category `gorm:"foreignKey:CategoryID;references:CategoryID;constraint:OnDelete:CASCADE"`
	Nominee  Nominee  `gorm:"foreignKey:NomineeID;references:NomineeID;constraint:OnDelete:CASCADE"`
}

func (NomineeCategory) TableName() string {
	return "nominee_categories"
}
