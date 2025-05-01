package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)


type Nominee struct {
	NomineeID   uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string         `gorm:"not null"`
	Description string
	SampleWorks datatypes.JSON `gorm:"type:jsonb"`
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Categories  []Category     `gorm:"many2many:nominee_categories;"`
}