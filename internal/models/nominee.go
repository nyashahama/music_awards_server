package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	//"gorm.io/datatypes"
)

type Nominee struct {
	NomineeID   uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string          `gorm:"not null;index"`
	Description string          `gorm:"type:text"`
	SampleWorks json.RawMessage `gorm:"type:jsonb"`
	ImageURL    string          `gorm:"type:varchar(500)"`
	IsActive    bool            `gorm:"not null;default:true"` // Enable/disable nominee
	CreatedAt   time.Time       `gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime"`
	Categories  []Category      `gorm:"many2many:nominee_categories;"`
	Votes       []Vote          `gorm:"foreignKey:NomineeID;constraint:OnDelete:CASCADE"`
}
