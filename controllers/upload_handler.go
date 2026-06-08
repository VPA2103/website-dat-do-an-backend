package controllers

import (
	"context"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/dto"
)

func UploadHandler(c *gin.Context) {
	// Lấy file từ request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	//// Mở file
	//src, err := file.Open()
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
	//	return
	//}
	//defer src.Close()

	// Upload lên Cloudinary
	uploadResult, err := config.CLD.Upload.Upload(context.Background(), file, uploader.UploadParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Upload success",
		"url":     uploadResult.SecureURL,
	})
}

func GetImage(c *gin.Context) {
	var images []dto.HinhAnh
	if err := config.DB.Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product images"})
		return
	}
	c.JSON(http.StatusOK, images)
}
