package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID               uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Code             string            `gorm:"size:50;not null;unique" json:"code"`
	Name             string            `gorm:"size:255;not null" json:"name"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	DeletedAt        gorm.DeletedAt    `gorm:"index" json:"-"`
	BookTransactions []BookTransaction `gorm:"foreignKey:CustomerID" json:"book_transactions,omitempty"`
}
