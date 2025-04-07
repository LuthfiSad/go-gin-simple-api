package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	ID               uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Title            string            `gorm:"size:255;not null" json:"title"`
	Description      string            `gorm:"type:text" json:"description"`
	CoverID          *uuid.UUID        `json:"cover_id"`
	Cover            *Media            `gorm:"foreignKey:CoverID" json:"cover,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	DeletedAt        gorm.DeletedAt    `gorm:"index" json:"-"`
	BookStocks       []BookStock       `gorm:"foreignKey:BookID" json:"book_stocks,omitempty"`
	BookTransactions []BookTransaction `gorm:"foreignKey:BookID" json:"book_transactions,omitempty"`
}
