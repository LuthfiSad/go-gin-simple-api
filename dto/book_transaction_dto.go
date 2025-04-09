package dto

import (
	"time"

	"github.com/google/uuid"
)

type BookTransactionResponse struct {
	ID         uuid.UUID          `json:"id"`
	BookID     uuid.UUID          `json:"book_id"`
	Book       *BookRes           `json:"book,omitempty"`
	StockCode  string             `json:"stock_code"`
	BookStock  *BookStockResponse `json:"book_stock,omitempty"`
	CustomerID uuid.UUID          `json:"customer_id"`
	Customer   *CustomerResponse  `json:"customer,omitempty"`
	DueDate    time.Time          `json:"due_date"`
	Status     string             `json:"status"`
	BorrowedAt *time.Time         `json:"borrowed_at,omitempty"`
	ReturnAt   *time.Time         `json:"return_at,omitempty"`
	Charges    []ChargeResponse   `json:"charges,omitempty"`
}

type BookTransactionCreateRequest struct {
	// BookID     uuid.UUID `json:"book_id" validate:"required"`
	StockCode  string    `json:"stock_code" validate:"required"`
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	Status     string    `json:"status" validate:"required,oneof=Borrowed Returned Overdue"`
}

type BookTransactionUpdateRequest struct {
	// BookID     uuid.UUID  `json:"book_id,omitempty"`
	StockCode  string     `json:"stock_code,omitempty"`
	CustomerID uuid.UUID  `json:"customer_id,omitempty"`
	DueDate    *time.Time `json:"due_date,omitempty"`
	Status     string     `json:"status,omitempty" validate:"omitempty,oneof=Borrowed Returned Overdue"`
	ReturnAt   *time.Time `json:"return_at,omitempty"`
}

type BookTransactionStatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=Borrowed Returned Overdue"`
}

type BookTransactionReturnRequest struct {
	ReturnAt *time.Time `json:"return_at,omitempty"`
}
