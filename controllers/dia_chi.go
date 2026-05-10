package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

type input struct {
	HoTen   string `form:"ho_ten" json:"ho_ten"`
	SDT     string `form:"sdt" json:"sdt"`
	DiaChi  string `form:"dia_chi" json:"dia_chi"`
	MacDinh bool   `form:"mac_dinh" json:"mac_dinh"`
}

func CreateDiaChi(c *gin.Context) {
	var dc models.DiaChi

	if err := c.ShouldBindJSON(&dc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	dc.CreatedAt = time.Now()
	dc.UpdatedAt = time.Now()

	// nếu là địa chỉ mặc định thì reset các địa chỉ khác
	if dc.MacDinh {
		config.DB.Model(&models.DiaChi{}).
			Where("ma_nguoi_dung = ?", dc.MaNguoiDung).
			Update("mac_dinh", false)
	}

	if err := config.DB.Create(&dc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể tạo địa chỉ",
		})
		return
	}

	c.JSON(http.StatusOK, dc)
}

func GetDiaChiByUser(c *gin.Context) {
	maNguoiDung := c.Param("ma_nguoi_dung")

	var list []models.DiaChi

	if err := config.DB.
		Where("ma_nguoi_dung = ?", maNguoiDung).
		Order("mac_dinh DESC").
		Find(&list).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không lấy được dữ liệu",
		})
		return
	}

	c.JSON(http.StatusOK, list)
}

func GetDiaChiByID(c *gin.Context) {
	id := c.Param("id")

	var dc models.DiaChi

	if err := config.DB.First(&dc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy địa chỉ",
		})
		return
	}

	c.JSON(http.StatusOK, dc)
}

func UpdateDiaChi(c *gin.Context) {
	id := c.Param("id")

	var dc models.DiaChi
	if err := config.DB.First(&dc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy địa chỉ",
		})
		return
	}

	var input input
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	// update từng field an toàn
	dc.HoTen = input.HoTen
	dc.SDT = input.SDT
	dc.DiaChi = input.DiaChi
	dc.MacDinh = input.MacDinh
	dc.UpdatedAt = time.Now()

	// nếu set mặc định -> reset các địa chỉ khác
	if input.MacDinh {
		config.DB.Model(&models.DiaChi{}).
			Where("ma_nguoi_dung = ? AND id != ?", dc.MaNguoiDung, dc.ID).
			Update("mac_dinh", false)
	}

	if err := config.DB.Save(&dc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể cập nhật",
		})
		return
	}

	c.JSON(http.StatusOK, dc)
}

func DeleteDiaChi(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.DiaChi{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể xóa",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa thành công",
	})
}
func SetDiaChiMacDinh(c *gin.Context) {
	id := c.Param("id")

	var dc models.DiaChi

	// kiểm tra địa chỉ tồn tại
	if err := config.DB.First(&dc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy địa chỉ",
		})
		return
	}

	// transaction để đảm bảo dữ liệu đồng bộ
	tx := config.DB.Begin()

	// reset tất cả địa chỉ của user
	if err := tx.Model(&models.DiaChi{}).
		Where("ma_nguoi_dung = ?", dc.MaNguoiDung).
		Update("mac_dinh", false).Error; err != nil {

		tx.Rollback()

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể cập nhật địa chỉ",
		})
		return
	}

	// set địa chỉ hiện tại thành mặc định
	if err := tx.Model(&dc).
		Updates(map[string]interface{}{
			"mac_dinh":  true,
			"updated_at": time.Now(),
		}).Error; err != nil {

		tx.Rollback()

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể đặt mặc định",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Đặt địa chỉ mặc định thành công",
		"data":    dc,
	})
}