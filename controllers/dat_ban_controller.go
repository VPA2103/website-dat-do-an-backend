package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/utils"
)

func CreateDatBan(c *gin.Context) {
	var input models.DatBan
	userID, _ := c.Get("user_id")

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ÉP logic nghiệp vụ
	datban := models.DatBan{
		TenKhachHang: input.TenKhachHang,
		Email:        input.Email,
		SDT:          input.SDT,
		GhiChu:       input.GhiChu,
		MaBanAn:      input.MaBanAn,
		Ngay:         input.Ngay,
		Gio:          input.Gio,
		TrangThai:    "dang_xu_ly",
		MaNguoiDung:  userID.(uint),
		// IDNhanVienXacNhan = nil
	}

	if err := config.DB.Create(&datban).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo đặt bàn"})
		return
	}

	// 🔔 GỬI MAIL SAU KHI ĐẶT BÀN THÀNH CÔNG
	go func(db models.DatBan) {
		var ban models.BanAn
		config.DB.First(&ban, db.MaBanAn)

		if err := utils.SendMailDatBan(db.Email, utils.DatBanMailInfo{
			TenKhachHang: db.TenKhachHang,
			MaDatBan:     db.MaDatBan,
			Ngay:         db.Ngay,
			Gio:          db.Gio,
			TenBan:       ban.TenBan,
			Email:        db.Email,
			GhiChu:       db.GhiChu,
		}); err != nil {
			log.Println("Send mail dat ban error:", err)
		}
	}(datban)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Đặt bàn thành công",
		"data":    datban,
	})
}

func GetAllDatBan(c *gin.Context) {
	var datbans []models.DatBan

	if err := config.DB.Find(&datbans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể lấy danh sách đặt bàn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": datbans,
	})
}

func GetDatBanByID(c *gin.Context) {
	id := c.Param("id")
	var datban models.DatBan

	if err := config.DB.Preload("NhanVienXacNhan").First(&datban, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy đặt bàn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": datban,
	})
}

func UpdateDatBan(c *gin.Context) {
	id := c.Param("id")
	var datban models.DatBan

	if err := config.DB.First(&datban, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy đặt bàn"})
		return
	}

	var input struct {
		TenKhachHang string `json:"ten_khach_hang"`
		SDT          string `json:"sdt"`
		GhiChu       string `json:"ghi_chu"`
		Ngay         string `json:"ngay"`
		Gio          string `json:"gio"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&datban).Updates(input)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật đặt bàn thành công",
		"data":    datban,
	})
}

func XacNhanDatBan(c *gin.Context) {
	id := c.Param("id")

	var datban models.DatBan
	if err := config.DB.First(&datban, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy đặt bàn"})
		return
	}

	// tránh xác nhận lại
	if datban.TrangThai == "da_xac_nhan" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Đặt bàn đã được xác nhận"})
		return
	}

	// update trạng thái
	if err := config.DB.Model(&datban).Updates(map[string]interface{}{
		"trang_thai": "da_xac_nhan",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xác nhận đặt bàn"})
		return
	}

	// load lại thông tin bàn (optional)
	var ban models.BanAn
	config.DB.First(&ban, datban.MaBanAn)

	// gửi email async
	go func(db models.DatBan, tenBan string) {
		err := utils.SendMailDatBanXacNhan(db.Email, utils.DatBanXacNhanMailInfo{
			TenKhachHang: db.TenKhachHang,
			MaDatBan:     db.MaDatBan,
			Ngay:         db.Ngay,
			Gio:          db.Gio,
			TenBan:       tenBan,
			Email:        db.Email,
			GhiChu:       db.GhiChu,
		})

		if err != nil {
			log.Println("Send mail xác nhận đặt bàn lỗi:", err)
		}
	}(datban, ban.TenBan)

	c.JSON(http.StatusOK, gin.H{
		"message": "Xác nhận đặt bàn thành công",
	})
}

func DeleteDatBan(c *gin.Context) {
	id := c.Param("id")
	var datban models.DatBan

	if err := config.DB.First(&datban, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy đặt bàn",
		})
		return
	}

	config.DB.Delete(&datban)

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa đặt bàn thành công",
	})
}
func GetDatBanCuaNguoiDung(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Không xác thực người dùng",
		})
		return
	}

	userID := userIDRaw.(uint) // hoặc uint64 / int tùy bạn lưu

	var datbans []models.DatBan

	if err := config.DB.
		Where("ma_nguoi_dung = ?", userID).
		Find(&datbans).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể lấy danh sách đặt bàn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": datbans,
	})
}
func HuyDatBan(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Không xác thực người dùng",
		})
		return
	}

	var datban models.DatBan
	if err := config.DB.First(&datban, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy đặt bàn",
		})
		return
	}

	// chỉ chủ đặt bàn mới được hủy
	if datban.MaNguoiDung != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Bạn không có quyền hủy đặt bàn này",
		})
		return
	}

	// không cho hủy nếu đã xác nhận
	if datban.TrangThai == "da_xac_nhan" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Đặt bàn đã được xác nhận, không thể hủy",
		})
		return
	}

	// update trạng thái hủy
	if err := config.DB.Model(&datban).Updates(map[string]interface{}{
		"trang_thai": "da_huy",
	}).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể hủy đặt bàn",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hủy đặt bàn thành công",
	})
}
