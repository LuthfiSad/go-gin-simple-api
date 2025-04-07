package repository

import (
	"errors"
	"go-gin-simple-api/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(email string) (*model.User, error)
	Create(user *model.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}
