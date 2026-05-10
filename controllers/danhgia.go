package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

type DanhGiaInput struct {
	MaNguoiDung uint   `json:"ma_nguoi_dung" binding:"required"`
	MaMonAn     uint   `json:"ma_mon_an" binding:"required"`
	SoSao       int    `json:"so_sao" binding:"required,min=1,max=5"`
	NoiDung     string `json:"noi_dung"`
}

func CreateDanhGia(c *gin.Context) {
	var input DanhGiaInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	dg := models.DanhGia{
		MaNguoiDung: input.MaNguoiDung,
		MaMonAn:     input.MaMonAn,
		SoSao:       input.SoSao,
		NoiDung:     input.NoiDung,
	}

	if err := config.DB.Create(&dg).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dg)
}

func GetDanhSachDanhGia(c *gin.Context) {
	var list []models.DanhGia

	if err := config.DB.
		Preload("NguoiDung").
		Preload("MonAn").
		Find(&list).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, list)
}

func GetDanhGiaByID(c *gin.Context) {
	id := c.Param("id")

	var dg models.DanhGia

	if err := config.DB.
		Preload("NguoiDung").
		Preload("MonAn").
		First(&dg, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy đánh giá"})
		return
	}

	c.JSON(200, dg)
}

func UpdateDanhGia(c *gin.Context) {
	id := c.Param("id")

	var dg models.DanhGia
	if err := config.DB.First(&dg, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy"})
		return
	}

	var input DanhGiaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	dg.SoSao = input.SoSao
	dg.NoiDung = input.NoiDung

	if err := config.DB.Save(&dg).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dg)
}

func DeleteDanhGia(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.DanhGia{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Đã xoá đánh giá"})
}