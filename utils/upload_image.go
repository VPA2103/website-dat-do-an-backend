package utils

import (
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

// function trong một gói utils hoặc service
func UploadAndSaveImage(c *gin.Context, fieldName string, folder string, ownerID uint, ownerType string) error {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return err // Trả về lỗi nếu không có file hoặc lỗi đọc file
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Upload lên Cloudinary
	uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return err
	}

	// Lưu vào Database
	img := models.HinhAnh{
		OwnerID:   ownerID,
		OwnerType: ownerType,
		Url:  uploadResult.SecureURL,
	}

	return config.DB.Create(&img).Error
}
