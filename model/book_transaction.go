package model

import (
	"time"

	"github.com/google/uuid"
)

type BookTransaction struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	BookID     uuid.UUID  `gorm:"not null" json:"book_id"`
	Book       Book       `gorm:"foreignKey:BookID" json:"book,omitempty"`
	StockCode  string     `gorm:"size:50;not null" json:"stock_code"`
	BookStock  BookStock  `gorm:"foreignKey:StockCode;references:Code" json:"book_stock,omitempty"`
	CustomerID uuid.UUID  `gorm:"not null" json:"customer_id"`
	Customer   Customer   `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	DueDate    time.Time  `json:"due_date"`
	Status     string     `gorm:"size:50;not null" json:"status"` // Borrowed, Returned, Overdue
	BorrowedAt *time.Time `json:"borrowed_at"`
	ReturnAt   *time.Time `json:"return_at"`
	Charges    []Charge   `gorm:"foreignKey:BookTransactionID" json:"charges,omitempty"`
}

const (
	StatusBTBorrowed = "Borrowed"
	StatusBTReturned = "Returned"
	StatusBTOverdue  = "Overdue"
)
