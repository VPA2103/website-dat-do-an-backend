package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/services/send_mail"
)

func CreateDatBan(c *gin.Context) {
	var input models.DatBan

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
		// IDNhanVienXacNhan = nil
	}

	if err := config.DB.Create(&datban).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo đặt bàn"})
		return
	}

	// 🔔 GỬI MAIL SAU KHI ĐẶT BÀN THÀNH CÔNG
	go func(email string) {
		if err := send_mail.SendDatBanMail(email); err != nil {
			log.Println("❌ Gửi mail thất bại:", err)
		}
	}(datban.Email)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Đặt bàn thành công",
		"data":    datban,
	})
}

func GetAllDatBan(c *gin.Context) {
	var datbans []models.DatBan

	if err := config.DB.Preload("NhanVienXacNhan").Find(&datbans).Error; err != nil {
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

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không xác định được nhân viên"})
		return
	}
	nhanVienID := userID.(uint)

	var datban models.DatBan
	if err := config.DB.First(&datban, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy đặt bàn"})
		return
	}

	if datban.IDNhanVienXacNhan != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Đặt bàn đã được xác nhận"})
		return
	}

	config.DB.Model(&datban).Updates(map[string]interface{}{
		"id_nhan_vien_xac_nhan": nhanVienID,
		"trang_thai":            "da_xac_nhan",
	})

	c.JSON(http.StatusOK, gin.H{"message": "Xác nhận đặt bàn thành công"})
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
