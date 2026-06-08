package controllers

import (
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

func CreateLoaiMonAn(c *gin.Context) {
	var loaimonan models.LoaiMonAn

	// Lấy dữ liệu form-data
	if err := c.ShouldBind(&loaimonan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ: " + err.Error(),
		})
		return
	}

	// Kiểm tra tên loại món ăn
	if loaimonan.TenLoaiMonAn == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tên loại món ăn không được để trống",
		})
		return
	}

	// Upload ảnh nếu có
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Không thể mở file ảnh",
			})
			return
		}
		defer src.Close()

		// Upload Cloudinary
		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
			Folder: "loaimonan",
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Upload ảnh thất bại: " + err.Error(),
			})
			return
		}

		loaimonan.AnhLoaiMonAn = uploadResult.SecureURL
	}

	// Lưu database
	if err := config.DB.Create(&loaimonan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể tạo loại món ăn: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo loại món ăn thành công",
		"data":    loaimonan,
	})
}

func GetAllLoaiMonAn(c *gin.Context) {
	var list []models.LoaiMonAn

	if err := config.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể lấy danh sách loại món ăn: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

func GetLoaiMonAnByID(c *gin.Context) {
	id := c.Param("id")
	var loaimonan models.LoaiMonAn

	if err := config.DB.First(&loaimonan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy loại món ăn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": loaimonan,
	})
}

func UpdateLoaiMonAn(c *gin.Context) {
	id := c.Param("id")
	var loaimonan models.LoaiMonAn

	// Kiểm tra tồn tại
	if err := config.DB.First(&loaimonan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy loại món ăn",
		})
		return
	}

	// Bind dữ liệu mới
	if err := c.ShouldBind(&loaimonan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ: " + err.Error(),
		})
		return
	}

	// Upload ảnh mới nếu có
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Không thể mở file ảnh",
			})
			return
		}
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
			Folder: "loaimonan",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Upload ảnh thất bại: " + err.Error(),
			})
			return
		}

		loaimonan.AnhLoaiMonAn = uploadResult.SecureURL
	}

	// Cập nhật DB
	if err := config.DB.Save(&loaimonan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể cập nhật loại món ăn: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật thành công",
		"data":    loaimonan,
	})
}

func DeleteLoaiMonAn(c *gin.Context) {
	id := c.Param("id")
	var loaimonan models.LoaiMonAn

	// Kiểm tra tồn tại
	if err := config.DB.First(&loaimonan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy loại món ăn",
		})
		return
	}

	// Xóa DB
	if err := config.DB.Delete(&loaimonan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể xóa loại món ăn: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa loại món ăn thành công",
	})
}
