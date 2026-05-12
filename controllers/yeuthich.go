package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

type YeuThichInput struct {
	MaMonAn uint `json:"ma_mon_an" binding:"required"`
}

func AddMonAnYeuThich(c *gin.Context) {
	var input YeuThichInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	maNguoiDungAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Không tìm thấy user từ token"})
		return
	}

	maNguoiDung := maNguoiDungAny.(uint)

	// 🔥 CHECK TRÙNG TRƯỚC KHI INSERT
	var existing models.YeuThich
	err := config.DB.Where(
		"ma_nguoi_dung = ? AND ma_mon_an = ?",
		maNguoiDung,
		input.MaMonAn,
	).First(&existing).Error

	if err == nil {
		// đã tồn tại
		c.JSON(200, gin.H{"message": "Đã tồn tại trong yêu thích"})
		return
	}

	yt := models.YeuThich{
		MaNguoiDung: maNguoiDung,
		MaMonAn:     input.MaMonAn,
	}

	if err := config.DB.Create(&yt).Error; err != nil {
		c.JSON(500, gin.H{"error": "Lỗi DB"})
		return
	}

	c.JSON(200, yt)
}

func GetAllYeuThich(c *gin.Context) {
	var list []models.YeuThich

	if err := config.DB.
		Preload("NguoiDung").
		Preload("MonAn.AnhMonAn").
		Find(&list).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, list)
}

func GetYeuThichByUser(c *gin.Context) {
	userID := c.Param("id")

	var list []models.YeuThich

	if err := config.DB.
		Where("ma_nguoi_dung = ?", userID).
		Preload("MonAn.AnhMonAn").
		Find(&list).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, list)
}

func DeleteYeuThich(c *gin.Context) {
	monID := c.Param("ma_mon_an")

	// lấy user từ token
	maNguoiDungAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Không tìm thấy user từ token"})
		return
	}

	maNguoiDung, ok := maNguoiDungAny.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "Sai kiểu dữ liệu user"})
		return
	}

	// chỉ xoá món của chính user đó
	if err := config.DB.
		Where("ma_nguoi_dung = ? AND ma_mon_an = ?", maNguoiDung, monID).
		Delete(&models.YeuThich{}).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Đã xoá yêu thích"})
}
