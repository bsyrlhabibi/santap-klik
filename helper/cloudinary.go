package helper

import (
	"context"
	"io"
	"santapKlik/configs"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(folder string, file io.Reader) (string, error) {
	cfg := configs.InitConfig()
	cld, err := cloudinary.NewFromURL(cfg.CLOUDINARY) // Use the CLOUDINARY_URL from the configuration

	if err != nil {
		return "", err
	}

	ctx := context.Background()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folder,
	})

	if err != nil {
		return "", err
	}

	return uploadResult.PublicID, nil
}
