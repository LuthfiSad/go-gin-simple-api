package repository

import (
	"fmt"
	"go-gin-simple-api/lib"
	"go-gin-simple-api/model"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookStockRepository interface {
	FindAll(page, perPage int, search string, filter lib.FilterParams) ([]model.BookStock, int64, error)
	FindByCode(code string) (*model.BookStock, error)
	FindByBookID(bookID uuid.UUID) ([]model.BookStock, error)
	FindAvailableByBookID(bookID uuid.UUID) ([]model.BookStock, error)
	Create(bookStock *model.BookStock) error
	Update(bookStock *model.BookStock) error
	Delete(code string) error
	UpdateStatus(code, status string) error
}

type bookStockRepository struct {
	db *gorm.DB
}

func NewBookStockRepository(db *gorm.DB) BookStockRepository {
	return &bookStockRepository{db}
}

func (r *bookStockRepository) FindAll(page, perPage int, search string, filter lib.FilterParams) ([]model.BookStock, int64, error) {
	var bookStocks []model.BookStock
	var total int64

	query := r.db.Model(&model.BookStock{}).Preload("Book.Cover")

	// Join with Book to enable searching by book title
	query = query.Joins("LEFT JOIN books ON book_stocks.book_id = books.id")

	// Apply search if provided
	if search != "" {
		query = query.Where("book_stocks.code LIKE ? OR book_stocks.status LIKE ? OR books.title LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Apply filters
	if len(filter) > 0 {
		for _, f := range filter {
			switch f.Operator {
			case lib.IsEqual:
				query = query.Where(fmt.Sprintf("book_stocks.%s = ?", f.Field), f.Value)
			case lib.IsNotEqual:
				query = query.Where(fmt.Sprintf("book_stocks.%s != ?", f.Field), f.Value)
			case lib.IsGreaterThan:
				query = query.Where(fmt.Sprintf("book_stocks.%s > ?", f.Field), f.Value)
			case lib.IsGreaterEqual:
				query = query.Where(fmt.Sprintf("book_stocks.%s >= ?", f.Field), f.Value)
			case lib.IsLessThan:
				query = query.Where(fmt.Sprintf("book_stocks.%s < ?", f.Field), f.Value)
			case lib.IsLessEqual:
				query = query.Where(fmt.Sprintf("book_stocks.%s <= ?", f.Field), f.Value)
			case lib.IsContain:
				query = query.Where(fmt.Sprintf("book_stocks.%s LIKE ?", f.Field), "%"+fmt.Sprintf("%v", f.Value)+"%")
			case lib.IsBeginWith:
				query = query.Where(fmt.Sprintf("book_stocks.%s LIKE ?", f.Field), fmt.Sprintf("%v", f.Value)+"%")
			case lib.IsEndWith:
				query = query.Where(fmt.Sprintf("book_stocks.%s LIKE ?", f.Field), "%"+fmt.Sprintf("%v", f.Value))
			case lib.IsIn:
				if values, ok := f.Value.([]interface{}); ok {
					query = query.Where(fmt.Sprintf("book_stocks.%s IN ?", f.Field), values)
				} else if str, ok := f.Value.(string); ok {
					values := strings.Split(str, ",")
					query = query.Where(fmt.Sprintf("book_stocks.%s IN ?", f.Field), values)
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
	query = query.Preload("Book")

	// Execute query
	if err := query.Find(&bookStocks).Error; err != nil {
		return nil, 0, err
	}

	return bookStocks, total, nil
}

func (r *bookStockRepository) FindByCode(code string) (*model.BookStock, error) {
	var bookStock model.BookStock
	if err := r.db.Preload("Book").Preload("Book.Cover").First(&bookStock, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return &bookStock, nil
}

func (r *bookStockRepository) FindByBookID(bookID uuid.UUID) ([]model.BookStock, error) {
	var bookStocks []model.BookStock
	if err := r.db.Preload("Book").Preload("Book.Cover").Where("book_id = ?", bookID).Find(&bookStocks).Error; err != nil {
		return nil, err
	}
	return bookStocks, nil
}

func (r *bookStockRepository) FindAvailableByBookID(bookID uuid.UUID) ([]model.BookStock, error) {
	var bookStocks []model.BookStock
	if err := r.db.Preload("Book").Preload("Book.Cover").Where("book_id = ? AND status = ?", bookID, model.StatusAvailable).Find(&bookStocks).Error; err != nil {
		return nil, err
	}
	return bookStocks, nil
}

func (r *bookStockRepository) Create(bookStock *model.BookStock) error {
	return r.db.Create(bookStock).Error
}

func (r *bookStockRepository) Update(bookStock *model.BookStock) error {
	return r.db.Save(bookStock).Error
}

func (r *bookStockRepository) Delete(code string) error {
	return r.db.Delete(&model.BookStock{}, "code = ?", code).Error
}

func (r *bookStockRepository) UpdateStatus(code, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	return r.db.Model(&model.BookStock{}).Where("code = ?", code).Updates(updates).Error
}
