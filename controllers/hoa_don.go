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
	HoTen  string         `json:"ho_ten"`
	SDT    string         `json:"sdt"`
	DiaChi string         `json:"dia_chi"`
	GhiChu string         `json:"ghi_chu"`
	MonAns []MonDatInput  `json:"mon_ans"`
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

	// transaction
	tx := config.DB.Begin()

	hoaDon := models.HoaDon{
		HoTen:     input.HoTen,
		SDT:       input.SDT,
		DiaChi:    input.DiaChi,
		GhiChu:    input.GhiChu,
		Ngay:      time.Now(),
		TrangThai: "cho_xac_nhan",
		TongTien:  0,
	}

	// tạo hóa đơn
	if err := tx.Create(&hoaDon).Error; err != nil {
		tx.Rollback()

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể tạo hóa đơn",
		})
		return
	}

	var tongTien float64

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

		thanhTien := float64(item.SoLuong) * monAn.GiaTien

		chiTiet := models.ChiTietHoaDon{
			MaHoaDon:  hoaDon.MaHD,
			MaMonAn:   item.MaMonAn,
			SoLuong:   item.SoLuong,
			DonGia:    monAn.GiaTien,
			ThanhTien: thanhTien,
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

		tongTien += thanhTien
	}

	// cập nhật tổng tiền
	if err := tx.Model(&hoaDon).
		Update("tong_tien", tongTien).Error; err != nil {

		tx.Rollback()

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể cập nhật tổng tiền",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Đặt đồ ăn thành công",
		"data":    hoaDon,
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
