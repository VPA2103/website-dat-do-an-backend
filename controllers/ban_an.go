package controllers

import (
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

func CreateBanAn(c *gin.Context) {
	var ban models.BanAn

	// ✅ Bind form data
	if err := c.ShouldBind(&ban); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu form không hợp lệ: " + err.Error()})
		return
	}

	// ✅ Mặc định trạng thái là "Trống"
	//if ban.TrangThai != 0 {
	//	defaultTrangThai := 0
	//	ban.TrangThai = defaultTrangThai
	//}

	// ✅ Tạo record trong DB trước để có MaBan
	if err := config.DB.Create(&ban).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo bàn ăn: " + err.Error()})
		return
	}

	// ✅ Tạo URL menu
	// menuURL := fmt.Sprintf(
	// 	"http://localhost:4200/#/customer/goimon/menu?table=%d",
	// 	ban.MaBan,
	// )

	// ✅ Tạo QR
	// qrBytes, err := utils.GenerateQRBytes(menuURL)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo mã QR: " + err.Error()})
	// 	return
	// }

	// ✅ Upload QR trực tiếp lên Cloudinary
	// uploadResult, err := config.CLD.Upload.Upload(c, bytes.NewReader(qrBytes), uploader.UploadParams{
	// 	Folder:   "banan_qr",
	// 	PublicID: fmt.Sprintf("qr_ban_%d", ban.MaBan),
	// })

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Upload QR thất bại: " + err.Error(),
	// 	})
	// 	return
	// }

	// ban.Anh_QR = uploadResult.SecureURL
	config.DB.Save(&ban)

	// ✅ Upload ảnh bàn (nếu có)
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()

			uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
				Folder: "banan",
			})
			if err == nil {
				img := models.HinhAnh{
					OwnerID:   ban.MaBanAn,
					OwnerType: "ban_an",
					Url:  uploadResult.SecureURL,
				}
				config.DB.Create(&img)
			}
		}
	}

	config.DB.Preload("AnhBan").First(&ban, ban.MaBanAn)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo bàn ăn thành công",
		"data":    ban,
	})
}

// Lấy tất cả bàn ăn kèm ảnh
func GetAllBanAn(c *gin.Context) {
	var dsBanAn []models.BanAn

	// ✅ Preload ảnh bàn (quan hệ polymorphic)
	if err := config.DB.Preload("AnhBan").Find(&dsBanAn).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách bàn ăn: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lấy danh sách bàn ăn thành công",
		"data":    dsBanAn,
	})
}

func GetBanAnByID(c *gin.Context) {
	id := c.Param("id")

	var banan models.BanAn

	// 🔥 Query đúng: WHERE id = ? + Preload ảnh
	if err := config.DB.Preload("AnhBan").First(&banan, "ma_ban = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy bàn ăn với ID " + id,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lấy thông tin bàn ăn thành công",
		"data":    banan,
	})
}

// ✅ Cập nhật thông tin bàn ăn
func UpdateBanAn(c *gin.Context) {
	id := c.Param("id")
	var ban models.BanAn

	// 1️⃣ Tìm bàn ăn
	if err := config.DB.First(&ban, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bàn ăn"})
		return
	}

	// 2️⃣ Bind dữ liệu form
	var input models.BanAn
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu gửi lên không hợp lệ: " + err.Error(),
		})
		return
	}

	// 3️⃣ Update text (AN TOÀN)
	ban.TenBan = input.TenBan
	ban.SoChoNgoi = input.SoChoNgoi
	ban.TrangThai = input.TrangThai

	if err := config.DB.Save(&ban).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể cập nhật bàn ăn: " + err.Error(),
		})
		return
	}

	// 4️⃣ Upload ảnh mới (nếu có)
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không mở được file ảnh"})
			return
		}
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(
			c,
			src,
			uploader.UploadParams{
				Folder: "banan",
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload ảnh lỗi"})
			return
		}

		// 🔥 XÓA TOÀN BỘ ẢNH CŨ CỦA BÀN ĂN
		config.DB.
			Where("owner_id = ? AND owner_type = ?", ban.MaBanAn, "ban_an").
			Delete(&models.HinhAnh{})

		// 🔥 THÊM ẢNH MỚI
		config.DB.Create(&models.HinhAnh{
			OwnerID:   ban.MaBanAn,
			OwnerType: "ban_an",
			Url:  uploadResult.SecureURL,
		})
	}

	// 5️⃣ Load lại quan hệ ảnh
	config.DB.Preload("AnhBan").First(&ban, ban.MaBanAn)

	// 6️⃣ Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật bàn ăn thành công",
		"data":    ban,
	})
}

// ✅ Xóa bàn ăn
func DeleteBanAn(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var ban models.BanAn
	if err := config.DB.First(&ban, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bàn ăn"})
		return
	}

	// 🔹 Xóa ảnh liên quan (nếu có)
	config.DB.Where("owner_id = ? AND owner_type = ?", id, "ban_an").Delete(&models.HinhAnh{})

	// 🔹 Xóa bàn ăn
	if err := config.DB.Delete(&ban).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xóa bàn ăn: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa bàn ăn thành công",
	})
}
