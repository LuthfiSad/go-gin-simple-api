package repository

import (
	"go-gin-simple-api/lib"
	"go-gin-simple-api/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookRepository interface {
	FindBooks(page, perPage int, search string, filter lib.FilterParams) ([]model.Book, int64, error)
	FindByID(id uuid.UUID) (*model.Book, error)
	Create(book *model.Book) error
	Update(book *model.Book) error
	Delete(id uuid.UUID) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *bookRepository {
	return &bookRepository{
		db: db,
	}
}

func (r *bookRepository) FindBooks(page, perPage int, search string, filter lib.FilterParams) ([]model.Book, int64, error) {
	var books []model.Book
	var total int64

	offset := (page - 1) * perPage
	query := r.db.Model(&model.Book{}).Preload("Cover")

	// Apply search if provided
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply filters
	if len(filter) > 0 {
		for _, f := range filter {
			switch f.Operator {
			case lib.IsEqual:
				query = query.Where(f.Field+" = ?", f.Value)
			case lib.IsNotEqual:
				query = query.Where(f.Field+" != ?", f.Value)
			case lib.IsGreaterThan:
				query = query.Where(f.Field+" > ?", f.Value)
			case lib.IsGreaterEqual:
				query = query.Where(f.Field+" >= ?", f.Value)
			case lib.IsLessThan:
				query = query.Where(f.Field+" < ?", f.Value)
			case lib.IsLessEqual:
				query = query.Where(f.Field+" <= ?", f.Value)
			case lib.IsContain:
				query = query.Where(f.Field+" ILIKE ?", "%"+f.Value.(string)+"%")
			case lib.IsBeginWith:
				query = query.Where(f.Field+" ILIKE ?", f.Value.(string)+"%")
			case lib.IsEndWith:
				query = query.Where(f.Field+" ILIKE ?", "%"+f.Value.(string))
			case lib.IsIn:
				query = query.Where(f.Field+" IN ?", f.Value)
			}
		}
	}

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and get results
	if err := query.Limit(perPage).Offset(offset).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (r *bookRepository) FindByID(id uuid.UUID) (*model.Book, error) {
	var book model.Book
	if err := r.db.Preload("Cover").First(&book, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) Create(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) Update(book *model.Book) error {
	return r.db.Save(book).Error
}

func (r *bookRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Book{}, "id = ?", id).Error
}
