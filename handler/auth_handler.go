package handler

import (
	"go-gin-simple-api/dto"
	"go-gin-simple-api/service"
	"go-gin-simple-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.AuthReq
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

	// Authenticate user
	res, err := h.authService.Authenticate(c, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ResponseError{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Login successful",
		Data:    res,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
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

	// Register user
	userData, err := h.authService.Register(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseData{
		Status:  http.StatusCreated,
		Message: "Registration successful",
		Data:    userData,
	})
}
