package handler

import (
	"go-gin-simple-api/dto"
	"go-gin-simple-api/lib"
	"go-gin-simple-api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MediaHandler struct {
	mediaService service.MediaService
}

func NewMediaHandler(mediaService service.MediaService) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
	}
}

// RegisterRoutes registers the media routes
// func (h *MediaHandler) RegisterRoutes(router *gin.RouterGroup) {
// 	media := router.Group("/media")
// 	{
// 		media.GET("", h.GetMedias)
// 		media.POST("", h.UploadMedia)
// 		media.GET("/:id", h.GetMedia)
// 		media.DELETE("/:id", h.DeleteMedia)
// 	}
// }

// ListMedia handles retrieving a list of media
func (h *MediaHandler) ListMedia(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.DefaultQuery("search", "")
	filterStr := c.DefaultQuery("filter", "")

	// Parse filters
	filter := lib.ParseFilterString(filterStr)

	result, err := h.mediaService.ListMedia(page, perPage, search, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch media",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UploadMedia handles uploading a new media file
func (h *MediaHandler) UploadMedia(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "No file uploaded",
		})
		return
	}

	media, err := h.mediaService.UploadMedia(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to upload media",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseData{
		Status:  http.StatusCreated,
		Message: "Media uploaded successfully",
		Data:    media,
	})
}

// GetMedia handles retrieving a specific media by ID
func (h *MediaHandler) GetMedia(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid media ID",
		})
		return
	}

	media, err := h.mediaService.GetMedia(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ResponseError{
			Status:  http.StatusNotFound,
			Message: "Media not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseData{
		Status:  http.StatusOK,
		Message: "Media retrieved successfully",
		Data:    media,
	})
}

// DeleteMedia handles deleting a media file
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid media ID",
		})
		return
	}

	if err := h.mediaService.DeleteMedia(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete media",
			Error:   map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Status:  http.StatusOK,
		Message: "Media deleted successfully",
	})
}
