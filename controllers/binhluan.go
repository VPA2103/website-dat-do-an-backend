package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

type CreateBinhLuanInput struct {
	MaNguoiDung uint   `json:"ma_nguoi_dung" binding:"required"`
	MaMonAn     uint   `json:"ma_mon_an" binding:"required"`
	NoiDung     string `json:"noi_dung" binding:"required"`
	MaCha       *uint  `json:"ma_cha"`
}

type UpdateBinhLuanInput struct {
	NoiDung string `json:"noi_dung"`
}

func CreateBinhLuan(c *gin.Context) {
	var input CreateBinhLuanInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	binhLuan := models.BinhLuan{
		MaNguoiDung: input.MaNguoiDung,
		MaMonAn:     input.MaMonAn,
		NoiDung:     input.NoiDung,
		MaCha:       input.MaCha,
	}

	if err := config.DB.Create(&binhLuan).Error; err != nil {
		c.JSON(500, gin.H{"error": "Không thể tạo bình luận"})
		return
	}

	c.JSON(200, gin.H{"data": binhLuan})
}

func GetBinhLuanByMonAn(c *gin.Context) {
	maMon := c.Param("ma_mon_an")

	var binhLuans []models.BinhLuan

	err := config.DB.
		Where("ma_mon_an = ? AND ma_cha IS NULL", maMon).
		Preload("BinhLuans").
		Find(&binhLuans).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Lỗi khi lấy dữ liệu"})
		return
	}

	c.JSON(200, gin.H{"data": binhLuans})
}

func GetBinhLuanByID(c *gin.Context) {
	id := c.Param("id")

	var binhLuan models.BinhLuan

	if err := config.DB.
		Preload("BinhLuans").
		First(&binhLuan, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy bình luận"})
		return
	}

	c.JSON(200, gin.H{"data": binhLuan})
}

func UpdateBinhLuan(c *gin.Context) {
	id := c.Param("id")

	var binhLuan models.BinhLuan

	if err := config.DB.First(&binhLuan, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy bình luận"})
		return
	}

	var input UpdateBinhLuanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	config.DB.Model(&binhLuan).Update("noi_dung", input.NoiDung)

	c.JSON(200, gin.H{"data": binhLuan})
}

func DeleteBinhLuan(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.BinhLuan{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Không thể xóa"})
		return
	}

	c.JSON(200, gin.H{"message": "Đã xóa"})
}