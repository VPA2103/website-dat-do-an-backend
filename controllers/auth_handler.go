package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/utils"

	//"github.com/vpa/quanlynhahang-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// Struct cho request login
type LoginInput struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng nhập email và mật khẩu"})
		return
	}

	// =========================
	// ✅ Check khách hàng trước
	// =========================
	var nd models.NguoiDung
	if err := config.DB.Where("email = ?", input.Email).First(&nd).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email không tồn tại"})
		return
	}

	// check password
	if bcrypt.CompareHashAndPassword([]byte(nd.MatKhau), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mật khẩu không đúng"})
		return
	}

	role := nd.LoaiNguoiDung
	redirect := "/user"
	if role == "admin" {
		redirect = "/admin"
	}

	token, err := utils.GenerateToken(nd.MaNguoiDung, nd.Email, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Đăng nhập thành công",
		"role":     role,
		"redirect": redirect,
		"token":    token,
		"data":     nd,
	})
}

func Register(c *gin.Context) {
	var input struct {
		HoTen   string `json:"name" form:"name" binding:"required"`
		Email   string `json:"email" form:"email" binding:"required,email"`
		MatKhau string `json:"password" form:"password" binding:"required"`
		SDT     string `json:"sdt" form:"sdt"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng nhập đầy đủ thông tin"})
		return
	}

	// Kiểm tra trùng email
	var existingKH models.NguoiDung
	if err := config.DB.Where("email = ?", input.Email).First(&existingKH).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email đã tồn tại trong hệ thống"})
		return
	}
	var existingNV models.NguoiDung
	if err := config.DB.Where("email = ?", input.Email).First(&existingNV).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email đã tồn tại trong hệ thống"})
		return
	}

	// Mã hoá mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.MatKhau), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi mã hoá mật khẩu"})
		return
	}

	newKH := models.NguoiDung{
		HoTen:         input.HoTen,
		Email:         input.Email,
		MatKhau:       string(hashedPassword),
		SDT:           input.SDT,
		LoaiNguoiDung: "user",
	}

	if err := config.DB.Create(&newKH).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := utils.GenerateToken(newKH.MaNguoiDung, newKH.Email, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng ký thành công",

		"redirect": "/user",
		"token":    token,
		"user": gin.H{
			"id":    newKH.MaNguoiDung,
			"hoten": newKH.HoTen,
			"email": newKH.Email,
			"sdt":   newKH.SDT,
			"role":  newKH.LoaiNguoiDung,
		},
	})
}

// Handler admin
func AdminDashboard(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome Admin Dashboard"})
}

func GetProfile(c *gin.Context) {
	c.JSON(200, gin.H{"message": "User profile"})
}
