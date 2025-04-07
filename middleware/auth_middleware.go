package middleware

import (
	"go-gin-simple-api/config"
	"go-gin-simple-api/dto"
	"go-gin-simple-api/repository"
	"go-gin-simple-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth(r repository.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ResponseError{Status: http.StatusUnauthorized, Message: "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.ResponseError{Status: http.StatusUnauthorized, Message: "Invalid token format"})
			c.Abort()
			return
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ResponseError{Status: http.StatusInternalServerError, Message: "Failed to load config"})
			c.Abort()
			return
		}

		userData, err := utils.ValidateToken(parts[1], cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ResponseError{Status: http.StatusUnauthorized, Message: "Invalid or expired token"})
			c.Abort()
			return
		}

		if user, errExist := r.FindByEmail(userData.Email); errExist != nil || user.ID != userData.ID {
			c.JSON(http.StatusNotFound, dto.ResponseError{
				Status:  http.StatusNotFound,
				Message: "User not found",
			})
			c.Abort()
			return
		}

		// Set user data in context for use in handlers
		c.Set("userData", userData)
		c.Next()
	}
}

func RoleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userData, exists := c.Get("userData")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.ResponseError{Status: http.StatusUnauthorized, Message: "User data not found"})
			c.Abort()
			return
		}

		user, ok := userData.(dto.UserData)
		if !ok {
			c.JSON(http.StatusInternalServerError, dto.ResponseError{Status: http.StatusInternalServerError, Message: "Invalid user data"})
			c.Abort()
			return
		}

		// Check if user role is in the allowed roles
		allowed := false
		for _, role := range roles {
			if user.Role == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, dto.ResponseError{Status: http.StatusForbidden, Message: "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
