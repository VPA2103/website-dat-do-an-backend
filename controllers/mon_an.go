package controllers

import (
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

// ======================= CREATE =======================
func CreateMonAn(c *gin.Context) {
	var monan models.MonAn

	// Bind dữ liệu form
	if err := c.ShouldBind(&monan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	// Validate
	if monan.TenMonAn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên món ăn không được để trống"})
		return
	}

	if monan.MoTa == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mô tả món ăn không được để trống"})
		return
	}

	// Tạo trước để lấy ID
	if err := config.DB.Create(&monan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo món ăn: " + err.Error()})
		return
	}

	// Upload ảnh món ăn nếu có
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()

			uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
				Folder: "monan",
			})

			if err == nil {
				img := models.HinhAnh{
					OwnerID:   monan.MaMonAn,
					OwnerType: "mon_an",
					Url:       uploadResult.SecureURL,
				}
				config.DB.Create(&img)
			}
		}
	}

	// Lấy món ăn kèm ảnh trả về client
	config.DB.Preload("AnhMonAn").First(&monan, monan.MaMonAn)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo món ăn thành công",
		"data":    monan,
	})
}

// ======================= GET ALL =======================
func GetAllMonAn(c *gin.Context) {
	var list []models.MonAn
	config.DB.Preload("AnhMonAn").Find(&list)

	c.JSON(http.StatusOK, gin.H{"data": list})
}

func GetMonAnByID(c *gin.Context) {
	id := c.Param("id")
	var monan models.MonAn

	if err := config.DB.Preload("AnhMonAn").First(&monan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy món ăn"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": monan})
}

// ======================= UPDATE =======================
func UpdateMonAn(c *gin.Context) {
	id := c.Param("id")
	var monan models.MonAn

	// 1. Tìm món ăn
	if err := config.DB.First(&monan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy món ăn"})
		return
	}

	// 2. Bind & update text (AN TOÀN)
	var input models.MonAn
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.MoTa == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Mô tả món ăn không được để trống",
		})
		return
	}

	config.DB.Model(&monan).Updates(input)

	// 3. Upload ảnh mới (nếu có)
	file, err := c.FormFile("image")
	if err == nil {
		src, _ := file.Open()
		defer src.Close()

		upload, err := config.CLD.Upload.Upload(
			c,
			src,
			uploader.UploadParams{
				Folder: "monan",
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload ảnh lỗi"})
			return
		}

		// 🔥 XÓA TẤT CẢ ẢNH CŨ
		config.DB.
			Where("owner_id = ? AND owner_type = ?", monan.MaMonAn, "mon_an").
			Delete(&models.HinhAnh{})

		// 🔥 THÊM ẢNH MỚI
		config.DB.Create(&models.HinhAnh{
			Url:       upload.SecureURL,
			OwnerID:   monan.MaMonAn,
			OwnerType: "mon_an",
		})
	}

	// 4. Load lại quan hệ ảnh
	config.DB.Preload("AnhMonAn").First(&monan, id)

	// 5. Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật món ăn thành công",
		"data":    monan,
	})
}

// ======================= DELETE =======================
func DeleteMonAn(c *gin.Context) {
	id := c.Param("id")
	var monan models.MonAn

	if err := config.DB.First(&monan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy món ăn"})
		return
	}

	// Xóa ảnh thuộc món ăn
	config.DB.Where("owner_id = ? AND owner_type = ?", id, "mon_an").Delete(&models.HinhAnh{})

	// Xóa món ăn
	config.DB.Delete(&monan)

	c.JSON(http.StatusOK, gin.H{"message": "Xóa món ăn thành công"})
}
