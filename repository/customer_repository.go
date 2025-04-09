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

type CustomerRepository interface {
	FindAll(page, perPage int, search string, filter lib.FilterParams) ([]model.Customer, int64, error)
	FindByID(id uuid.UUID) (*model.Customer, error)
	FindByCode(code string) (*model.Customer, error)
	Create(customer *model.Customer) error
	Update(customer *model.Customer) error
	Delete(id uuid.UUID) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db}
}

func (r *customerRepository) FindAll(page, perPage int, search string, filter lib.FilterParams) ([]model.Customer, int64, error) {
	var customers []model.Customer
	var total int64

	query := r.db.Model(&model.Customer{})

	// Apply search if provided
	if search != "" {
		query = query.Where("code LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply filters
	if len(filter) > 0 {
		for _, f := range filter {
			switch f.Operator {
			case lib.IsEqual:
				query = query.Where(fmt.Sprintf("%s = ?", f.Field), f.Value)
			case lib.IsNotEqual:
				query = query.Where(fmt.Sprintf("%s != ?", f.Field), f.Value)
			case lib.IsGreaterThan:
				query = query.Where(fmt.Sprintf("%s > ?", f.Field), f.Value)
			case lib.IsGreaterEqual:
				query = query.Where(fmt.Sprintf("%s >= ?", f.Field), f.Value)
			case lib.IsLessThan:
				query = query.Where(fmt.Sprintf("%s < ?", f.Field), f.Value)
			case lib.IsLessEqual:
				query = query.Where(fmt.Sprintf("%s <= ?", f.Field), f.Value)
			case lib.IsContain:
				query = query.Where(fmt.Sprintf("%s LIKE ?", f.Field), "%"+fmt.Sprintf("%v", f.Value)+"%")
			case lib.IsBeginWith:
				query = query.Where(fmt.Sprintf("%s LIKE ?", f.Field), fmt.Sprintf("%v", f.Value)+"%")
			case lib.IsEndWith:
				query = query.Where(fmt.Sprintf("%s LIKE ?", f.Field), "%"+fmt.Sprintf("%v", f.Value))
			case lib.IsIn:
				if values, ok := f.Value.([]interface{}); ok {
					query = query.Where(fmt.Sprintf("%s IN ?", f.Field), values)
				} else if str, ok := f.Value.(string); ok {
					values := strings.Split(str, ",")
					query = query.Where(fmt.Sprintf("%s IN ?", f.Field), values)
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

	// Execute query
	if err := query.Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *customerRepository) FindByID(id uuid.UUID) (*model.Customer, error) {
	var customer model.Customer
	if err := r.db.First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) FindByCode(code string) (*model.Customer, error) {
	var customer model.Customer
	if err := r.db.First(&customer, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Create(customer *model.Customer) error {
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	return r.db.Create(customer).Error
}

func (r *customerRepository) Update(customer *model.Customer) error {
	customer.UpdatedAt = time.Now()
	return r.db.Save(customer).Error
}

func (r *customerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Customer{}, "id = ?", id).Error
}
