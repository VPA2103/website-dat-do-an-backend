package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/dto"
)

type UpdateNhaHangDTO struct {
	TenNhaHang string `form:"ten_nha_hang"`
	TrangThai  int    `form:"trang_thai"`
	DiaChi   string    `form:"dia_chi"`
	SoTaiKhoan   int   `form:"so_tai_khoan"`
	NganHang   int     `form:"ngan_hang"`

}

func (h *ChatHandler) CreateNhaHang(c *gin.Context) {
	var nhahang dto.NhaHang

	// 1️⃣ Bind
	if err := c.ShouldBind(&nhahang); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 2️⃣ Validate
	if nhahang.TenNhaHang == "" || nhahang.DiaChi == "" {
		c.JSON(400, gin.H{"error": "Thiếu dữ liệu"})
		return
	}
	if nhahang.SoTaiKhoan == 0 || nhahang.NganHang == "" {
		c.JSON(400, gin.H{"error": "Thiếu thông tin ngân hàng"})
		return
	}

	// 3️⃣ Lấy user
	userAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, gin.H{"error": "Chưa đăng nhập"})
		return
	}
	nhahang.MaNguoiDung = userAny.(uint)

	// 4️⃣ Save domain
	if err := config.DB.Create(&nhahang).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// =========================
	// 🧠 BUILD DOCUMENT
	// =========================
	document := fmt.Sprintf(
		"Nhà hàng: %s\nĐịa chỉ: %s\nNgân hàng: %s\nSố tài khoản: %d",
		nhahang.TenNhaHang,
		nhahang.DiaChi,
		nhahang.NganHang,
		nhahang.SoTaiKhoan,
	)

	// =========================
	// 🔥 EMBEDDING
	// =========================
	embedding, err := h.llm.Embed(c.Request.Context(), document)
	if err != nil {
		log.Println("embed error:", err)
	}

	// =========================
	// 📦 METADATA
	// =========================
	metaJSON, _ := json.Marshal(map[string]any{
		"type": "nha_hang",
		"id":   nhahang.MaNhaHang,
		"name": nhahang.TenNhaHang,
	})

	// =========================
	// 📦 VECTOR STRING
	// =========================
	vectorStr := vectorToString(embedding)

	// =========================
	// INSERT VECTOR DB (RAW SQL)
	// =========================
	if len(embedding) > 0 {
		result := config.DB.Exec(`
			INSERT INTO menu_embeddings (id, document, metadata, embedding)
			VALUES ($1, $2, $3, $4)
		`,
			fmt.Sprintf("nha_hang_%d", nhahang.MaNhaHang),
			document,
			string(metaJSON),
			vectorStr,
		)

		if result.Error != nil {
			log.Println("vector insert error:", result.Error)
		}
	}

	// =========================
	// 🖼️ UPLOAD IMAGE
	// =========================
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()

			uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
				Folder: "nhahang",
			})

			if err == nil {
				config.DB.Create(&dto.HinhAnh{
					OwnerID:   nhahang.MaNhaHang,
					OwnerType: "nha_hang",
					Url:       uploadResult.SecureURL,
				})
			}
		}
	}

	// =========================
	// RESPONSE
	// =========================
	config.DB.Preload("AnhNhaHang").First(&nhahang, nhahang.MaNhaHang)

	c.JSON(201, gin.H{
		"message": "Tạo nhà hàng + embedding thành công",
		"data":    nhahang,
	})
}
func GetAllNhaHang(c *gin.Context) {
	var list []dto.NhaHang

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
	var nhahang dto.NhaHang

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

	var nhaHangs []dto.NhaHang

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
	var nhahang dto.NhaHang

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
			Delete(&dto.HinhAnh{})

		src, _ := file.Open()
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
			Folder: "nhahang",
		})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		hinhAnh := dto.HinhAnh{
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
	var nhahang dto.NhaHang

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