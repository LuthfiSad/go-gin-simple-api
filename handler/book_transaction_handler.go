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

type BookTransactionHandler struct {
	bookTransactionService service.BookTransactionService
}

func NewBookTransactionHandler(bookTransactionService service.BookTransactionService) *BookTransactionHandler {
	return &BookTransactionHandler{
		bookTransactionService: bookTransactionService,
	}
}

// ListBookTransactions handles retrieving all book transactions with pagination, search, and filter
func (h *BookTransactionHandler) ListBookTransactions(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")
	filterStr := c.Query("filter")

	// Parse filters
	filters := lib.ParseFilterString(filterStr)

	// Get book transactions
	result, err := h.bookTransactionService.ListBookTransactions(page, perPage, search, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve book transactions",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetBookTransaction handles retrieving a book transaction by ID
func (h *BookTransactionHandler) GetBookTransaction(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid transaction ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	transaction, err := h.bookTransactionService.GetBookTransaction(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ResponseError{
			Status:  http.StatusNotFound,
			Message: "Book transaction not found",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Book transaction retrieved successfully",
		Data:    transaction,
	})
}

// GetByCustomerID handles retrieving book transactions by customer ID
func (h *BookTransactionHandler) GetByCustomerID(c *gin.Context) {
	idStr := c.Param("customer_id")

	customerID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid customer ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	transactions, err := h.bookTransactionService.GetByCustomerID(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve book transactions",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Book transactions retrieved successfully",
		Data:    transactions,
	})
}

// GetByBookID handles retrieving book transactions by book ID
func (h *BookTransactionHandler) GetByBookID(c *gin.Context) {
	idStr := c.Param("book_id")

	bookID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid book ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	transactions, err := h.bookTransactionService.GetByBookID(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve book transactions",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Book transactions retrieved successfully",
		Data:    transactions,
	})
}

// GetByStockCode handles retrieving book transactions by stock code
func (h *BookTransactionHandler) GetByStockCode(c *gin.Context) {
	stockCode := c.Param("stock_code")

	transactions, err := h.bookTransactionService.GetByStockCode(stockCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve book transactions",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Book transactions retrieved successfully",
		Data:    transactions,
	})
}

// CreateBookTransaction handles creating a new book transaction
func (h *BookTransactionHandler) CreateBookTransaction(c *gin.Context) {
	var req dto.BookTransactionCreateRequest
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

	transaction, err := h.bookTransactionService.CreateBookTransaction(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create book transaction",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseData{
		Status:  http.StatusCreated,
		Message: "Book transaction created successfully",
		Data:    transaction,
	})
}

// UpdateBookTransaction handles updating a book transaction
func (h *BookTransactionHandler) UpdateBookTransaction(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid transaction ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	var req dto.BookTransactionUpdateRequest
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

	transaction, err := h.bookTransactionService.UpdateBookTransaction(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update book transaction",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Book transaction updated successfully",
		Data:    transaction,
	})
}

// DeleteBookTransaction handles deleting a book transaction
func (h *BookTransactionHandler) DeleteBookTransaction(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid transaction ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	err = h.bookTransactionService.DeleteBookTransaction(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete book transaction",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Status:  http.StatusOK,
		Message: "Book transaction deleted successfully",
	})
}

// UpdateStatus handles updating a book transaction's status
func (h *BookTransactionHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid transaction ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	var req dto.BookTransactionStatusUpdateRequest
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

	transaction, err := h.bookTransactionService.UpdateStatus(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update book transaction status",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Book transaction status updated successfully",
		Data:    transaction,
	})
}

// ReturnBook handles returning a book
func (h *BookTransactionHandler) ReturnBook(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid transaction ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	var req dto.BookTransactionReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	transaction, err := h.bookTransactionService.ReturnBook(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to return book",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Book returned successfully",
		Data:    transaction,
	})
}

// GetOverdueTransactions handles retrieving all overdue transactions
func (h *BookTransactionHandler) GetOverdueTransactions(c *gin.Context) {
	transactions, err := h.bookTransactionService.GetOverdueTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve overdue transactions",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Overdue transactions retrieved successfully",
		Data:    transactions,
	})
}
