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

type ChargeService interface {
	GetAll(page, perPage int, search string, filter lib.FilterParams) (*dto.PaginatedResponseData[[]dto.ChargeResponse], error)
	GetByID(id uuid.UUID) (*dto.ChargeResponse, error)
	GetByBookTransactionID(bookTransactionID uuid.UUID) ([]dto.ChargeResponse, error)
	GetByUserID(userID uuid.UUID) ([]dto.ChargeResponse, error)
	Create(userEmail string, req dto.ChargeCreateRequest) (*dto.ChargeResponse, error)
	Update(id uuid.UUID, req dto.ChargeUpdateRequest) (*dto.ChargeResponse, error)
	Delete(id uuid.UUID) error
}

type chargeService struct {
	repository          repository.ChargeRepository
	bookTransactionRepo repository.BookTransactionRepository
	userRepo            repository.AuthRepository
}

func NewChargeService(
	repository repository.ChargeRepository,
	bookTransactionRepo repository.BookTransactionRepository,
	userRepo repository.AuthRepository,
) ChargeService {
	return &chargeService{
		repository:          repository,
		bookTransactionRepo: bookTransactionRepo,
		userRepo:            userRepo,
	}
}

func (s *chargeService) GetAll(page, perPage int, search string, filter lib.FilterParams) (*dto.PaginatedResponseData[[]dto.ChargeResponse], error) {
	charges, total, err := s.repository.FindAll(page, perPage, search, filter)
	if err != nil {
		return nil, err
	}

	chargeResponses := make([]dto.ChargeResponse, 0)
	for _, charge := range charges {
		chargeResponses = append(chargeResponses, mapToChargeResponse(&charge))
	}

	// Calculate total pages
	totalPages := int64(total) / int64(perPage)
	if int64(total)%int64(perPage) > 0 {
		totalPages++
	}

	return &dto.PaginatedResponseData[[]dto.ChargeResponse]{
		Status:  200,
		Message: "Charges retrieved successfully",
		Data:    chargeResponses,
		Meta: dto.PaginationMeta{
			Page:        page,
			PerPage:     perPage,
			TotalItems:  total,
			TotalPages:  totalPages,
			ItemsOnPage: int64(len(chargeResponses)),
		},
	}, nil
}

func (s *chargeService) GetByID(id uuid.UUID) (*dto.ChargeResponse, error) {
	charge, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := mapToChargeResponse(charge)
	return &response, nil
}

func (s *chargeService) GetByBookTransactionID(bookTransactionID uuid.UUID) ([]dto.ChargeResponse, error) {
	charges, err := s.repository.FindByBookTransactionID(bookTransactionID)
	if err != nil {
		return nil, err
	}

	var responses []dto.ChargeResponse
	for _, charge := range charges {
		responses = append(responses, mapToChargeResponse(&charge))
	}

	return responses, nil
}

func (s *chargeService) GetByUserID(userID uuid.UUID) ([]dto.ChargeResponse, error) {
	charges, err := s.repository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.ChargeResponse
	for _, charge := range charges {
		responses = append(responses, mapToChargeResponse(&charge))
	}

	return responses, nil
}

func (s *chargeService) Create(userEmail string, req dto.ChargeCreateRequest) (*dto.ChargeResponse, error) {
	// Check if book transaction exists
	transaction, err := s.bookTransactionRepo.FindByID(req.BookTransactionID)
	if err != nil {
		return nil, errors.New("book transaction not found")
	}

	// Check if user exists
	user, errExists := s.userRepo.FindByEmail(userEmail)
	if errExists != nil {
		return nil, errors.New("user not found")
	}

	// Calculate total charge
	total := float64(req.DaysLate) * req.DailyLateFee

	charge := model.Charge{
		ID:                uuid.New(),
		BookTransactionID: transaction.ID,
		DaysLate:          req.DaysLate,
		DailyLateFee:      req.DailyLateFee,
		Total:             total,
		UserID:            user.ID,
		CreatedAt:         time.Now(),
	}

	if err := s.repository.Create(&charge); err != nil {
		return nil, err
	}

	// Get the created charge with relationships loaded
	createdCharge, err := s.repository.FindByID(charge.ID)
	if err != nil {
		return nil, err
	}

	response := mapToChargeResponse(createdCharge)
	return &response, nil
}

func (s *chargeService) Update(id uuid.UUID, req dto.ChargeUpdateRequest) (*dto.ChargeResponse, error) {
	// Check if charge exists
	charge, err := s.repository.FindByID(id)
	if err != nil {
		return nil, errors.New("charge not found")
	}

	// Update fields if provided
	if req.DaysLate != nil {
		charge.DaysLate = *req.DaysLate
	}

	if req.DailyLateFee != nil {
		charge.DailyLateFee = *req.DailyLateFee
	}

	// Recalculate total
	charge.Total = float64(charge.DaysLate) * charge.DailyLateFee

	if err := s.repository.Update(charge); err != nil {
		return nil, err
	}

	response := mapToChargeResponse(charge)
	return &response, nil
}

func (s *chargeService) Delete(id uuid.UUID) error {
	// Check if charge exists
	_, err := s.repository.FindByID(id)
	if err != nil {
		return errors.New("charge not found")
	}

	// Delete charge
	return s.repository.Delete(id)
}

// Helper function to map a Charge entity to a ChargeResponse DTO
func mapToChargeResponse(charge *model.Charge) dto.ChargeResponse {
	response := dto.ChargeResponse{
		ID:                charge.ID,
		BookTransactionID: charge.BookTransactionID,
		DaysLate:          charge.DaysLate,
		DailyLateFee:      charge.DailyLateFee,
		Total:             charge.Total,
		UserID:            charge.UserID,
		CreatedAt:         charge.CreatedAt,
	}

	if charge.BookTransaction.ID != uuid.Nil {
		bookTransaction := mapToBookTransactionResponse(&charge.BookTransaction)
		response.BookTransaction = &bookTransaction
	}

	if charge.User.ID != uuid.Nil {
		response.User = &dto.UserData{
			ID:    charge.User.ID,
			Name:  charge.User.Name,
			Email: charge.User.Email,
			Role:  charge.User.Role,
		}
	}

	return response
}
