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

	// backend tự tính
	var tongTienServer float64

	// tạo hóa đơn trước
	hoaDon := models.HoaDon{
		HoTen:     input.HoTen,
		SDT:       input.SDT,
		DiaChi:    input.DiaChi,
		GhiChu:    input.GhiChu,
		Ngay:      time.Now(),
		TrangThai: "cho_xac_nhan",
	}

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

		thanhTien := monAn.GiaTien * float64(item.SoLuong)

		tongTienServer += thanhTien

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
	}

	// so sánh FE và BE
	if input.TongTien != tongTienServer {

		tx.Rollback()

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tổng tiền không hợp lệ",
		})
		return
	}

	// cập nhật tổng tiền thật từ backend
	if err := tx.Model(&hoaDon).
		Update("tong_tien", tongTienServer).Error; err != nil {

		tx.Rollback()

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể cập nhật tổng tiền",
		})
		return
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

func GetHoaDons(c *gin.Context) {

	var hoaDons []models.HoaDon

	if err := config.DB.
		Preload("ChiTietHoaDons").
		Order("ma_hd DESC").
		Find(&hoaDons).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể lấy danh sách hóa đơn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": hoaDons,
	})
}

func GetHoaDonByID(c *gin.Context) {

	id := c.Param("id")

	var hoaDon models.HoaDon

	if err := config.DB.
		Preload("ChiTietHoaDons").
		First(&hoaDon, "ma_hd = ?", id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy hóa đơn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": hoaDon,
	})
}

func UpdateTrangThaiHoaDon(c *gin.Context) {

	id := c.Param("id")

	var input struct {
		TrangThai string `json:"trang_thai"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	// Validate trạng thái hợp lệ
	validTrangThai := map[string]bool{
		"cho_xac_nhan":  true,
		"dang_chuan_bi": true,
		"dang_giao":     true,
		"da_giao":       true,
		"da_huy":        true,
		"da_thanh_toan": true,
	}

	if !validTrangThai[input.TrangThai] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Trạng thái không hợp lệ",
			"trang_thai_hop_le": []string{
				"cho_xac_nhan",
				"dang_chuan_bi",
				"dang_giao",
				"da_giao",
				"da_huy",
				"da_thanh_toan",
			},
		})
		return
	}

	var hoaDon models.HoaDon

	if err := config.DB.First(&hoaDon, "ma_hd = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy hóa đơn",
		})
		return
	}

	if err := config.DB.Model(&hoaDon).
		Update("trang_thai", input.TrangThai).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể cập nhật trạng thái",
		})
		return
	}

	if err := config.DB.First(&hoaDon, "ma_hd = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể tải lại hóa đơn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật trạng thái thành công",
		"hoa_don": hoaDon,
	})
}

func HuyHoaDon(c *gin.Context) {

	id := c.Param("id")

	var hoaDon models.HoaDon

	if err := config.DB.
		First(&hoaDon, "ma_hd = ?", id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy hóa đơn",
		})
		return
	}

	if hoaDon.TrangThai == "da_giao" {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Không thể hủy hóa đơn đã giao",
		})
		return
	}

	if err := config.DB.Model(&hoaDon).
		Update("trang_thai", "da_huy").Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể hủy hóa đơn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hủy hóa đơn thành công",
	})
}

func GetHoaDonByTrangThai(c *gin.Context) {

	trangThai := c.Query("trang_thai")

	var hoaDons []models.HoaDon

	if err := config.DB.
		Where("trang_thai = ?", trangThai).
		Preload("ChiTietHoaDons").
		Find(&hoaDons).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể lấy hóa đơn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": hoaDons,
	})
}

func UpdateHoaDon(c *gin.Context) {

	id := c.Param("id")

	var input struct {
		HoTen  string `json:"ho_ten"`
		SDT    string `json:"sdt"`
		DiaChi string `json:"dia_chi"`
		GhiChu string `json:"ghi_chu"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	var hoaDon models.HoaDon

	if err := config.DB.
		First(&hoaDon, "ma_hd = ?", id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy hóa đơn",
		})
		return
	}

	if err := config.DB.Model(&hoaDon).Updates(models.HoaDon{
		HoTen:  input.HoTen,
		SDT:    input.SDT,
		DiaChi: input.DiaChi,
		GhiChu: input.GhiChu,
	}).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể cập nhật hóa đơn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật hóa đơn thành công",
	})
}

func ThanhToanHoaDon(c *gin.Context) {

	id := c.Param("ma_hd")

	var hoaDon models.HoaDon

	// check tồn tại
	if err := config.DB.First(&hoaDon, "ma_hd = ?", id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Hóa đơn không tồn tại"})
		return
	}

	// gọi service (nên truyền ID luôn)
	if err := services.CloseHoaDon(hoaDon.MaHD); err != nil {
		c.JSON(500, gin.H{"error": "Không thể thanh toán"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Xác nhận đã thanh toán",
	})
}
