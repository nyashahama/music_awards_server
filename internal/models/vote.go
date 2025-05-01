package models

import (
	"time"

	"github.com/google/uuid"
)


type Vote struct {
	VoteID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_vote_user_category"`
	NomineeID   uuid.UUID `gorm:"type:uuid;not null"`
	CategoryID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_vote_user_category"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	
	// Associations
	User         User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Nominee      Nominee        `gorm:"foreignKey:NomineeID;constraint:OnDelete:CASCADE;"`
	Category     Category       `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;"`
	NomineeCategory NomineeCategory `gorm:"foreignKey:NomineeID,CategoryID;references:NomineeID,CategoryID;constraint:OnDelete:CASCADE;"`
}
