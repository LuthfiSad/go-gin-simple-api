package service

import (
	"context"
	"errors"
	"go-gin-simple-api/config"
	"go-gin-simple-api/dto"
	"go-gin-simple-api/model"
	"go-gin-simple-api/repository"
	"go-gin-simple-api/utils"
)

type AuthService interface {
	Authenticate(ctx context.Context, req dto.AuthReq) (dto.AuthRes, error)
	Register(ctx context.Context, req dto.RegisterReq) (dto.UserData, error)
	Validate(ctx context.Context, tokenString string) (dto.UserData, error)
}

type authService struct {
	repo repository.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo repository.AuthRepository) *authService {
	cfg, _ := config.LoadConfig()
	return &authService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *authService) Authenticate(ctx context.Context, req dto.AuthReq) (dto.AuthRes, error) {
	// Validate request
	if validationErrors := utils.Validate(req); len(validationErrors) > 0 {
		return dto.AuthRes{}, errors.New("validation failed")
	}

	// Find user by email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return dto.AuthRes{}, errors.New("invalid credentials")
	}

	// Verify password
	if err := utils.VerifyPassword(user.Password, req.Password); err != nil {
		return dto.AuthRes{}, errors.New("invalid credentials")
	}

	// Generate user data
	userData := dto.UserData{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	// Generate token
	token, err := utils.GenerateToken(userData, s.cfg)
	if err != nil {
		return dto.AuthRes{}, err
	}

	return dto.AuthRes{
		Token: token,
		// User:  userData,
	}, nil
}

func (s *authService) Register(ctx context.Context, req dto.RegisterReq) (dto.UserData, error) {
	// Validate request
	if validationErrors := utils.Validate(req); len(validationErrors) > 0 {
		return dto.UserData{}, errors.New("validation failed")
	}

	// Check if user already exists
	existingUser, _ := s.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return dto.UserData{}, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto.UserData{}, err
	}

	// Create user
	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user", // Default role
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	if err := s.repo.Create(&user); err != nil {
		return dto.UserData{}, err
	}

	// Generate user data
	userData := dto.UserData{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	return userData, nil
}

func (s *authService) Validate(ctx context.Context, tokenString string) (dto.UserData, error) {
	return utils.ValidateToken(tokenString, s.cfg)
}
