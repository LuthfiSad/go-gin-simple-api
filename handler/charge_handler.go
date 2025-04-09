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

type ChargeHandler struct {
	chargeService service.ChargeService
}

func NewChargeHandler(chargeService service.ChargeService) *ChargeHandler {
	return &ChargeHandler{
		chargeService: chargeService,
	}
}

// GetAll handles retrieving all charges with pagination, search, and filter
func (h *ChargeHandler) GetAll(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")
	filterStr := c.Query("filter")

	// Parse filters
	filters := lib.ParseFilterString(filterStr)

	// Get charges
	result, err := h.chargeService.GetAll(page, perPage, search, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve charges",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetByID handles retrieving a charge by ID
func (h *ChargeHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid charge ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	charge, err := h.chargeService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ResponseError{
			Status:  http.StatusNotFound,
			Message: "Charge not found",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Charge retrieved successfully",
		Data:    charge,
	})
}

// GetByBookTransactionID handles retrieving charges by book transaction ID
func (h *ChargeHandler) GetByBookTransactionID(c *gin.Context) {
	idStr := c.Param("transaction_id")

	transactionID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid transaction ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	charges, err := h.chargeService.GetByBookTransactionID(transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve charges",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Charges retrieved successfully",
		Data:    charges,
	})
}

// GetByUserID handles retrieving charges by user ID
func (h *ChargeHandler) GetByUserID(c *gin.Context) {
	idStr := c.Param("user_id")

	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	charges, err := h.chargeService.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve charges",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Charges retrieved successfully",
		Data:    charges,
	})
}

// Create handles creating a new charge
func (h *ChargeHandler) Create(c *gin.Context) {
	var req dto.ChargeCreateRequest
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

	userData, exists := c.Get("userData")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ResponseError{Status: http.StatusUnauthorized, Message: "User data not found"})
		return
	}

	user, ok := userData.(dto.UserData)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Status: http.StatusInternalServerError, Message: "Invalid user data"})
		return
	}

	charge, err := h.chargeService.Create(user.Email, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create charge",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}
	// Create handler for Charge (continuation)
	c.JSON(http.StatusCreated, dto.ResponseData{
		Status:  http.StatusCreated,
		Message: "Charge created successfully",
		Data:    charge,
	})
}

// Update handles updating a charge
func (h *ChargeHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid charge ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	var req dto.ChargeUpdateRequest
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

	charge, err := h.chargeService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update charge",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Charge updated successfully",
		Data:    charge,
	})
}

// Delete handles deleting a charge
func (h *ChargeHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid charge ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	err = h.chargeService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete charge",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Status:  http.StatusOK,
		Message: "Charge deleted successfully",
	})
}
