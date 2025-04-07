package service

import (
	"errors"
	"go-gin-simple-api/dto"
	"go-gin-simple-api/lib"
	"go-gin-simple-api/model"
	"go-gin-simple-api/repository"

	"github.com/google/uuid"
)

type BookService interface {
	GetBooks(page, perPage int, search string, filter lib.FilterParams) (*dto.PaginatedResponseData[[]dto.BookRes], error)
	GetBookByID(id uuid.UUID) (*dto.BookRes, error)
	CreateBook(req dto.BookCreateReq) (*dto.BookRes, error)
	UpdateBook(id uuid.UUID, req dto.BookUpdateReq) (*dto.BookRes, error)
	DeleteBook(id uuid.UUID) error
	DeleteBookCover(id uuid.UUID) error
}

type bookService struct {
	repo repository.BookRepository
	// mediaRepo repository.MediaRepository
}

func NewBookService(repo repository.BookRepository) *bookService {
	return &bookService{
		repo: repo,
		// mediaRepo: mediaRepo,
	}
}

func (s *bookService) GetBooks(page, perPage int, search string, filter lib.FilterParams) (*dto.PaginatedResponseData[[]dto.BookRes], error) {
	books, total, err := s.repo.FindBooks(page, perPage, search, filter)
	if err != nil {
		return nil, err
	}

	bookResponses := make([]dto.BookRes, 0)
	for _, book := range books {
		bookResponses = append(bookResponses, mapBookToResponse(&book))
	}

	totalPages := (total + int64(perPage) - 1) / int64(perPage)
	if totalPages == 0 {
		totalPages = 1
	}

	return &dto.PaginatedResponseData[[]dto.BookRes]{
		Status:  200,
		Message: "Books retrieved successfully",
		Data:    bookResponses,
		Meta: dto.PaginationMeta{
			Page:        page,
			PerPage:     perPage,
			TotalPages:  totalPages,
			TotalItems:  total,
			ItemsOnPage: int64(len(bookResponses)),
		},
	}, nil
}

func (s *bookService) GetBookByID(id uuid.UUID) (*dto.BookRes, error) {
	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := mapBookToResponse(book)
	return &response, nil
}

func (s *bookService) CreateBook(req dto.BookCreateReq) (*dto.BookRes, error) {
	book := model.Book{
		Title:       req.Title,
		Description: req.Description,
	}

	if req.CoverID != nil {
		// cover, err := s.mediaRepo.GetMediaByID(*req.CoverID)
		// if err != nil {
		// 	return nil, err
		// }
		book.CoverID = req.CoverID
		// book.CoverID = &cover.ID
	}

	if err := s.repo.Create(&book); err != nil {
		return nil, err
	}

	response := mapBookToResponse(&book)
	return &response, nil
}

func (s *bookService) UpdateBook(id uuid.UUID, req dto.BookUpdateReq) (*dto.BookRes, error) {
	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Title != "" {
		book.Title = req.Title
	}

	if req.Description != "" {
		book.Description = req.Description
	}

	if req.CoverID != nil {
		// cover, err := s.mediaRepo.GetMediaByID(*req.CoverID)
		// if err != nil {
		// 	return nil, err
		// }
		book.CoverID = req.CoverID
		// book.CoverID = &cover.ID
	}

	if err := s.repo.Update(book); err != nil {
		return nil, err
	}

	response := mapBookToResponse(book)
	return &response, nil
}

func (s *bookService) DeleteBook(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *bookService) DeleteBookCover(id uuid.UUID) error {
	book, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if book.CoverID == nil {
		return errors.New("book has no cover")
	}

	// Reset cover ID
	book.CoverID = nil
	return s.repo.Update(book)
}

// Helper function to map domain book to DTO response
func mapBookToResponse(book *model.Book) dto.BookRes {
	response := dto.BookRes{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
	}

	if book.Cover != nil {
		response.Cover = &model.Media{
			ID:   book.Cover.ID,
			Path: book.Cover.Path,
		}
	}

	return response
}
