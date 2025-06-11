package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	CategoryID  uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"unique;not null"`
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Votes       []Vote    `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;"`
	Nominees    []Nominee `gorm:"many2many:nominee_categories;"`
}
