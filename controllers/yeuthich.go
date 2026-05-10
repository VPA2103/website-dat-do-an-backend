package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

type YeuThichInput struct {
	MaNguoiDung uint `json:"ma_nguoi_dung" binding:"required"`
	MaMonAn     uint `json:"ma_mon_an" binding:"required"`
}

func CreateYeuThich(c *gin.Context) {
	var input YeuThichInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	yt := models.YeuThich{
		MaNguoiDung: input.MaNguoiDung,
		MaMonAn:     input.MaMonAn,
	}

	// tránh duplicate (unique index)
	if err := config.DB.Create(&yt).Error; err != nil {
		c.JSON(500, gin.H{"error": "Món đã được yêu thích hoặc lỗi DB"})
		return
	}

	c.JSON(200, yt)
}

func GetAllYeuThich(c *gin.Context) {
	var list []models.YeuThich

	if err := config.DB.
		Preload("NguoiDung").
		Preload("MonAn").
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
		Preload("MonAn").
		Find(&list).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, list)
}

func DeleteYeuThich(c *gin.Context) {
	userID := c.Query("user")
	monID := c.Query("mon")

	if err := config.DB.
		Where("ma_nguoi_dung = ? AND ma_mon_an = ?", userID, monID).
		Delete(&models.YeuThich{}).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Đã xoá yêu thích"})
}