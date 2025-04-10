package repository

import (
	"go-gin-simple-api/lib"
	"go-gin-simple-api/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaRepository interface {
	FindAll(page, perPage int, search string, filter lib.FilterParams) ([]model.Media, int64, error)
	FindByID(id uuid.UUID) (*model.Media, error)
	Create(media *model.Media) error
	Update(media *model.Media) error
	Delete(id uuid.UUID) error
	IsMediaUsed(id uuid.UUID) (bool, error)
}

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &mediaRepository{
		db: db,
	}
}

func (r *mediaRepository) FindAll(page, perPage int, search string, filter lib.FilterParams) ([]model.Media, int64, error) {
	var media []model.Media
	var total int64

	offset := (page - 1) * perPage
	query := r.db.Model(&model.Media{}).Preload("Books")

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
	if err := query.Limit(perPage).Offset(offset).Find(&media).Error; err != nil {
		return nil, 0, err
	}

	return media, total, nil
}

func (r *mediaRepository) FindByID(id uuid.UUID) (*model.Media, error) {
	var media model.Media
	if err := r.db.Preload("Books").First(&media, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *mediaRepository) Create(media *model.Media) error {
	return r.db.Create(media).Error
}

func (r *mediaRepository) Update(media *model.Media) error {
	return r.db.Save(media).Error
}

func (r *mediaRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Media{}, "id = ?", id).Error
}

// IsMediaUsed checks if a media is referenced by any book
// func (r *mediaRepository) IsMediaUsed(id uuid.UUID) ([]model.Media, bool, error) {
func (r *mediaRepository) IsMediaUsed(id uuid.UUID) (bool, error) {
	// var media []model.Media
	var count int64
	if err := r.db.Model(&model.Book{}).Where("cover_id = ?", id).Count(&count).Error; err != nil {
		// return nil, false, err
		return false, err
	}
	// if err := r.db.Model(&model.Book{}).Where("cover_id = ?", id).Find(&media).Error; err != nil {
	// 	return nil, false, err
	// }
	// return media, count > 0, nil
	return count > 0, nil
}
