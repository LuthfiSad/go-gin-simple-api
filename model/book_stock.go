package model

import (
	"github.com/google/uuid"
)

type BookStock struct {
	Code             string            `gorm:"primaryKey;size:50" json:"code"`
	BookID           uuid.UUID         `gorm:"not null" json:"book_id"`
	Book             Book              `gorm:"foreignKey:BookID" json:"book"`
	Status           string            `gorm:"size:50;not null" json:"status"` // Available, Borrowed, Damaged, Lost
	BookTransactions []BookTransaction `gorm:"foreignKey:StockCode;references:Code" json:"book_transactions,omitempty"`
}

const (
	StatusAvailable = "Available"
	StatusBorrowed  = "Borrowed"
	StatusDamaged   = "Damaged"
	StatusLost      = "Lost"
)
