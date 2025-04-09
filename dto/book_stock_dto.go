package dto

import (
	"github.com/google/uuid"
)

// BookStockResponse represents the response for book stock data
type BookStockResponse struct {
	Code   string    `json:"code"`
	BookID uuid.UUID `json:"book_id"`
	Book   *BookRes  `json:"book,omitempty"`
	Status string    `json:"status"`
}

// BookStockCreateRequest represents the request to create a book stock
type BookStockCreateRequest struct {
	Code   string    `json:"code" validate:"required,min=3,max=50"`
	BookID uuid.UUID `json:"book_id" validate:"required"`
	Status string    `json:"status" validate:"omitempty,oneof=Available Borrowed Damaged Lost"`
}

// BookStockUpdateRequest represents the request to update a book stock
type BookStockUpdateRequest struct {
	BookID uuid.UUID `json:"book_id,omitempty"`
	Status string    `json:"status,omitempty" validate:"omitempty,oneof=Available Borrowed Damaged Lost"`
}

// BookStockStatusUpdateRequest represents the request to update a book stock's status
type BookStockStatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=Available Borrowed Damaged Lost"`
}
