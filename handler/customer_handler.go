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

type CustomerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler(customerService service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

// GetAll handles retrieving all customers with pagination, search, and filter
func (h *CustomerHandler) GetAll(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")
	filterStr := c.Query("filter")

	// Parse filters
	filters := lib.ParseFilterString(filterStr)

	// Get customers
	result, err := h.customerService.GetAll(page, perPage, search, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve customers",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetByID handles retrieving a customer by ID
func (h *CustomerHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid customer ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	customer, err := h.customerService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ResponseError{
			Status:  http.StatusNotFound,
			Message: "Customer not found",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Customer retrieved successfully",
		Data:    customer,
	})
}

// GetByIDWithTransactions handles retrieving a customer with their book transactions
func (h *CustomerHandler) GetByIDWithTransactions(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid customer ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	customer, err := h.customerService.GetByIDWithTransactions(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ResponseError{
			Status:  http.StatusNotFound,
			Message: "Customer not found",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Customer with transactions retrieved successfully",
		Data:    customer,
	})
}

// GetByCode handles retrieving a customer by code
func (h *CustomerHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	customer, err := h.customerService.GetByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ResponseError{
			Status:  http.StatusNotFound,
			Message: "Customer not found",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Customer retrieved successfully",
		Data:    customer,
	})
}

// Create handles creating a new customer
func (h *CustomerHandler) Create(c *gin.Context) {
	var req dto.CustomerCreateRequest
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

	customer, err := h.customerService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create customer",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseData{
		Status:  http.StatusCreated,
		Message: "Customer created successfully",
		Data:    customer,
	})
}

// Update handles updating a customer
func (h *CustomerHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid customer ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	var req dto.CustomerUpdateRequest
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

	customer, err := h.customerService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update customer",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Customer updated successfully",
		Data:    customer,
	})
}

// Delete handles deleting a customer
func (h *CustomerHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid customer ID format",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	err = h.customerService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete customer",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Status:  http.StatusOK,
		Message: "Customer deleted successfully",
	})
}
