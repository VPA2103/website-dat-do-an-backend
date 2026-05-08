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
	var kh models.NguoiDung
	if err := config.DB.Where("email = ?", input.Email).First(&kh).Error; err == nil {
		// So sánh mật khẩu (đã mã hoá)
		if err := bcrypt.CompareHashAndPassword([]byte(kh.MatKhau), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Mật khẩu không đúng"})
			return
		}

		// Tạo token
		token, err := utils.GenerateToken(kh.MaNguoiDung, kh.Email, "user")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Đăng nhập thành công",
			"role":     "user",
			"redirect": "/user",
			"token":    token,
			"data":     kh,
		})
		return
	}

	// =========================
	// ✅ Check nhân viên
	// =========================
	var nv models.NguoiDung
	if err := config.DB.Where("email = ?", input.Email).First(&nv).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email không tồn tại"})
		return
	}

	// Nếu mật khẩu không được mã hoá thì so sánh trực tiếp
	// Nếu đã mã hoá thì dùng bcrypt.CompareHashAndPassword
	if bcrypt.CompareHashAndPassword([]byte(nv.MatKhau), []byte(input.Password)) != nil &&
		nv.MatKhau != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mật khẩu không đúng"})
		return
	}

	redirect := "/account"
	if nv.LoaiNguoiDung == "admin" {
		redirect = "/admin"
	} else if nv.LoaiNguoiDung == "user" {
		redirect = "/user/home"
	}

	token, err := utils.GenerateToken(nv.MaNguoiDung, nv.Email, nv.LoaiNguoiDung)
	//token, err := utils.generateToken(nv.MaNV, nv.Email, nv.LoaiNhanVien)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Đăng nhập thành công",
		"role":     nv.LoaiNguoiDung,
		"redirect": redirect,
		"token":    token,
		"data":     nv,
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
		HoTen:   input.HoTen,
		Email:   input.Email,
		MatKhau: string(hashedPassword),
		SDT:     input.SDT,
	}

	if err := config.DB.Create(&newKH).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo tài khoản"})
		return
	}

	token, err := utils.GenerateToken(newKH.MaNguoiDung, newKH.Email, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Đăng ký thành công",
		"role":     "user",
		"redirect": "/user",
		"token":    token,
		"user": gin.H{
			"id":    newKH.MaNguoiDung,
			"hoten": newKH.HoTen,
			"email": newKH.Email,
			"sdt":   newKH.SDT,
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
