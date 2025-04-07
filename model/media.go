package model

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Path      string    `gorm:"size:255;not null" json:"path"`
	PublicID  string    `gorm:"size:255;not null" json:"public_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Books     []Book    `gorm:"foreignKey:CoverID" json:"books,omitempty"`
}
