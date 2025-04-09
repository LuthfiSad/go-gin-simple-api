package service

import (
	"errors"
	"fmt"
	"go-gin-simple-api/dto"
	"go-gin-simple-api/lib"
	"go-gin-simple-api/model"
	"go-gin-simple-api/repository"

	"github.com/google/uuid"
)

type BookStockService interface {
	GetAll(page, perPage int, search string, filter lib.FilterParams) (*dto.PaginatedResponseData[[]dto.BookStockResponse], error)
	GetByCode(code string) (*dto.BookStockResponse, error)
	GetByBookID(bookID uuid.UUID) ([]dto.BookStockResponse, error)
	GetAvailableByBookID(bookID uuid.UUID) ([]dto.BookStockResponse, error)
	Create(bookStockRequest dto.BookStockCreateRequest) (*dto.BookStockResponse, error)
	Update(code string, bookStockRequest dto.BookStockUpdateRequest) (*dto.BookStockResponse, error)
	Delete(code string) error
	UpdateStatus(code string, req dto.BookStockStatusUpdateRequest) (*dto.BookStockResponse, error)
}

type bookStockService struct {
	repository repository.BookStockRepository
	bookRepo   repository.BookRepository
}

func NewBookStockService(repository repository.BookStockRepository, bookRepo repository.BookRepository) BookStockService {
	return &bookStockService{
		repository: repository,
		bookRepo:   bookRepo,
	}
}

func (s *bookStockService) GetAll(page, perPage int, search string, filter lib.FilterParams) (*dto.PaginatedResponseData[[]dto.BookStockResponse], error) {
	bookStocks, total, err := s.repository.FindAll(page, perPage, search, filter)
	if err != nil {
		return nil, err
	}

	bookStockResponses := make([]dto.BookStockResponse, 0)
	for _, bookStock := range bookStocks {
		bookStockResponses = append(bookStockResponses, mapToBookStockResponse(&bookStock))
	}

	// Calculate total pages
	totalPages := int64(total) / int64(perPage)
	if int64(total)%int64(perPage) > 0 {
		totalPages++
	}

	return &dto.PaginatedResponseData[[]dto.BookStockResponse]{
		Status:  200,
		Message: "Book stocks retrieved successfully",
		Data:    bookStockResponses,
		Meta: dto.PaginationMeta{
			Page:        page,
			PerPage:     perPage,
			TotalItems:  total,
			TotalPages:  totalPages,
			ItemsOnPage: int64(len(bookStockResponses)),
		},
	}, nil
}

func (s *bookStockService) GetByCode(code string) (*dto.BookStockResponse, error) {
	bookStock, err := s.repository.FindByCode(code)
	if err != nil {
		return nil, err
	}

	response := mapToBookStockResponse(bookStock)
	return &response, nil
}

func (s *bookStockService) GetByBookID(bookID uuid.UUID) ([]dto.BookStockResponse, error) {
	bookStocks, err := s.repository.FindByBookID(bookID)
	if err != nil {
		return nil, err
	}

	var bookStockResponses []dto.BookStockResponse
	for _, bookStock := range bookStocks {
		bookStockResponses = append(bookStockResponses, mapToBookStockResponse(&bookStock))
	}

	return bookStockResponses, nil
}

func (s *bookStockService) GetAvailableByBookID(bookID uuid.UUID) ([]dto.BookStockResponse, error) {
	bookStocks, err := s.repository.FindAvailableByBookID(bookID)
	if err != nil {
		return nil, err
	}

	var bookStockResponses []dto.BookStockResponse
	for _, bookStock := range bookStocks {
		bookStockResponses = append(bookStockResponses, mapToBookStockResponse(&bookStock))
	}

	return bookStockResponses, nil
}

func (s *bookStockService) Create(req dto.BookStockCreateRequest) (*dto.BookStockResponse, error) {
	// Check if book exists
	_, err := s.bookRepo.FindByID(req.BookID)
	if err != nil {
		return nil, errors.New("book not found")
	}

	// Check if code already exists
	existingStock, err := s.repository.FindByCode(req.Code)
	if err == nil && existingStock != nil {
		return nil, errors.New("code already exists")
	}

	bookStock := model.BookStock{
		Code:   req.Code,
		BookID: req.BookID,
		Status: model.StatusAvailable,
	}

	if req.Status != "" {
		bookStock.Status = req.Status
	}

	fmt.Println(bookStock)

	if err := s.repository.Create(&bookStock); err != nil {
		return nil, err
	}

	response := mapToBookStockResponse(&bookStock)
	return &response, nil
}

func (s *bookStockService) Update(code string, req dto.BookStockUpdateRequest) (*dto.BookStockResponse, error) {
	// Check if book stock exists
	bookStock, err := s.repository.FindByCode(code)
	if err != nil {
		return nil, errors.New("book stock not found")
	}
	fmt.Println("bookStock", bookStock)
	fmt.Printf("Book: %+v\n", bookStock.Book)

	// Check if book exists when book_id is provided
	if req.BookID != uuid.Nil {
		book, err := s.bookRepo.FindByID(req.BookID)
		if err != nil {
			return nil, errors.New("book not found")
		}
		fmt.Println("book", book)
		bookStock.BookID = book.ID
		bookStock.Book = *book
	}

	// Update fields if provided
	if req.Status != "" {
		bookStock.Status = req.Status
	}

	if err := s.repository.Update(bookStock); err != nil {
		return nil, err
	}

	response := mapToBookStockResponse(bookStock)
	return &response, nil
}

func (s *bookStockService) Delete(code string) error {
	// Check if book stock exists
	_, err := s.repository.FindByCode(code)
	if err != nil {
		return errors.New("book stock not found")
	}

	// Delete book stock
	return s.repository.Delete(code)
}

func (s *bookStockService) UpdateStatus(code string, req dto.BookStockStatusUpdateRequest) (*dto.BookStockResponse, error) {
	// Check if book stock exists
	bookStock, err := s.repository.FindByCode(code)
	if err != nil {
		return nil, errors.New("book stock not found")
	}

	// Update status
	if err := s.repository.UpdateStatus(code, req.Status); err != nil {
		return nil, err
	}

	response := mapToBookStockResponse(bookStock)
	return &response, nil
}

// Helper function to map a BookStock entity to a BookStockResponse DTO
func mapToBookStockResponse(bookStock *model.BookStock) dto.BookStockResponse {
	response := dto.BookStockResponse{
		Code:   bookStock.Code,
		BookID: bookStock.BookID,
		Status: bookStock.Status,
		// BorrowedID: bookStock.BorrowedID,
	}

	// if bookStock.BorrowedAt != nil {
	// 	response.BorrowedAt = bookStock.BorrowedAt
	// }

	if bookStock.Book.ID != uuid.Nil {
		response.Book = &dto.BookRes{
			ID:          bookStock.Book.ID,
			Title:       bookStock.Book.Title,
			Description: bookStock.Book.Description,
			CoverURL:    bookStock.Book.Cover.Path,
		}

		// if bookStock.Book.CoverID != nil {
		// 	response.Book.Cover.ID = *bookStock.Book.CoverID
		// 	if bookStock.Book.Cover != nil {
		// 		response.Book.Cover = bookStock.Book.Cover
		// 	}
		// }
	}

	// if bookStock.BorrowedID != nil {
	// 	response.User = &dto.UserData{
	// 		ID:    bookStock.User.ID,
	// 		Name:  bookStock.User.Name,
	// 		Email: bookStock.User.Email,
	// 		Role:  bookStock.User.Role,
	// 	}
	// }

	return response
}
