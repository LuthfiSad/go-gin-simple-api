// handler/book_stock_handler.go

package handler

import (
	"go-gin-simple-api/dto"
	"go-gin-simple-api/lib"
	"go-gin-simple-api/service"
	"go-gin-simple-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookStockHandler struct {
	bookStockService service.BookStockService
}

func NewBookStockHandler(bookStockService service.BookStockService) *BookStockHandler {
	return &BookStockHandler{
		bookStockService: bookStockService,
	}
}

// ListBookStocks handles retrieving all book stocks with pagination, search, and filter
func (h *BookStockHandler) ListBookStocks(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")
	filterStr := c.Query("filter")

	// Parse filters
	filters := lib.ParseFilterString(filterStr)

	// Get book stocks
	result, err := h.bookStockService.ListBookStocks(page, perPage, search, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve book stocks",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetBookStock handles retrieving a book stock by its code
func (h *BookStockHandler) GetBookStock(c *gin.Context) {
	code := c.Param("code")

	bookStock, err := h.bookStockService.GetBookStock(code)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ResponseError{
			Status:  http.StatusNotFound,
			Message: "Book stock not found",
			Error:   map[string]string{"error": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Book stock retrieved successfully",
		Data:    bookStock,
	})
}

// GetByBookID handles retrieving book stocks by book ID
func (h *BookStockHandler) GetByBookID(c *gin.Context) {
	idStr := c.Param("book_id")

	bookID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid book ID format",
			Error:   map[string]string{"error": err.Error()}})
		return
	}

	bookStocks, err := h.bookStockService.GetByBookID(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve book stocks",
			Error:   map[string]string{"error": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Book stocks retrieved successfully",
		Data:    bookStocks,
	})
}

// GetAvailableByBookID handles retrieving available book stocks by book ID
func (h *BookStockHandler) GetAvailableByBookID(c *gin.Context) {
	idStr := c.Param("book_id")

	bookID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid book ID format",
			Error:   map[string]string{"error": err.Error()}})
		return
	}

	bookStocks, err := h.bookStockService.GetAvailableByBookID(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve available book stocks",
			Error:   map[string]string{"error": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Available book stocks retrieved successfully",
		Data:    bookStocks,
	})
}

// CreateBookStock handles creating a new book stock
func (h *BookStockHandler) CreateBookStock(c *gin.Context) {
	var req dto.BookStockCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   map[string]string{"error": err.Error()}})
		return
	}

	// Validate request
	if validationErrors := utils.Validate(req); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Error:   validationErrors,
		})
		return
	}

	_, err := h.bookStockService.CreateBookStock(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create book stock",
			Error:   map[string]string{"error": err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseSuccess{
		Status:  http.StatusCreated,
		Message: "Book stock created successfully",
	})
}

// UpdateBookStock handles updating a book stock
// UpdateBookStock handles updating a book stock
func (h *BookStockHandler) UpdateBookStock(c *gin.Context) {
	code := c.Param("code")

	var req dto.BookStockUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	// Validate request
	if validationErrors := utils.Validate(req); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Error:   validationErrors,
		})
		return
	}

	bookStock, err := h.bookStockService.UpdateBookStock(code, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update book stock",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Book stock updated successfully",
		Data:    bookStock,
	})
}

// DeleteBookStock handles deleting a book stock
func (h *BookStockHandler) DeleteBookStock(c *gin.Context) {
	code := c.Param("code")

	err := h.bookStockService.DeleteBookStock(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete book stock",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Status:  http.StatusOK,
		Message: "Book stock deleted successfully",
	})
}

// UpdateStatus handles updating a book stock's status
func (h *BookStockHandler) UpdateStatus(c *gin.Context) {
	code := c.Param("code")

	var req dto.BookStockStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	// Validate request
	if validationErrors := utils.Validate(req); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Error:   validationErrors,
		})
		return
	}

	// userData, exists := c.Get("userData")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, dto.ResponseError{Status: http.StatusUnauthorized, Message: "User data not found"})
	// 	return
	// }

	// user, ok := userData.(dto.UserData)
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, dto.ResponseError{Status: http.StatusInternalServerError, Message: "Invalid user data"})
	// 	return
	// }

	bookStock, err := h.bookStockService.UpdateStatus(code, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update book stock status",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Book stock status updated successfully",
		Data:    bookStock,
	})
}
