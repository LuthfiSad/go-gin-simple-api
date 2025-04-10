package repository

import (
	"fmt"
	"go-gin-simple-api/lib"
	"go-gin-simple-api/model"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChargeRepository interface {
	FindAll(page, perPage int, search string, filter lib.FilterParams) ([]model.Charge, int64, error)
	FindByID(id uuid.UUID) (*model.Charge, error)
	FindByBookTransactionID(bookTransactionID uuid.UUID) ([]model.Charge, error)
	FindByUserID(userID uuid.UUID) ([]model.Charge, error)
	Create(charge *model.Charge) error
	Update(charge *model.Charge) error
	Delete(id uuid.UUID) error
}

type chargeRepository struct {
	db *gorm.DB
}

func NewChargeRepository(db *gorm.DB) ChargeRepository {
	return &chargeRepository{db}
}

func (r *chargeRepository) FindAll(page, perPage int, search string, filter lib.FilterParams) ([]model.Charge, int64, error) {
	var charges []model.Charge
	var total int64

	query := r.db.Model(&model.Charge{})

	// Join with related tables to enable search
	query = query.Joins("LEFT JOIN book_transactions ON charges.book_transaction_id = book_transactions.id")
	query = query.Joins("LEFT JOIN users ON charges.user_id = users.id")

	// Apply search if provided
	if search != "" {
		query = query.Where("users.name LIKE ? OR book_transactions.id::text LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// Apply filters
	if len(filter) > 0 {
		for _, f := range filter {
			switch f.Operator {
			case lib.IsEqual:
				query = query.Where(fmt.Sprintf("charges.%s = ?", f.Field), f.Value)
			case lib.IsNotEqual:
				query = query.Where(fmt.Sprintf("charges.%s != ?", f.Field), f.Value)
			case lib.IsGreaterThan:
				query = query.Where(fmt.Sprintf("charges.%s > ?", f.Field), f.Value)
			case lib.IsGreaterEqual:
				query = query.Where(fmt.Sprintf("charges.%s >= ?", f.Field), f.Value)
			case lib.IsLessThan:
				query = query.Where(fmt.Sprintf("charges.%s < ?", f.Field), f.Value)
			case lib.IsLessEqual:
				query = query.Where(fmt.Sprintf("charges.%s <= ?", f.Field), f.Value)
			case lib.IsContain:
				query = query.Where(fmt.Sprintf("charges.%s LIKE ?", f.Field), "%"+fmt.Sprintf("%v", f.Value)+"%")
			case lib.IsBeginWith:
				query = query.Where(fmt.Sprintf("charges.%s LIKE ?", f.Field), fmt.Sprintf("%v", f.Value)+"%")
			case lib.IsEndWith:
				query = query.Where(fmt.Sprintf("charges.%s LIKE ?", f.Field), "%"+fmt.Sprintf("%v", f.Value))
			case lib.IsIn:
				if values, ok := f.Value.([]interface{}); ok {
					query = query.Where(fmt.Sprintf("charges.%s IN ?", f.Field), values)
				} else if str, ok := f.Value.(string); ok {
					values := strings.Split(str, ",")
					query = query.Where(fmt.Sprintf("charges.%s IN ?", f.Field), values)
				}
			}
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * perPage
	if page > 0 && perPage > 0 {
		query = query.Offset(offset).Limit(perPage)
	}

	// Preload relationships
	query = query.Preload("BookTransaction.Book").Preload("BookTransaction.Customer").Preload("User")

	// Execute query
	if err := query.Find(&charges).Error; err != nil {
		return nil, 0, err
	}

	return charges, total, nil
}

func (r *chargeRepository) FindByID(id uuid.UUID) (*model.Charge, error) {
	var charge model.Charge
	if err := r.db.Preload("BookTransaction").Preload("BookTransaction.Book").Preload("BookTransaction.Book.Cover").Preload("BookTransaction.Customer").Preload("User").First(&charge, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &charge, nil
}

func (r *chargeRepository) FindByBookTransactionID(bookTransactionID uuid.UUID) ([]model.Charge, error) {
	var charges []model.Charge
	if err := r.db.Preload("BookTransaction.Book").Preload("BookTransaction.Book.Cover").Preload("BookTransaction.Customer").Preload("User").Where("book_transaction_id = ?", bookTransactionID).Find(&charges).Error; err != nil {
		return nil, err
	}
	return charges, nil
}

func (r *chargeRepository) FindByUserID(userID uuid.UUID) ([]model.Charge, error) {
	var charges []model.Charge
	if err := r.db.Preload("BookTransaction.Book").Preload("BookTransaction.Book.Cover").Preload("BookTransaction.Customer").Preload("User").Where("user_id = ?", userID).Find(&charges).Error; err != nil {
		return nil, err
	}
	return charges, nil
}

func (r *chargeRepository) Create(charge *model.Charge) error {
	charge.CreatedAt = time.Now()
	return r.db.Create(charge).Error
}

func (r *chargeRepository) Update(charge *model.Charge) error {
	return r.db.Save(charge).Error
}

func (r *chargeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Charge{}, "id = ?", id).Error
}
