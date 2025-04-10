package service

import (
	"context"
	"errors"
	"go-gin-simple-api/dto"
	"go-gin-simple-api/lib"
	"go-gin-simple-api/model"
	"go-gin-simple-api/repository"
	"go-gin-simple-api/utils"
	"io/ioutil"
	"mime/multipart"

	"github.com/google/uuid"
)

type MediaService interface {
	ListMedia(page, perPage int, search string, filter lib.FilterParams) (*dto.PaginatedResponseData[[]dto.MediaRes], error)
	UploadMedia(ctx context.Context, file *multipart.FileHeader) (*dto.MediaRes, error)
	GetMedia(id uuid.UUID) (*dto.MediaRes, error)
	DeleteMedia(ctx context.Context, id uuid.UUID) error
}

type mediaService struct {
	repo       repository.MediaRepository
	repoBook   repository.BookRepository
	cloudinary *lib.CloudinaryService
}

func NewMediaService(repo repository.MediaRepository, repoBook repository.BookRepository, cloudinary *lib.CloudinaryService) MediaService {
	return &mediaService{
		repo:       repo,
		repoBook:   repoBook,
		cloudinary: cloudinary,
	}
}

func (s *mediaService) ListMedia(page, perPage int, search string, filter lib.FilterParams) (*dto.PaginatedResponseData[[]dto.MediaRes], error) {
	media, total, err := s.repo.FindAll(page, perPage, search, filter)
	if err != nil {
		return nil, err
	}

	mediaRes := make([]dto.MediaRes, 0)
	for _, m := range media {
		mediaRes = append(mediaRes, mapMediaToResponse(&m))
	}

	totalPages := (total + int64(perPage) - 1) / int64(perPage)
	if totalPages == 0 {
		totalPages = 1
	}

	return &dto.PaginatedResponseData[[]dto.MediaRes]{
		Status:  200,
		Message: "Medias retrieved successfully",
		Data:    mediaRes,
		Meta: dto.PaginationMeta{
			Page:        page,
			PerPage:     perPage,
			TotalPages:  totalPages,
			TotalItems:  total,
			ItemsOnPage: int64(len(mediaRes)),
		},
	}, nil
}

// UploadMedia uploads a media file and stores its reference
func (s *mediaService) UploadMedia(ctx context.Context, file *multipart.FileHeader) (*dto.MediaRes, error) {
	// Check file type
	if !utils.IsValidImageType(file.Header.Get("Content-Type")) {
		return nil, errors.New("invalid file type. only images are allowed")
	}
	// Open uploaded file
	openedFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer openedFile.Close()

	// Read file content
	fileBytes, err := ioutil.ReadAll(openedFile)
	if err != nil {
		return nil, err
	}

	// Upload to Cloudinary
	path, publicID, err := s.cloudinary.UploadImage(fileBytes, "media")
	if err != nil {
		return nil, err
	}

	// Create media record
	media := model.Media{
		Path:     path,
		PublicID: publicID,
	}

	if err := s.repo.Create(&media); err != nil {
		return nil, err
	}

	response := mapMediaToResponse(&media)
	return &response, nil
}

// GetMedia retrieves a media by ID
func (s *mediaService) GetMedia(id uuid.UUID) (*dto.MediaRes, error) {
	media, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	res := mapMediaToResponse(media)
	return &res, nil
}

// DeleteMedia deletes a media file
func (s *mediaService) DeleteMedia(ctx context.Context, id uuid.UUID) error {
	media, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// remove media from books
	// for _, book := range media.Books {
	// 	book.CoverID = nil
	// 	book.Cover = nil
	// 	if err := s.repoBook.Update(&book); err != nil {
	// 		return err
	// 	}
	// }

	// Check if the media is used by any book
	used, err := s.repo.IsMediaUsed(id)
	if err != nil {
		return err
	}
	if used {
		return errors.New("media is currently in use and cannot be deleted")
	}

	// Delete from Cloudinary
	if err := s.cloudinary.DeleteImage(media.PublicID); err != nil {
		return err
	}

	// Delete from database
	return s.repo.Delete(id)
}

func mapMediaToResponse(media *model.Media) dto.MediaRes {
	response := dto.MediaRes{
		ID:        media.ID,
		Path:      media.Path,
		CreatedAt: media.CreatedAt,
		UpdatedAt: media.UpdatedAt,
	}

	if media.Books != nil {
		response.Books = media.Books
	}

	return response
}
