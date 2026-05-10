package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/services"
)

type MonDatInput struct {
	MaMonAn uint   `json:"ma_mon_an"`
	SoLuong int    `json:"so_luong"`
	GhiChu  string `json:"ghi_chu"`
}

type DatDoAnInput struct {
	HoTen    string        `json:"ho_ten"`
	SDT      string        `json:"sdt"`
	DiaChi   string        `json:"dia_chi"`
	GhiChu   string        `json:"ghi_chu"`
	TongTien float64       `json:"tong_tien"`
	MonAns   []MonDatInput `json:"mon_ans"`
}

func DatDoAn(c *gin.Context) {

	var input DatDoAnInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	// validate
	if input.HoTen == "" || input.SDT == "" || input.DiaChi == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Thiếu thông tin khách hàng",
		})
		return
	}

	if len(input.MonAns) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Chưa có món ăn",
		})
		return
	}

	tx := config.DB.Begin()

	hoaDon := models.HoaDon{
		HoTen:     input.HoTen,
		SDT:       input.SDT,
		DiaChi:    input.DiaChi,
		GhiChu:    input.GhiChu,
		Ngay:      time.Now(),
		TrangThai: "cho_xac_nhan",
		TongTien:  input.TongTien, // lấy từ FE
	}

	// tạo hóa đơn
	if err := tx.Create(&hoaDon).Error; err != nil {

		tx.Rollback()

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể tạo hóa đơn",
		})
		return
	}

	// thêm món
	for _, item := range input.MonAns {

		if item.SoLuong <= 0 {
			continue
		}

		var monAn models.MonAn

		if err := tx.First(&monAn, "ma_mon_an = ?", item.MaMonAn).Error; err != nil {

			tx.Rollback()

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Món ăn không tồn tại",
			})
			return
		}

		chiTiet := models.ChiTietHoaDon{
			MaHoaDon:  hoaDon.MaHD,
			MaMonAn:   item.MaMonAn,
			SoLuong:   item.SoLuong,
			DonGia:    monAn.GiaTien,
			TrangThai: "cho_xac_nhan",
			GhiChu:    item.GhiChu,
		}

		if err := tx.Create(&chiTiet).Error; err != nil {

			tx.Rollback()

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Không thể thêm món ăn",
			})
			return
		}
	}

	tx.Commit()

	var result models.HoaDon

	if err := config.DB.
		Preload("ChiTietHoaDons").
		First(&result, "ma_hd = ?", hoaDon.MaHD).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể lấy hóa đơn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đặt đồ ăn thành công",
		"data":    result,
	})
}

func XoaHoaDon(c *gin.Context) {

	id := c.Param("id")

	var hoaDon models.HoaDon

	// kiểm tra hóa đơn tồn tại
	if err := config.DB.
		First(&hoaDon, "ma_hd = ?", id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Hóa đơn không tồn tại",
		})
		return
	}

	tx := config.DB.Begin()

	// xóa chi tiết hóa đơn trước
	if err := tx.
		Where("ma_hoa_don = ?", id).
		Delete(&models.ChiTietHoaDon{}).Error; err != nil {

		tx.Rollback()

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể xóa chi tiết hóa đơn",
		})
		return
	}

	// xóa hóa đơn
	if err := tx.
		Delete(&hoaDon).Error; err != nil {

		tx.Rollback()

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể xóa hóa đơn",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa hóa đơn thành công",
	})
}

func ThanhToanHoaDon(c *gin.Context) {

	var req struct {
		MaBan uint `json:"ma_ban"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	if err := services.CloseHoaDon(req.MaBan); err != nil {
		c.JSON(500, gin.H{"error": "Không thể thanh toán"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Thanh toán thành công",
	})
}
