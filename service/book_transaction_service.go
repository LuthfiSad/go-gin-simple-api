package service

import (
	"errors"
	"go-gin-simple-api/dto"
	"go-gin-simple-api/lib"
	"go-gin-simple-api/model"
	"go-gin-simple-api/repository"
	"time"

	"github.com/google/uuid"
)

type BookTransactionService interface {
	GetAll(page, perPage int, search string, filter lib.FilterParams) (*dto.PaginatedResponseData[[]dto.BookTransactionResponse], error)
	GetByID(id uuid.UUID) (*dto.BookTransactionResponse, error)
	GetByCustomerID(customerID uuid.UUID) ([]dto.BookTransactionResponse, error)
	GetByBookID(bookID uuid.UUID) ([]dto.BookTransactionResponse, error)
	GetByStockCode(stockCode string) ([]dto.BookTransactionResponse, error)
	Create(req dto.BookTransactionCreateRequest) (*dto.BookTransactionResponse, error)
	Update(id uuid.UUID, req dto.BookTransactionUpdateRequest) (*dto.BookTransactionResponse, error)
	Delete(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, req dto.BookTransactionStatusUpdateRequest) (*dto.BookTransactionResponse, error)
	ReturnBook(id uuid.UUID, req dto.BookTransactionReturnRequest) (*dto.BookTransactionResponse, error)
	GetOverdueTransactions() ([]dto.BookTransactionResponse, error)
}

type bookTransactionService struct {
	repository    repository.BookTransactionRepository
	bookRepo      repository.BookRepository
	bookStockRepo repository.BookStockRepository
	customerRepo  repository.CustomerRepository
}

func NewBookTransactionService(
	repository repository.BookTransactionRepository,
	bookRepo repository.BookRepository,
	bookStockRepo repository.BookStockRepository,
	customerRepo repository.CustomerRepository,
) BookTransactionService {
	return &bookTransactionService{
		repository:    repository,
		bookRepo:      bookRepo,
		bookStockRepo: bookStockRepo,
		customerRepo:  customerRepo,
	}
}

func (s *bookTransactionService) GetAll(page, perPage int, search string, filter lib.FilterParams) (*dto.PaginatedResponseData[[]dto.BookTransactionResponse], error) {
	transactions, total, err := s.repository.FindAll(page, perPage, search, filter)
	if err != nil {
		return nil, err
	}

	transactionResponses := make([]dto.BookTransactionResponse, 0)
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, mapToBookTransactionResponse(&transaction))
	}

	// Calculate total pages
	totalPages := int64(total) / int64(perPage)
	if int64(total)%int64(perPage) > 0 {
		totalPages++
	}

	return &dto.PaginatedResponseData[[]dto.BookTransactionResponse]{
		Status:  200,
		Message: "Book transactions retrieved successfully",
		Data:    transactionResponses,
		Meta: dto.PaginationMeta{
			Page:        page,
			PerPage:     perPage,
			TotalItems:  total,
			TotalPages:  totalPages,
			ItemsOnPage: int64(len(transactionResponses)),
		},
	}, nil
}

func (s *bookTransactionService) GetByID(id uuid.UUID) (*dto.BookTransactionResponse, error) {
	transaction, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := mapToBookTransactionResponse(transaction)
	return &response, nil
}

func (s *bookTransactionService) GetByCustomerID(customerID uuid.UUID) ([]dto.BookTransactionResponse, error) {
	transactions, err := s.repository.FindByCustomerID(customerID)
	if err != nil {
		return nil, err
	}

	var responses []dto.BookTransactionResponse
	for _, transaction := range transactions {
		responses = append(responses, mapToBookTransactionResponse(&transaction))
	}

	return responses, nil
}

func (s *bookTransactionService) GetByBookID(bookID uuid.UUID) ([]dto.BookTransactionResponse, error) {
	transactions, err := s.repository.FindByBookID(bookID)
	if err != nil {
		return nil, err
	}

	var responses []dto.BookTransactionResponse
	for _, transaction := range transactions {
		responses = append(responses, mapToBookTransactionResponse(&transaction))
	}

	return responses, nil
}

func (s *bookTransactionService) GetByStockCode(stockCode string) ([]dto.BookTransactionResponse, error) {
	transactions, err := s.repository.FindByStockCode(stockCode)
	if err != nil {
		return nil, err
	}

	var responses []dto.BookTransactionResponse
	for _, transaction := range transactions {
		responses = append(responses, mapToBookTransactionResponse(&transaction))
	}

	return responses, nil
}

func (s *bookTransactionService) Create(req dto.BookTransactionCreateRequest) (*dto.BookTransactionResponse, error) {
	// Check if book exists
	// book, err := s.bookRepo.FindByID(req.BookID)
	// if err != nil {
	// 	return nil, errors.New("book not found")
	// }

	// Check if book stock exists
	bookStock, err := s.bookStockRepo.FindByCode(req.StockCode)
	if err != nil {
		return nil, errors.New("book stock not found")
	}

	// Check if book stock is available
	if bookStock.Status != model.StatusAvailable {
		return nil, errors.New("book stock is not available")
	}

	// Check if customer exists
	customer, err := s.customerRepo.FindByID(req.CustomerID)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	now := time.Now()
	transaction := model.BookTransaction{
		ID:         uuid.New(),
		BookID:     bookStock.BookID,
		StockCode:  req.StockCode,
		CustomerID: customer.ID,
		DueDate:    time.Now().AddDate(0, 0, 7),
		Status:     req.Status,
		BorrowedAt: &now,
	}

	if err := s.repository.Create(&transaction); err != nil {
		return nil, err
	}

	// Update book stock status to borrowed
	if err := s.bookStockRepo.UpdateStatus(req.StockCode, model.StatusBorrowed); err != nil {
		return nil, errors.New("failed to update book stock status")
	}

	response := mapToBookTransactionResponse(&transaction)
	return &response, nil
}

func (s *bookTransactionService) Update(id uuid.UUID, req dto.BookTransactionUpdateRequest) (*dto.BookTransactionResponse, error) {
	// Check if transaction exists
	transaction, err := s.repository.FindByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	// // Check if book exists when book_id is provided
	// if req.BookID != uuid.Nil {
	// 	_, err := s.bookRepo.FindByID(req.BookID)
	// 	if err != nil {
	// 		return nil, errors.New("book not found")
	// 	}
	// 	transaction.BookID = req.BookID
	// }

	// Check if book stock exists when stock_code is provided
	if req.StockCode != "" {
		stockCode, err := s.bookStockRepo.FindByCode(req.StockCode)
		if err != nil {
			return nil, errors.New("book stock not found")
		}

		// If changing stock code, update old stock status back to Available
		if transaction.StockCode != req.StockCode {
			if err := s.bookStockRepo.UpdateStatus(transaction.StockCode, model.StatusAvailable); err != nil {
				return nil, errors.New("failed to update original book stock status")
			}

			// Set new stock to Borrowed
			if err := s.bookStockRepo.UpdateStatus(req.StockCode, model.StatusBorrowed); err != nil {
				return nil, errors.New("failed to update new book stock status")
			}
		}

		transaction.BookID = stockCode.BookID
		transaction.StockCode = req.StockCode
	}

	// Check if customer exists when customer_id is provided
	if req.CustomerID != uuid.Nil {
		_, err := s.customerRepo.FindByID(req.CustomerID)
		if err != nil {
			return nil, errors.New("customer not found")
		}
		transaction.CustomerID = req.CustomerID
	}

	// Update fields if provided
	// transaction.DueDate = time.Now().AddDate(0, 0, 7)
	if req.DueDate != nil {
		transaction.DueDate = *req.DueDate
	}

	if req.Status != "" {
		transaction.Status = req.Status

		// If status changed to Returned, update book stock status to Available
		if req.Status == model.StatusBTReturned && (transaction.ReturnAt == nil || req.ReturnAt != nil) {
			if err := s.bookStockRepo.UpdateStatus(transaction.StockCode, model.StatusAvailable); err != nil {
				return nil, errors.New("failed to update book stock status")
			}

			// Set return date if not provided
			if req.ReturnAt == nil {
				now := time.Now()
				transaction.ReturnAt = &now
			} else {
				transaction.ReturnAt = req.ReturnAt
			}
		}
	}

	if req.ReturnAt != nil {
		transaction.ReturnAt = req.ReturnAt
	}

	if err := s.repository.Update(transaction); err != nil {
		return nil, err
	}

	response := mapToBookTransactionResponse(transaction)
	return &response, nil
}

func (s *bookTransactionService) Delete(id uuid.UUID) error {
	// Check if transaction exists
	transaction, err := s.repository.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	// If the transaction is active (Borrowed or Overdue), restore book stock status to Available
	if transaction.Status == model.StatusBTBorrowed || transaction.Status == model.StatusBTOverdue {
		if err := s.bookStockRepo.UpdateStatus(transaction.StockCode, model.StatusAvailable); err != nil {
			return errors.New("failed to update book stock status")
		}
	}

	// Delete transaction
	return s.repository.Delete(id)
}

func (s *bookTransactionService) UpdateStatus(id uuid.UUID, req dto.BookTransactionStatusUpdateRequest) (*dto.BookTransactionResponse, error) {
	// Check if transaction exists
	transaction, err := s.repository.FindByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	// Update status
	if err := s.repository.UpdateStatus(id, req.Status); err != nil {
		return nil, err
	}

	// Update book stock status based on transaction status change
	if req.Status == model.StatusBTBorrowed {
		if err := s.bookStockRepo.UpdateStatus(transaction.StockCode, model.StatusBorrowed); err != nil {
			return nil, errors.New("failed to update book stock status")
		}
	} else if req.Status == model.StatusBTReturned {
		if err := s.bookStockRepo.UpdateStatus(transaction.StockCode, model.StatusAvailable); err != nil {
			return nil, errors.New("failed to update book stock status")
		}

		// Set return date
		now := time.Now()
		transaction.ReturnAt = &now
		errUpdate := s.repository.Update(transaction)
		if errUpdate != nil {
			return nil, errUpdate
		}
	}

	transaction.Status = req.Status
	response := mapToBookTransactionResponse(transaction)
	return &response, nil
}

func (s *bookTransactionService) ReturnBook(id uuid.UUID, req dto.BookTransactionReturnRequest) (*dto.BookTransactionResponse, error) {
	// Check if transaction exists
	transaction, err := s.repository.FindByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	// Check if transaction is already returned
	if transaction.Status == model.StatusBTReturned {
		return nil, errors.New("book is already returned")
	}

	returnAt := time.Now()
	if req.ReturnAt != nil {
		returnAt = *req.ReturnAt
	}

	// Return book
	if err := s.repository.ReturnBook(id, returnAt); err != nil {
		return nil, err
	}

	// Update book stock status to Available
	if err := s.bookStockRepo.UpdateStatus(transaction.StockCode, model.StatusAvailable); err != nil {
		return nil, errors.New("failed to update book stock status")
	}

	// Get updated transaction
	updatedTransaction, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := mapToBookTransactionResponse(updatedTransaction)
	return &response, nil
}

func (s *bookTransactionService) GetOverdueTransactions() ([]dto.BookTransactionResponse, error) {
	transactions, err := s.repository.FindOverdueTransactions()
	if err != nil {
		return nil, err
	}

	// Update status to Overdue for all overdue transactions
	var responses []dto.BookTransactionResponse
	for _, transaction := range transactions {
		if transaction.Status != model.StatusBTOverdue {
			// Update status to Overdue
			err := s.repository.UpdateStatus(transaction.ID, model.StatusBTOverdue)
			if err != nil {
				return nil, err
			}
			transaction.Status = model.StatusBTOverdue
		}
		responses = append(responses, mapToBookTransactionResponse(&transaction))
	}

	return responses, nil
}

// Helper function to map a BookTransaction entity to a BookTransactionResponse DTO
func mapToBookTransactionResponse(transaction *model.BookTransaction) dto.BookTransactionResponse {
	response := dto.BookTransactionResponse{
		ID:         transaction.ID,
		BookID:     transaction.BookID,
		StockCode:  transaction.StockCode,
		CustomerID: transaction.CustomerID,
		DueDate:    transaction.DueDate,
		Status:     transaction.Status,
		BorrowedAt: transaction.BorrowedAt,
		ReturnAt:   transaction.ReturnAt,
	}

	if transaction.Book.ID != uuid.Nil {
		response.Book = &dto.BookRes{
			ID:          transaction.Book.ID,
			Title:       transaction.Book.Title,
			Description: transaction.Book.Description,
		}

		if transaction.Book.Cover.Path != "" {
			response.Book.CoverURL = transaction.Book.Cover.Path
		}
	}

	if transaction.BookStock.Code != "" {
		response.BookStock = &dto.BookStockResponse{
			Code:   transaction.BookStock.Code,
			BookID: transaction.BookStock.BookID,
			Status: transaction.BookStock.Status,
		}
	}

	if transaction.Customer.ID != uuid.Nil {
		response.Customer = &dto.CustomerResponse{
			ID:   transaction.Customer.ID,
			Code: transaction.Customer.Code,
			Name: transaction.Customer.Name,
		}
	}

	if len(transaction.Charges) > 0 {
		response.Charges = make([]dto.ChargeResponse, 0)
		for _, charge := range transaction.Charges {
			response.Charges = append(response.Charges, dto.ChargeResponse{
				ID:                charge.ID,
				BookTransactionID: charge.BookTransactionID,
				DaysLate:          charge.DaysLate,
				DailyLateFee:      charge.DailyLateFee,
				Total:             charge.Total,
				UserID:            charge.UserID,
				CreatedAt:         charge.CreatedAt,
			})
		}
	}

	return response
}
