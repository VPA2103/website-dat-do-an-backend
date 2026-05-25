package controllers

import (
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

type UpdateNhaHangDTO struct {
	TenNhaHang string `form:"ten_nha_hang"`
	TrangThai  int    `form:"trang_thai"`
	DiaChi   string    `form:"dia_chi"`
	SoTaiKhoan   int   `form:"so_tai_khoan"`
	NganHang   int     `form:"ngan_hang"`

}

func CreateNhaHang(c *gin.Context) {
	var nhahang models.NhaHang

	// ======================
	// BIND FORM DATA
	// ======================
	if err := c.ShouldBind(&nhahang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ: " + err.Error(),
		})
		return
	}

	if nhahang.TenNhaHang == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tên nhà hàng không được để trống",
		})
		return
	}
	if nhahang.DiaChi == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Địa chỉ nhà hàng không được để trống",
		})
		return
	}
	if nhahang.SoTaiKhoan == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Số tài khoản không được để trống",
		})
		return
	}
	if nhahang.NganHang == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ngân hàng không được để trống",
		})
		return
	}

	// ======================
	// GET USER FROM TOKEN
	// ======================
	userAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Chưa đăng nhập",
		})
		return
	}
	nhahang.MaNguoiDung = userAny.(uint)

	// ======================
	// CREATE NHÀ HÀNG (LẤY ID)
	// ======================
	if err := config.DB.Create(&nhahang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể tạo nhà hàng: " + err.Error(),
		})
		return
	}

	// ======================
	// UPLOAD ẢNH (GIỐNG BanAn)
	// ======================
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()

			uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
				Folder: "nhahang",
			})
			if err == nil {
				img := models.HinhAnh{
					OwnerID:   nhahang.MaNhaHang,
					OwnerType: "nha_hang",
					Url:       uploadResult.SecureURL,
				}
				config.DB.Create(&img)
			}
		}
	}

	// ======================
	// PRELOAD ẢNH TRẢ VỀ
	// ======================
	config.DB.Preload("AnhNhaHang").First(&nhahang, nhahang.MaNhaHang)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo nhà hàng thành công",
		"data":    nhahang,
	})
}
func GetAllNhaHang(c *gin.Context) {
	var list []models.NhaHang

	// ✅ Preload ảnh nhà hàng (polymorphic)
	if err := config.DB.
		Preload("AnhNhaHang").
		Find(&list).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể lấy danh sách nhà hàng: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lấy danh sách nhà hàng thành công",
		"data":    list,
	})
}

func GetNhaHangByID(c *gin.Context) {
	id := c.Param("id")
	var nhahang models.NhaHang

	if err := config.DB.First(&nhahang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy nhà hàng",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": nhahang,
	})
}

func GetNhaHangByUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"message": "Không xác định được user",
		})
		return
	}

	var nhaHangs []models.NhaHang

	err := config.DB.
		Where("ma_nguoi_dung = ?", userID).
		Preload("AnhNhaHang").
		Find(&nhaHangs).Error

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Lỗi khi lấy nhà hàng",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": nhaHangs,
		"message": "Lấy nhà hàng theo user thành công",
	})
}

func UpdateNhaHang(c *gin.Context) {
	id := c.Param("id")
	var nhahang models.NhaHang

	if err := config.DB.First(&nhahang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy nhà hàng",
		})
		return
	}

	// ===== CHECK USER (CHỦ NHÀ HÀNG)
	userAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Chưa đăng nhập"})
		return
	}
	userID := userAny.(uint)

	if nhahang.MaNguoiDung != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Không có quyền sửa"})
		return
	}

	// ===== BIND DATA
	if err := c.ShouldBind(&nhahang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ===== UPLOAD ẢNH MỚI
	file, err := c.FormFile("image")
	if err == nil && file != nil {

		// ❌ XÓA ẢNH CŨ
		config.DB.
			Where("owner_id = ? AND owner_type = ?", nhahang.MaNhaHang, "nha_hang").
			Delete(&models.HinhAnh{})

		src, _ := file.Open()
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
			Folder: "nhahang",
		})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		hinhAnh := models.HinhAnh{
			Url:       uploadResult.SecureURL,
			OwnerID:   nhahang.MaNhaHang,
			OwnerType: "nha_hang",
		}
		config.DB.Create(&hinhAnh)
	}

	if err := config.DB.Save(&nhahang).Error; err != nil {
		c.JSON(500, gin.H{"error": "Không thể cập nhật"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Cập nhật thành công",
	})
}

func DeleteNhaHang(c *gin.Context) {
	id := c.Param("id")
	var nhahang models.NhaHang

	if err := config.DB.First(&nhahang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy nhà hàng",
		})
		return
	}

	if err := config.DB.Delete(&nhahang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể xóa nhà hàng",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa nhà hàng thành công",
	})
}