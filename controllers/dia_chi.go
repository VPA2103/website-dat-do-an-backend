package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/dto"
	"github.com/vpa/quanlynhahang-backend/models"
)

func CreateDiaChi(c *gin.Context) {

	var input dto.DiaChiInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	dc := models.DiaChi{
		HoTen:       input.HoTen,
		SDT:         input.SDT,
		DiaChi:      input.DiaChi,
		// Latitude:    input.Latitude,
		// Longitude:   input.Longitude,
		MacDinh:     input.MacDinh,
		MaNguoiDung: userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if dc.MacDinh {
		config.DB.Model(&models.DiaChi{}).
			Where("ma_nguoi_dung = ?", userID).
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

	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	if dc.MaNguoiDung != userID {
		c.JSON(403, gin.H{
			"error": "Không có quyền",
		})
		return
	}

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

	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	if dc.MaNguoiDung != userID {
		c.JSON(403, gin.H{
			"error": "Không có quyền",
		})
		return
	}

	var input dto.DiaChiInput
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
	dc.Latitude = input.Latitude
	dc.Longitude = input.Longitude
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
			"mac_dinh":   true,
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

func GetGoogleMapDirection(c *gin.Context) {

	id := c.Param("id")

	var diaChi models.DiaChi

	if err := config.DB.First(&diaChi, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Không tìm thấy địa chỉ",
		})
		return
	}

	// tọa độ quán
	shopLat := 10.762622
	shopLng := 106.660172

	mapURL := fmt.Sprintf(
		"https://www.google.com/maps/dir/%f,%f/%f,%f",
		shopLat,
		shopLng,
		diaChi.Latitude,
		diaChi.Longitude,
	)

	c.JSON(200, gin.H{
		"map_url": mapURL,
	})
}
