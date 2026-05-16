package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/internal/dto"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/models"
	"gorm.io/gorm"
)

type HoaDonController struct {
	Hub *websocket.Hub
}

func NewHoaDonController(hub *websocket.Hub) *HoaDonController {
	return &HoaDonController{
		Hub: hub,
	}
}

type MonDatInput struct {
	MaMonAn uint   `json:"ma_mon_an"`
	SoLuong int    `json:"so_luong"`
	GhiChu  string `json:"ghi_chu"`
}

type DatDoAnInput struct {
	HoTen       string        `json:"ho_ten"`
	SDT         string        `json:"sdt"`
	DiaChi      string        `json:"dia_chi"`
	GhiChu      string        `json:"ghi_chu"`
	CodeGiamGia string        `json:"code_giam_gia"`
	MonAns      []MonDatInput `json:"mon_ans"`
}

func (ctrl *HoaDonController) DatDoAn(c *gin.Context) {

	var input DatDoAnInput

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	// lấy user từ middleware
	maNguoiDungAny, exists := c.Get("user_id")

	if !exists {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Vui lòng đăng nhập",
		})
		return
	}

	maNguoiDung, ok := maNguoiDungAny.(uint)

	if !ok {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user_id không hợp lệ",
		})
		return
	}

	// validate input
	if input.HoTen == "" ||
		input.SDT == "" ||
		input.DiaChi == "" {

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

	// rollback nếu panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var tongTienServer float64

	// tạo hóa đơn
	hoaDon := models.HoaDon{
		MaNguoiDung: maNguoiDung,
		HoTen:       input.HoTen,
		SDT:         input.SDT,
		DiaChi:      input.DiaChi,
		GhiChu:      input.GhiChu,
		Ngay:        time.Now(),
		TrangThai:   "cho_xac_nhan",
	}

	if err := tx.Create(&hoaDon).Error; err != nil {

		tx.Rollback()

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể tạo hóa đơn",
		})
		return
	}

	// thêm món ăn
	for _, item := range input.MonAns {

		if item.SoLuong <= 0 {
			continue
		}

		var monAn models.MonAn

		if err := tx.
			First(&monAn, "ma_mon_an = ?", item.MaMonAn).Error; err != nil {

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
		}

		if err := tx.Create(&chiTiet).Error; err != nil {

			tx.Rollback()

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Không thể thêm món ăn",
			})
			return
		}
	}

	// =========================
	// xử lý mã giảm giá
	// =========================

	var tienGiam float64
	var giamGia models.GiamGia

	if input.CodeGiamGia != "" {

		err := tx.
			Where("code = ?", input.CodeGiamGia).
			First(&giamGia).Error

		if err != nil {

			tx.Rollback()

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Mã giảm giá không tồn tại",
			})
			return
		}

		// check active
		if !giamGia.IsActive {

			tx.Rollback()

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Mã giảm giá đã bị khóa",
			})
			return
		}

		now := time.Now()

		// check thời gian
		if now.Before(giamGia.NgayBatDau) ||
			now.After(giamGia.NgayKetThuc) {

			tx.Rollback()

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Mã giảm giá đã hết hạn",
			})
			return
		}

		// check đơn tối thiểu
		if tongTienServer < giamGia.DonToiThieu {

			tx.Rollback()

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Chưa đủ giá trị đơn tối thiểu",
			})
			return
		}

		// check giới hạn sử dụng
		if giamGia.GioiHanSuDung != nil &&
			giamGia.SoLanDaDung >= *giamGia.GioiHanSuDung {

			tx.Rollback()

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Mã giảm giá đã hết lượt sử dụng",
			})
			return
		}

		// tính giảm giá
		switch giamGia.LoaiGiamGia {

		case "percent":

			tienGiam =
				tongTienServer *
					giamGia.GiaTriGiam / 100

			// giới hạn giảm tối đa
			if giamGia.GiamToiDa > 0 &&
				tienGiam > giamGia.GiamToiDa {

				tienGiam = giamGia.GiamToiDa
			}

		case "fixed":

			tienGiam = giamGia.GiaTriGiam
		}

		// tránh âm tiền
		if tienGiam > tongTienServer {
			tienGiam = tongTienServer
		}

		// gắn voucher vào hóa đơn
		hoaDon.GiamGiaID = &giamGia.ID

		// tăng số lần sử dụng
		if err := tx.Model(&giamGia).
			Update(
				"so_lan_da_dung",
				gorm.Expr("so_lan_da_dung + ?", 1),
			).Error; err != nil {

			tx.Rollback()

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Không thể cập nhật mã giảm giá",
			})
			return
		}
	}

	// tổng cuối
	tongCuoi := tongTienServer - tienGiam

	// update hóa đơn
	updateData := map[string]interface{}{
		"tong_tien":   tongCuoi,
		"tam_tinh":    tongTienServer,
		"tien_giam":   tienGiam,
		"giam_gia_id": hoaDon.GiamGiaID,
	}

	if err := tx.
		Model(&hoaDon).
		Updates(updateData).Error; err != nil {

		tx.Rollback()

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể cập nhật hóa đơn",
		})
		return
	}

	// commit
	if err := tx.Commit().Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể lưu hóa đơn",
		})
		return
	}

	// lấy kết quả cuối
	var result models.HoaDon

	if err := config.DB.
		Preload("GiamGia").
		Preload("ChiTietHoaDons").
		First(&result, "ma_hd = ?", hoaDon.MaHD).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể lấy hóa đơn",
		})
		return
	}

	// realtime cho admin
	ctrl.Hub.BroadcastToRoom(0, dto.WSMessage{
		Type:    "new_hoa_don",
		Payload: result,
	})

	// realtime cho user
	ctrl.Hub.BroadcastToRoom(maNguoiDung, dto.WSMessage{
		Type:    "new_hoa_don_user",
		Payload: result,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Đặt đồ ăn thành công",
		"data":    result,
	})
}

func (ctrl *HoaDonController) XoaHoaDon(c *gin.Context) {

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

func (ctrl *HoaDonController) GetHoaDons(c *gin.Context) {

	var hoaDons []models.HoaDon

	if err := config.DB.
		Preload("ChiTietHoaDons").
		Preload("ChiTietHoaDons.MonAn").
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

func (ctrl *HoaDonController) GetHoaDonByID(c *gin.Context) {

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

func (ctrl *HoaDonController) UpdateTrangThaiHoaDon(c *gin.Context) {

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

	ctrl.Hub.BroadcastToRoom(0, dto.WSMessage{
		Type:    "update_trang_thai_hoa_don",
		Payload: hoaDon,
	})

	ctrl.Hub.BroadcastToRoom(hoaDon.MaNguoiDung, dto.WSMessage{
		Type:    "update_trang_thai_hoa_don_user",
		Payload: hoaDon,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật trạng thái thành công",
		"hoa_don": hoaDon,
	})
}

func (ctrl *HoaDonController) HuyHoaDon(c *gin.Context) {

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

	config.DB.First(&hoaDon, "ma_hd = ?", id)

	ctrl.Hub.BroadcastToRoom(0, dto.WSMessage{
		Type:    "cancel_hoa_don",
		Payload: hoaDon,
	})

	ctrl.Hub.BroadcastToRoom(hoaDon.MaNguoiDung, dto.WSMessage{
		Type:    "cancel_hoa_don_user",
		Payload: hoaDon,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Hủy hóa đơn thành công",
	})
}

func (ctrl *HoaDonController) GetHoaDonByTrangThai(c *gin.Context) {

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

func (ctrl *HoaDonController) UpdateHoaDon(c *gin.Context) {

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

func (ctrl *HoaDonController) GetHoaDonByNguoiDung(c *gin.Context) {

	maNguoiDungAny, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Chưa đăng nhập",
		})
		return
	}

	maNguoiDung := maNguoiDungAny.(uint)

	var hoaDons []models.HoaDon

	if err := config.DB.
		Where("ma_nguoi_dung = ?", maNguoiDung).
		Preload("ChiTietHoaDons").
		Preload("ChiTietHoaDons.MonAn").
		Order("ma_hd DESC").
		Find(&hoaDons).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": hoaDons,
	})
}
