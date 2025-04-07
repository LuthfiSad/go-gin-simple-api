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

type BookHandler struct {
	bookService service.BookService
}

func NewBookHandler(bookService service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

// RegisterRoutes registers the book routes
// func (h *BookHandler) RegisterRoutes(router *gin.RouterGroup) {
// 	books := router.Group("/books")
// 	{
// 		books.GET("", h.GetBooks)
// 		books.GET("/:id", h.GetBook)
// 		books.POST("", h.CreateBook)
// 		books.PUT("/:id", h.UpdateBook)
// 		books.DELETE("/:id", h.DeleteBook)
// 		books.DELETE("/:id/cover", h.DeleteBookCover)
// 	}
// }

// GetBooks handles retrieving all books with pagination, search, and filtering
func (h *BookHandler) GetBooks(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")
	filterStr := c.Query("filter")

	// Parse filters
	filters := lib.ParseFilterString(filterStr)

	// Get books
	result, err := h.bookService.GetBooks(page, perPage, search, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve books",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetBook handles retrieving a specific book by ID
func (h *BookHandler) GetBookByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid book ID",
		})
		return
	}

	book, err := h.bookService.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ResponseError{
			Status:  http.StatusNotFound,
			Message: "Book not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Book retrieved successfully",
		Data:    book,
	})
}

// CreateBook handles creating a new book
func (h *BookHandler) CreateBook(c *gin.Context) {
	var req dto.BookCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
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

	_, err := h.bookService.CreateBook(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create book",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseSuccess{
		Status:  http.StatusCreated,
		Message: "Book created successfully",
		// Data:    book,
	})
}

// UpdateBook handles updating an existing book
func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, errParse := uuid.Parse(idStr)
	if errParse != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid book ID",
		})
		return
	}

	var req dto.BookUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
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

	_, err := h.bookService.UpdateBook(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update book",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Status:  http.StatusOK,
		Message: "Book updated successfully",
		// Data:    book,
	})
}

// DeleteBook handles deleting a book
func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid book ID",
		})
		return
	}

	if err := h.bookService.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete book",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Status:  http.StatusOK,
		Message: "Book deleted successfully",
	})
}

// DeleteBookCover handles removing a book's cover
func (h *BookHandler) DeleteBookCover(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid book ID",
		})
		return
	}

	if err := h.bookService.DeleteBookCover(id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete book cover",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Status:  http.StatusOK,
		Message: "Book cover deleted successfully",
	})
}
