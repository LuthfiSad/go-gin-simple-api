package lib

import (
	"bytes"
	"context"
	"go-gin-simple-api/config"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	Cld *cloudinary.Cloudinary
}

func NewCloudinaryService(cfg *config.Config) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(
		cfg.CloudinaryName,
		cfg.CloudinaryKey,
		cfg.CloudinarySecret,
	)
	if err != nil {
		return nil, err
	}

	return &CloudinaryService{Cld: cld}, nil
}

func (c *CloudinaryService) UploadImage(file []byte, folder string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uploadParams := uploader.UploadParams{
		Folder: folder,
	}

	reader := bytes.NewReader(file)

	result, err := c.Cld.Upload.Upload(ctx, reader, uploadParams)
	if err != nil {
		return "", "", err
	}

	return result.SecureURL, result.PublicID, nil
}

func (c *CloudinaryService) DeleteImage(publicID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := c.Cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	return err
}
