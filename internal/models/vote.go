package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
    VoteID     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
    UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
    NomineeID  uuid.UUID `gorm:"type:uuid;not null;index"`
    CategoryID uuid.UUID `gorm:"type:uuid;not null;index"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`

    // enforce one vote per (user,category)
    // this creates a composite unique index
    _ struct{} `gorm:"uniqueIndex:idx_vote_user_category,columns:user_id,category_id"`

    // associations
    User     User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
    Nominee  Nominee  `gorm:"foreignKey:NomineeID;constraint:OnDelete:CASCADE;"`
    Category Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;"`
}
