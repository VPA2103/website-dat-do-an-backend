package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

type GioHangInput struct {
	MaMonAn uint `json:"ma_mon_an" binding:"required"`
	SoLuong int  `json:"so_luong" binding:"required"`
}

type UpdateSoLuongInput struct {
	SoLuong int `json:"so_luong" binding:"required"`
}

func AddToCart(c *gin.Context) {
	var input GioHangInput

	// validate json
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	// validate số lượng
	if input.SoLuong <= 0 {
		c.JSON(400, gin.H{
			"error": "Số lượng phải lớn hơn 0",
		})
		return
	}

	// lấy user từ token
	userAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"error": "Chưa đăng nhập",
		})
		return
	}

	userID, ok := userAny.(uint)
	if !ok {
		c.JSON(500, gin.H{
			"error": "Sai kiểu dữ liệu user",
		})
		return
	}

	// kiểm tra món ăn tồn tại
	var monAn models.MonAn

	if err := config.DB.
		Preload("AnhMonAn").
		Where("ma_mon_an = ?", input.MaMonAn).
		First(&monAn).Error; err != nil {

		c.JSON(404, gin.H{
			"error": "Không tìm thấy món ăn",
		})
		return
	}

	// kiểm tra đã tồn tại trong giỏ chưa
	var existing models.GioHang

	err := config.DB.
		Where(
			"ma_nguoi_dung = ? AND ma_mon_an = ?",
			userID,
			input.MaMonAn,
		).
		First(&existing).Error

	// =========================
	// ĐÃ TỒN TẠI -> UPDATE
	// =========================
	if err == nil {

		existing.SoLuong += input.SoLuong
		existing.GiaTien = existing.SoLuong * int(monAn.GiaTien)

		if err := config.DB.Save(&existing).Error; err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		// preload lại dữ liệu
		var result models.GioHang

		if err := config.DB.
			Preload("MonAn").
			Preload("MonAn.AnhMonAn").
			Where("ma_gio_hang = ?", existing.MaGioHang).
			First(&result).Error; err != nil {

			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Đã cập nhật số lượng",
			"data":    result,
		})

		return
	}

	// =========================
	// CHƯA TỒN TẠI -> CREATE
	// =========================
	gioHang := models.GioHang{
		MaNguoiDung: userID,
		MaMonAn:     input.MaMonAn,
		SoLuong:     input.SoLuong,
		GiaTien:     int(monAn.GiaTien) * input.SoLuong,
	}

	if err := config.DB.Create(&gioHang).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// preload lại dữ liệu
	var result models.GioHang

	if err := config.DB.
		Preload("MonAn").
		Preload("MonAn.AnhMonAn").
		Where("ma_gio_hang = ?", gioHang.MaGioHang).
		First(&result).Error; err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Thêm vào giỏ hàng thành công",
		"data":    result,
	})
}

func GetAllCart(c *gin.Context) {
	var list []models.GioHang

	if err := config.DB.
		Preload("MonAn.AnhMonAn").
		Find(&list).Error; err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, list)
}

func GetCartByUser(c *gin.Context) {
	userID := c.Param("id")

	// validate id
	if _, err := strconv.Atoi(userID); err != nil {
		c.JSON(400, gin.H{
			"error": "User id không hợp lệ",
		})
		return
	}

	var list []models.GioHang

	if err := config.DB.
		Where("ma_nguoi_dung = ?", userID).
		Preload("MonAn.AnhMonAn").
		Find(&list).Error; err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, list)
}

func UpdateSoLuongCart(c *gin.Context) {
	monID := c.Param("ma_mon_an")

	var input UpdateSoLuongInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	if input.SoLuong <= 0 {
		c.JSON(400, gin.H{
			"error": "Số lượng phải lớn hơn 0",
		})
		return
	}

	userAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"error": "Chưa đăng nhập",
		})
		return
	}

	userID := userAny.(uint)

	var cart models.GioHang

	if err := config.DB.
		Preload("MonAn").
		Preload("MonAn.AnhMonAn").
		Where(
			"ma_nguoi_dung = ? AND ma_mon_an = ?",
			userID,
			monID,
		).
		First(&cart).Error; err != nil {

		c.JSON(404, gin.H{
			"error": "Không tìm thấy sản phẩm trong giỏ",
		})
		return
	}

	var monAn models.MonAn

	config.DB.
		Where("ma_mon_an = ?", monID).
		First(&monAn)

	cart.SoLuong = input.SoLuong
	cart.GiaTien = int(monAn.GiaTien) * input.SoLuong

	if err := config.DB.Save(&cart).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, cart)
}

func DeleteCart(c *gin.Context) {
	monID := c.Param("ma_mon_an")

	userAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"error": "Chưa đăng nhập",
		})
		return
	}

	userID := userAny.(uint)

	result := config.DB.
		Where(
			"ma_nguoi_dung = ? AND ma_mon_an = ?",
			userID,
			monID,
		).
		Delete(&models.GioHang{})

	if result.Error != nil {
		c.JSON(500, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{
			"error": "Không tìm thấy sản phẩm trong giỏ",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Đã xoá khỏi giỏ hàng",
	})
}

func XoaGioHangNguoiDung(c *gin.Context) {

	maNguoiDungAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Không xác thực người dùng"})
		return
	}

	maNguoiDung := maNguoiDungAny.(uint)

	if err := config.DB.
		Where("ma_nguoi_dung = ?", maNguoiDung).
		Delete(&models.GioHang{}).Error; err != nil {

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Đã xóa giỏ hàng"})
}
