package utils

import (
	"errors"
	"fmt"
	"go-gin-simple-api/config"
	"go-gin-simple-api/dto"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(userData dto.UserData, cfg *config.Config) (string, error) {
	expiryHours, err := strconv.Atoi(cfg.JWTExpiryHours)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": userData.ID.String(),
		"email":   userData.Email,
		"name":    userData.Name,
		"role":    userData.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, cfg *config.Config) (dto.UserData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return dto.UserData{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			return dto.UserData{}, err
		}

		userData := dto.UserData{
			ID:    userID,
			Email: claims["email"].(string),
			Name:  claims["name"].(string),
			Role:  claims["role"].(string),
		}

		return userData, nil
	}

	return dto.UserData{}, errors.New("invalid token")
}
