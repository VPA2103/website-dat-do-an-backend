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

	// lấy user từ token
	maNguoiDungAny, exists := c.Get("ma_nguoi_dung")
	if !exists {
		c.JSON(401, gin.H{"error": "Không tìm thấy người dùng trong token"})
		return
	}

	maNguoiDung := maNguoiDungAny.(uint)

	binhLuan := models.BinhLuan{
		MaNguoiDung: maNguoiDung,
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
		Preload("BinhLuans.BinhLuans").
		Preload("NguoiDung").
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
		Preload("NguoiDung").
		First(&binhLuan, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy bình luận"})
		return
	}

	c.JSON(200, gin.H{"data": binhLuan})
}

func UpdateBinhLuan(c *gin.Context) {
	id := c.Param("id")

	maNguoiDungAny, _ := c.Get("ma_nguoi_dung")
	maNguoiDung := maNguoiDungAny.(uint)

	var binhLuan models.BinhLuan

	// 1. lấy data trước
	if err := config.DB.First(&binhLuan, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy bình luận"})
		return
	}

	// 2. check quyền sau
	if binhLuan.MaNguoiDung != maNguoiDung {
		c.JSON(403, gin.H{"error": "Không có quyền sửa bình luận này"})
		return
	}

	var input UpdateBinhLuanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// 3. update
	if err := config.DB.Model(&binhLuan).
		Update("noi_dung", input.NoiDung).Error; err != nil {
		c.JSON(500, gin.H{"error": "Không thể cập nhật"})
		return
	}

	c.JSON(200, gin.H{"data": binhLuan})
}

func DeleteBinhLuan(c *gin.Context) {
	id := c.Param("id")
	maNguoiDungAny, _ := c.Get("ma_nguoi_dung")

	maNguoiDung := maNguoiDungAny.(uint)

	var binhLuan models.BinhLuan
	if err := config.DB.First(&binhLuan, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy bình luận"})
		return
	}

	if binhLuan.MaNguoiDung != maNguoiDung {
		c.JSON(403, gin.H{"error": "Không có quyền xóa bình luận này"})
		return
	}

	if err := config.DB.Delete(&binhLuan, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Không thể xóa"})
		return
	}

	c.JSON(200, gin.H{"message": "Đã xóa"})
}
