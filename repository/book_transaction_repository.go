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

type BookTransactionRepository interface {
	FindAll(page, perPage int, search string, filter lib.FilterParams) ([]model.BookTransaction, int64, error)
	FindByID(id uuid.UUID) (*model.BookTransaction, error)
	FindByCustomerID(customerID uuid.UUID) ([]model.BookTransaction, error)
	FindByBookID(bookID uuid.UUID) ([]model.BookTransaction, error)
	FindByStockCode(stockCode string) ([]model.BookTransaction, error)
	FindActiveByStockCode(stockCode string) (*model.BookTransaction, error)
	Create(transaction *model.BookTransaction) error
	Update(transaction *model.BookTransaction) error
	Delete(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, status string) error
	ReturnBook(id uuid.UUID, returnAt time.Time) error
	FindOverdueTransactions() ([]model.BookTransaction, error)
}

type bookTransactionRepository struct {
	db *gorm.DB
}

func NewBookTransactionRepository(db *gorm.DB) BookTransactionRepository {
	return &bookTransactionRepository{db}
}

func (r *bookTransactionRepository) FindAll(page, perPage int, search string, filter lib.FilterParams) ([]model.BookTransaction, int64, error) {
	var transactions []model.BookTransaction
	var total int64

	query := r.db.Model(&model.BookTransaction{})

	// Join with related tables to enable search
	query = query.Joins("LEFT JOIN books ON book_transactions.book_id = books.id")
	query = query.Joins("LEFT JOIN book_stocks ON book_transactions.stock_code = book_stocks.code")
	query = query.Joins("LEFT JOIN customers ON book_transactions.customer_id = customers.id")

	// Apply search if provided
	if search != "" {
		query = query.Where("book_transactions.status LIKE ? OR books.title LIKE ? OR book_stocks.code LIKE ? OR customers.name LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Apply filters
	if len(filter) > 0 {
		for _, f := range filter {
			switch f.Operator {
			case lib.IsEqual:
				query = query.Where(fmt.Sprintf("book_transactions.%s = ?", f.Field), f.Value)
			case lib.IsNotEqual:
				query = query.Where(fmt.Sprintf("book_transactions.%s != ?", f.Field), f.Value)
			case lib.IsGreaterThan:
				query = query.Where(fmt.Sprintf("book_transactions.%s > ?", f.Field), f.Value)
			case lib.IsGreaterEqual:
				query = query.Where(fmt.Sprintf("book_transactions.%s >= ?", f.Field), f.Value)
			case lib.IsLessThan:
				query = query.Where(fmt.Sprintf("book_transactions.%s < ?", f.Field), f.Value)
			case lib.IsLessEqual:
				query = query.Where(fmt.Sprintf("book_transactions.%s <= ?", f.Field), f.Value)
			case lib.IsContain:
				query = query.Where(fmt.Sprintf("book_transactions.%s LIKE ?", f.Field), "%"+fmt.Sprintf("%v", f.Value)+"%")
			case lib.IsBeginWith:
				query = query.Where(fmt.Sprintf("book_transactions.%s LIKE ?", f.Field), fmt.Sprintf("%v", f.Value)+"%")
			case lib.IsEndWith:
				query = query.Where(fmt.Sprintf("book_transactions.%s LIKE ?", f.Field), "%"+fmt.Sprintf("%v", f.Value))
			case lib.IsIn:
				if values, ok := f.Value.([]interface{}); ok {
					query = query.Where(fmt.Sprintf("book_transactions.%s IN ?", f.Field), values)
				} else if str, ok := f.Value.(string); ok {
					values := strings.Split(str, ",")
					query = query.Where(fmt.Sprintf("book_transactions.%s IN ?", f.Field), values)
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
	query = query.Preload("Book").Preload("Book.Cover").Preload("BookStock").Preload("Customer").Preload("Charges")

	// Execute query
	if err := query.Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *bookTransactionRepository) FindByID(id uuid.UUID) (*model.BookTransaction, error) {
	var transaction model.BookTransaction
	if err := r.db.Preload("Book").Preload("Book.Cover").Preload("BookStock").Preload("Customer").Preload("Charges").First(&transaction, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *bookTransactionRepository) FindByCustomerID(customerID uuid.UUID) ([]model.BookTransaction, error) {
	var transactions []model.BookTransaction
	if err := r.db.Preload("Book").Preload("Book.Cover").Preload("BookStock").Preload("Customer").Where("customer_id = ?", customerID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *bookTransactionRepository) FindByBookID(bookID uuid.UUID) ([]model.BookTransaction, error) {
	var transactions []model.BookTransaction
	if err := r.db.Preload("Book").Preload("Book.Cover").Preload("BookStock").Preload("Customer").Where("book_id = ?", bookID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *bookTransactionRepository) FindByStockCode(stockCode string) ([]model.BookTransaction, error) {
	var transactions []model.BookTransaction
	if err := r.db.Preload("Book").Preload("Book.Cover").Preload("BookStock").Preload("Customer").Where("stock_code = ?", stockCode).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *bookTransactionRepository) FindActiveByStockCode(stockCode string) (*model.BookTransaction, error) {
	var transaction model.BookTransaction
	if err := r.db.Preload("Book").Preload("Book.Cover").Preload("BookStock").Preload("Customer").Where("stock_code = ? AND status IN ('Borrowed', 'Overdue')", stockCode).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *bookTransactionRepository) Create(transaction *model.BookTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *bookTransactionRepository) Update(transaction *model.BookTransaction) error {
	return r.db.Save(transaction).Error
}

func (r *bookTransactionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.BookTransaction{}, "id = ?", id).Error
}

func (r *bookTransactionRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&model.BookTransaction{}).Where("id = ?", id).Update("status", status).Error
}

func (r *bookTransactionRepository) ReturnBook(id uuid.UUID, returnAt time.Time) error {
	updates := map[string]interface{}{
		"status":    model.StatusBTReturned,
		"return_at": returnAt,
	}
	return r.db.Model(&model.BookTransaction{}).Where("id = ?", id).Updates(updates).Error
}

func (r *bookTransactionRepository) FindOverdueTransactions() ([]model.BookTransaction, error) {
	var transactions []model.BookTransaction
	now := time.Now()

	if err := r.db.Preload("Book").Preload("Book.Cover").Preload("BookStock").Preload("Customer").
		Where("status = 'Borrowed' AND due_date < ?", now).
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}
