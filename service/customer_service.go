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

type CustomerService interface {
	GetAll(page, perPage int, search string, filter lib.FilterParams) (*dto.PaginatedResponseData[[]dto.CustomerResponse], error)
	GetByID(id uuid.UUID) (*dto.CustomerResponse, error)
	GetByIDWithTransactions(id uuid.UUID) (*dto.CustomerWithTransactionsResponse, error)
	GetByCode(code string) (*dto.CustomerResponse, error)
	Create(req dto.CustomerCreateRequest) (*dto.CustomerResponse, error)
	Update(id uuid.UUID, req dto.CustomerUpdateRequest) (*dto.CustomerResponse, error)
	Delete(id uuid.UUID) error
}

type customerService struct {
	repository          repository.CustomerRepository
	bookTransactionRepo repository.BookTransactionRepository
}

func NewCustomerService(
	repository repository.CustomerRepository,
	bookTransactionRepo repository.BookTransactionRepository,
) CustomerService {
	return &customerService{
		repository:          repository,
		bookTransactionRepo: bookTransactionRepo,
	}
}

func (s *customerService) GetAll(page, perPage int, search string, filter lib.FilterParams) (*dto.PaginatedResponseData[[]dto.CustomerResponse], error) {
	customers, total, err := s.repository.FindAll(page, perPage, search, filter)
	if err != nil {
		return nil, err
	}

	customerResponses := make([]dto.CustomerResponse, 0)
	for _, customer := range customers {
		customerResponses = append(customerResponses, mapToCustomerResponse(&customer))
	}

	// Calculate total pages
	totalPages := int64(total) / int64(perPage)
	if int64(total)%int64(perPage) > 0 {
		totalPages++
	}

	return &dto.PaginatedResponseData[[]dto.CustomerResponse]{
		Status:  200,
		Message: "Customers retrieved successfully",
		Data:    customerResponses,
		Meta: dto.PaginationMeta{
			Page:        page,
			PerPage:     perPage,
			TotalItems:  total,
			TotalPages:  totalPages,
			ItemsOnPage: int64(len(customerResponses)),
		},
	}, nil
}

func (s *customerService) GetByID(id uuid.UUID) (*dto.CustomerResponse, error) {
	customer, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := mapToCustomerResponse(customer)
	return &response, nil
}

func (s *customerService) GetByIDWithTransactions(id uuid.UUID) (*dto.CustomerWithTransactionsResponse, error) {
	customer, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Get customer's book transactions
	transactions, err := s.bookTransactionRepo.FindByCustomerID(customer.ID)
	if err != nil {
		return nil, err
	}

	response := mapToCustomerWithTransactionsResponse(customer, transactions)
	return &response, nil
}

func (s *customerService) GetByCode(code string) (*dto.CustomerResponse, error) {
	customer, err := s.repository.FindByCode(code)
	if err != nil {
		return nil, err
	}

	response := mapToCustomerResponse(customer)
	return &response, nil
}

func (s *customerService) Create(req dto.CustomerCreateRequest) (*dto.CustomerResponse, error) {
	// Check if code already exists
	existingCustomer, err := s.repository.FindByCode(req.Code)
	if err == nil && existingCustomer != nil {
		return nil, errors.New("customer code already exists")
	}

	customer := model.Customer{
		ID:        uuid.New(),
		Code:      req.Code,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repository.Create(&customer); err != nil {
		return nil, err
	}

	response := mapToCustomerResponse(&customer)
	return &response, nil
}

func (s *customerService) Update(id uuid.UUID, req dto.CustomerUpdateRequest) (*dto.CustomerResponse, error) {
	// Check if customer exists
	customer, err := s.repository.FindByID(id)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	// Update fields if provided
	if req.Code != nil {
		// Check if the new code already exists (excluding current customer)
		existingCustomer, err := s.repository.FindByCode(*req.Code)
		if err == nil && existingCustomer != nil && existingCustomer.ID != customer.ID {
			return nil, errors.New("customer code already exists")
		}
		customer.Code = *req.Code
	}

	if req.Name != nil {
		customer.Name = *req.Name
	}

	if err := s.repository.Update(customer); err != nil {
		return nil, err
	}

	response := mapToCustomerResponse(customer)
	return &response, nil
}

func (s *customerService) Delete(id uuid.UUID) error {
	// Check if customer exists
	customer, err := s.repository.FindByID(id)
	if err != nil {
		return errors.New("customer not found")
	}

	// Check if customer has any transactions before deleting
	transactions, err := s.bookTransactionRepo.FindByCustomerID(customer.ID)
	if err == nil && len(transactions) > 0 {
		return errors.New("cannot delete customer with existing book transactions")
	}

	// Delete customer
	return s.repository.Delete(id)
}

// Helper function to map a Customer entity to a CustomerResponse DTO
func mapToCustomerResponse(customer *model.Customer) dto.CustomerResponse {
	return dto.CustomerResponse{
		ID:        customer.ID,
		Code:      customer.Code,
		Name:      customer.Name,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}

// Helper function to map a Customer entity with transactions to a CustomerWithTransactionsResponse DTO
func mapToCustomerWithTransactionsResponse(customer *model.Customer, transactions []model.BookTransaction) dto.CustomerWithTransactionsResponse {
	response := dto.CustomerWithTransactionsResponse{
		ID:        customer.ID,
		Code:      customer.Code,
		Name:      customer.Name,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}

	// Map book transactions
	if len(transactions) > 0 {
		bookTransactions := make([]dto.BookTransactionResponse, 0)
		for _, transaction := range transactions {
			bookTransactions = append(bookTransactions, mapToBookTransactionResponse(&transaction))
		}
		response.BookTransactions = bookTransactions
	}

	return response
}
