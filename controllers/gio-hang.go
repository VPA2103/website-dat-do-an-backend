package controllers

import (
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

type GioHangInput struct {
	MaMonAn uint                 `json:"ma_mon_an" binding:"required"`
	SoLuong int                  `json:"so_luong" binding:"required"`
	Options []GioHangOptionInput `json:"options"`
}
type GioHangOptionInput struct {
	MaNhomOption uint `json:"ma_nhom_option" binding:"required"`
	MaOptionItem uint `json:"ma_option_item" binding:"required"`
}



type UpdateSoLuongInput struct {
	SoLuong int `json:"so_luong" binding:"required"`
}

func AddToCart(c *gin.Context) {
	var input GioHangInput

	// =======================
	// VALIDATE INPUT
	// =======================
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	if input.SoLuong <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Số lượng phải > 0"})
		return
	}

	// =======================
	// GET USER
	// =======================
	userAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Chưa đăng nhập"})
		return
	}
	userID := userAny.(uint)

	// =======================
	// CHECK MÓN ĂN
	// =======================
	var monAn models.MonAn
	if err := config.DB.
		Where("ma_mon_an = ?", input.MaMonAn).
		First(&monAn).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy món ăn"})
		return
	}

	tx := config.DB.Begin()

	// =======================
	// TÍNH GIÁ OPTION
	// =======================
	var giaOption float64
	var optionItems []models.OptionItem

	for _, opt := range input.Options {
		var item models.OptionItem
		if err := tx.
			Preload("NhomOption").
			Where("ma_option_item = ?", opt.MaOptionItem).
			First(&item).Error; err != nil {

			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Option không hợp lệ"})
			return
		}

		giaOption += item.GiaThem
		optionItems = append(optionItems, item)
	}

	// =======================
	// TÍNH GIÁ CUỐI
	// =======================
	donGia := float64(monAn.GiaTien) + giaOption
	thanhTien := donGia * float64(input.SoLuong)

	// =======================
	// CREATE GIO HANG
	// =======================
	gioHang := models.GioHang{
		MaNguoiDung: userID,
		MaMonAn:     monAn.MaMonAn,
		SoLuong:     input.SoLuong,
		GiaTien:     int(thanhTien), // ✅ ÉP KIỂU
	}

	if err := tx.Create(&gioHang).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// =======================
	// INSERT OPTIONS
	// =======================
	for _, item := range optionItems {
		row := models.GioHangOption{
			MaGioHang:     gioHang.MaGioHang,
			MaNhomOption: item.MaNhomOption,
			MaOptionItem: item.MaOptionItem,
			TenNhomOption: item.NhomOption.TenNhom, // ✅ ĐÚNG FIELD
			TenOption:     item.TenOption,
			GiaThem:       int(item.GiaThem), // ✅ ÉP KIỂU
		}

		if err := tx.Create(&row).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	// =======================
	// RESPONSE
	// =======================
	var result models.GioHang
	config.DB.
		Preload("MonAn").
		Preload("Options").
		Where("ma_gio_hang = ?", gioHang.MaGioHang).
		First(&result)

	c.JSON(http.StatusCreated, gin.H{
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

	if _, err := strconv.Atoi(userID); err != nil {
		c.JSON(400, gin.H{
			"message": "User id không hợp lệ",
		})
		return
	}

	var list []models.GioHang

	if err := config.DB.
		Where("ma_nguoi_dung = ?", userID).
		Preload("MonAn").
		Preload("MonAn.AnhMonAn").
		Preload("Options").
		Preload("Options.OptionItem").
		Preload("Options.OptionItem.NhomOption"). // ✅ ĐÚNG
		Find(&list).Error; err != nil {

		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": list,
	})
}

func UpdateSoLuongCart(c *gin.Context) {
	cartID := c.Param("ma_gio_hang") // ✅ FIX HERE

	var input UpdateSoLuongInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	if input.SoLuong <= 0 {
		c.JSON(400, gin.H{"error": "Số lượng phải lớn hơn 0"})
		return
	}

	userAny, _ := c.Get("user_id")
	userID := userAny.(uint)

	var cart models.GioHang

	if err := config.DB.
		Where("ma_gio_hang = ? AND ma_nguoi_dung = ?", cartID, userID).
		First(&cart).Error; err != nil {

		c.JSON(404, gin.H{"error": "Không tìm thấy item"})
		return
	}

	var monAn models.MonAn
	config.DB.Where("ma_mon_an = ?", cart.MaMonAn).First(&monAn)

	cart.SoLuong = input.SoLuong
	cart.GiaTien = int(monAn.GiaTien) * input.SoLuong

	if err := config.DB.
		Model(&cart).
		Updates(map[string]interface{}{
			"so_luong": cart.SoLuong,
			"gia_tien": cart.GiaTien,
		}).Error; err != nil {

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "OK",
		"data": cart,
	})
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
