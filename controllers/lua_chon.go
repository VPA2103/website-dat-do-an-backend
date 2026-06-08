package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/dto"
)

// API tạo nhóm option


type UpdateNhomOptionRequest struct {
	TenNhom        string `json:"ten_nhom"`
	BatBuoc        bool   `json:"bat_buoc"`
	ChonNhieu      bool   `json:"chon_nhieu"`
	SoLuongToiDa   int    `json:"so_luong_toi_da"`
	SoLuongToiThieu int   `json:"so_luong_toi_thieu"`
}

func CreateNhomOption(c *gin.Context) {
	var input dto.NhomOption

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// validate
	if input.MaMonAn == 0 {
		c.JSON(400, gin.H{
			"error": "Thiếu mã món ăn",
		})
		return
	}

	if input.TenNhom == "" {
		c.JSON(400, gin.H{
			"error": "Tên nhóm không được để trống",
		})
		return
	}

	// kiểm tra món ăn tồn tại
	var monan dto.MonAn
	if err := config.DB.First(&monan, input.MaMonAn).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Món ăn không tồn tại",
		})
		return
	}

	input.TrangThai = 1

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Không thể tạo nhóm option",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Tạo nhóm option thành công",
		"data":    input,
	})
}

func GetAllNhomOption(c *gin.Context) {

	var nhoms []dto.NhomOption

	err := config.DB.
		Preload("OptionItems").
		Find(&nhoms).Error

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Không thể lấy danh sách nhóm option",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": nhoms,
	})
}

func GetNhomOptionByID(c *gin.Context) {

	id := c.Param("id")

	var nhom dto.NhomOption

	err := config.DB.
		Preload("OptionItems").
		First(&nhom, id).Error

	if err != nil {
		c.JSON(404, gin.H{
			"error": "Không tìm thấy nhóm option",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": nhom,
	})
}

func UpdateNhomOption(c *gin.Context) {
	id := c.Param("id")

	var nhom dto.NhomOption
	if err := config.DB.First(&nhom, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy nhóm option"})
		return
	}

	var input UpdateNhomOptionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	nhom.TenNhom = input.TenNhom
	nhom.BatBuoc = input.BatBuoc
	nhom.ChonNhieu = input.ChonNhieu
	nhom.SoLuongToiDa = input.SoLuongToiDa
	nhom.SoLuongToiThieu = input.SoLuongToiThieu

	if err := config.DB.Save(&nhom).Error; err != nil {
		c.JSON(500, gin.H{"error": "Cập nhật thất bại"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Cập nhật thành công",
		"data":    nhom,
	})
}

func DeleteNhomOption(c *gin.Context) {
	id := c.Param("id")

	var nhom dto.NhomOption

	if err := config.DB.First(&nhom, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy nhóm option"})
		return
	}

	// 1. xóa option items trước
	if err := config.DB.Where("ma_nhom_option = ?", id).
		Delete(&dto.OptionItem{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "Không thể xóa option items"})
		return
	}

	// 2. xóa nhóm
	if err := config.DB.Delete(&nhom).Error; err != nil {
		c.JSON(500, gin.H{"error": "Không thể xóa nhóm option"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Xóa nhóm option thành công",
	})
}
// API tạo option item

func CreateOptionItem(c *gin.Context) {
	var input dto.OptionItem

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if input.MaNhomOption == 0 {
		c.JSON(400, gin.H{
			"error": "Thiếu mã nhóm option",
		})
		return
	}

	if input.TenOption == "" {
		c.JSON(400, gin.H{
			"error": "Tên option không được để trống",
		})
		return
	}

	// kiểm tra nhóm option tồn tại
	var nhom dto.NhomOption
	if err := config.DB.First(&nhom, input.MaNhomOption).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Nhóm option không tồn tại",
		})
		return
	}

	input.TrangThai = 1

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Không thể tạo option item",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Tạo option item thành công",
		"data":    input,
	})
}

func GetAllOptionItem(c *gin.Context) {

	var items []dto.OptionItem

	if err := config.DB.Find(&items).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Không thể lấy option item",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": items,
	})
}

func GetOptionItemByID(c *gin.Context) {

	id := c.Param("id")

	var item dto.OptionItem

	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Không tìm thấy option item",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": item,
	})
}

func UpdateOptionItem(c *gin.Context) {

	id := c.Param("id")

	var item dto.OptionItem

	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Không tìm thấy option item",
		})
		return
	}

	var input dto.OptionItem

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	item.TenOption = input.TenOption
	item.GiaThem = input.GiaThem

	if err := config.DB.Save(&item).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Cập nhật thất bại",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Cập nhật thành công",
		"data":    item,
	})
}

func DeleteOptionItem(c *gin.Context) {

	id := c.Param("id")

	var item dto.OptionItem

	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Không tìm thấy option item",
		})
		return
	}

	config.DB.Delete(&item)

	c.JSON(200, gin.H{
		"message": "Xóa option item thành công",
	})
}

