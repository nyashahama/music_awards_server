package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	//"gorm.io/datatypes"
)

type Nominee struct {
	NomineeID   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"not null"`
	Description string
	SampleWorks json.RawMessage `gorm:"type:jsonb"`
	ImageURL    string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Categories []Category `gorm:"many2many:nominee_categories;joinForeignKey:NomineeID;joinReferences:CategoryID"`
}
