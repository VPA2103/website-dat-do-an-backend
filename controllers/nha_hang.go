package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

type UpdateNhaHangRequest struct {
	TenNhaHang string `form:"ten_nha_hang"`
	TrangThai  int    `form:"trang_thai"`
	DiaChi     string `form:"dia_chi"`

	SoTaiKhoan   string `form:"so_tai_khoan"`
	NganHang     string `form:"ngan_hang"`
	TenNguoiNhan string `form:"ten_nguoi_nhan"`

	GioMoCua   string `form:"gio_mo_cua"`
	GioDongCua string `form:"gio_dong_cua"`
	MoTa       string `form:"mo_ta"`
}

func (h *ChatHandler) CreateNhaHang(c *gin.Context) {
	var nhahang models.NhaHang

	if err := c.ShouldBind(&nhahang); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if nhahang.TenNhaHang == "" || nhahang.DiaChi == "" {
		c.JSON(400, gin.H{"error": "Thiếu dữ liệu"})
		return
	}

	// save DB
	if err := config.DB.Create(&nhahang).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// ================= EMBEDDING =================
	document := fmt.Sprintf(
		`Nhà hàng: %s
		Địa chỉ: %s
		Mô tả: %s
		Giờ mở cửa: %s
		Giờ đóng cửa: %s
		Ngân hàng: %s
		Số tài khoản: %d
		Tên người nhận: %s`,
		nhahang.TenNhaHang,
		nhahang.DiaChi,
		nhahang.MoTa,
		nhahang.GioMoCua,
		nhahang.GioDongCua,
		nhahang.NganHang,
		nhahang.SoTaiKhoan,
		nhahang.TenNguoiNhan,
	)

	embedding, err := h.llm.Embed(c.Request.Context(), document)
	if err == nil && len(embedding) > 0 {

		metaJSON, _ := json.Marshal(map[string]any{
			"type":    "nha_hang",
			"id":      nhahang.MaNhaHang,
			"name":    nhahang.TenNhaHang,
			"dia_chi": nhahang.DiaChi,
			"mo_ta":   nhahang.MoTa,
		})

		vectorStr := vectorToString(embedding)

		config.DB.Exec(`
			INSERT INTO menu_embeddings (id, document, metadata, embedding)
			VALUES ($1, $2, $3, $4)
		`,
			fmt.Sprintf("nha_hang_%d", nhahang.MaNhaHang),
			document,
			string(metaJSON),
			vectorStr,
		)
	}

	// ================= IMAGE =================
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, _ := file.Open()
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
			Folder: "nhahang",
		})

		if err == nil {
			config.DB.Create(&models.HinhAnh{
				OwnerID:   nhahang.MaNhaHang,
				OwnerType: "nha_hang",
				Url:       uploadResult.SecureURL,
			})
		}
	}

	config.DB.Preload("AnhNhaHang").First(&nhahang, nhahang.MaNhaHang)

	c.JSON(201, gin.H{
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

func (h *ChatHandler) UpdateNhaHang(c *gin.Context) {
	id := c.Param("id")

	var nhahang models.NhaHang

	// 1. FIND
	if err := config.DB.First(&nhahang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy nhà hàng",
		})
		return
	}

	// 3. BIND REQUEST (KHÔNG OVERWRITE ENTITY)
	var req UpdateNhaHangRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. UPDATE FIELD TAY
	nhahang.TenNhaHang = req.TenNhaHang
	nhahang.TrangThai = req.TrangThai
	nhahang.DiaChi = req.DiaChi
	nhahang.SoTaiKhoan = req.SoTaiKhoan
	nhahang.NganHang = req.NganHang
	nhahang.TenNguoiNhan = req.TenNguoiNhan
	nhahang.GioMoCua = req.GioMoCua
	nhahang.GioDongCua = req.GioDongCua
	nhahang.MoTa = req.MoTa

	// 5. HANDLE IMAGE
	file, err := c.FormFile("image")
	if err == nil && file != nil {

		// delete old image
		config.DB.
			Where("owner_id = ? AND owner_type = ?", nhahang.MaNhaHang, "nha_hang").
			Delete(&models.HinhAnh{})

		src, err := file.Open()
		if err == nil {
			defer src.Close()

			uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
				Folder: "nhahang",
			})

			if err != nil {
				c.JSON(500, gin.H{"error": "Upload ảnh thất bại"})
				return
			}

			config.DB.Create(&models.HinhAnh{
				Url:       uploadResult.SecureURL,
				OwnerID:   nhahang.MaNhaHang,
				OwnerType: "nha_hang",
			})
		}
	}

	// 6. SAVE DB
	if err := config.DB.Save(&nhahang).Error; err != nil {
		c.JSON(500, gin.H{"error": "Không thể cập nhật"})
		return
	}

	// 7. UPDATE EMBEDDING (QUAN TRỌNG)
	document := fmt.Sprintf(
		`Nhà hàng: %s
Địa chỉ: %s
Mô tả: %s
Giờ mở cửa: %s
Giờ đóng cửa: %s
Ngân hàng: %s
Số tài khoản: %d
Tên người nhận: %s`,
		nhahang.TenNhaHang,
		nhahang.DiaChi,
		nhahang.MoTa,
		nhahang.GioMoCua,
		nhahang.GioDongCua,
		nhahang.NganHang,
		nhahang.SoTaiKhoan,
		nhahang.TenNguoiNhan,
	)

	embedding, err := h.llm.Embed(c.Request.Context(), document)
	if err == nil && len(embedding) > 0 {

		metaJSON, _ := json.Marshal(map[string]any{
			"type":    "nha_hang",
			"id":      nhahang.MaNhaHang,
			"name":    nhahang.TenNhaHang,
			"dia_chi": nhahang.DiaChi,
			"mo_ta":   nhahang.MoTa,
		})

		vectorStr := vectorToString(embedding)

		config.DB.Exec(`
			UPDATE menu_embeddings
			SET document = $1,
			    metadata = $2,
			    embedding = $3
			WHERE id = $4
		`,
			document,
			string(metaJSON),
			vectorStr,
			fmt.Sprintf("nha_hang_%d", nhahang.MaNhaHang),
		)
	}

	// 8. RETURN FULL DATA
	config.DB.Preload("AnhNhaHang").First(&nhahang, nhahang.MaNhaHang)

	c.JSON(200, gin.H{
		"message": "Cập nhật thành công",
		"data":    nhahang,
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
