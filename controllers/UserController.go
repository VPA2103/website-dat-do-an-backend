package controllers

import (
	//"fmt"
	//"log"
	//"math"
	"net/http"
	//"strconv"
	//"time"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	//"github.com/vpa/quanlynhahang-backend/dto"
	//"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/models"
	//"github.com/vpa/quanlynhahang-backend/utils"
	//"gorm.io/gorm"
	//"github.com/xuri/excelize/v2"
)

func GetTongTienDaMua(c *gin.Context) {

	// lấy user_id từ middleware
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

	var tongTien float64

	err := config.DB.
		Model(&models.HoaDon{}).
		Where("ma_nguoi_dung = ?", maNguoiDung).
		Where("trang_thai_thanh_toan = ?", "da_thanh_toan").
		Select("COALESCE(SUM(tong_tien), 0)").
		Scan(&tongTien).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể lấy tổng tiền",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tong_tien_da_mua": tongTien,
	})
}