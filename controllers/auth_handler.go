package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

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

type OTPData struct {
	Code      string
	ExpiredAt time.Time
}

var OTPStore = map[string]OTPData{}


func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

type ForgotPasswordInput struct {
	Email string `json:"email"`
}


func SendOTP(c *gin.Context) {

	var input ForgotPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	var user models.NguoiDung

	if err := config.DB.
		Where("email = ?", input.Email).
		First(&user).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Email không tồn tại",
		})
		return
	}

	otp := GenerateOTP()

	OTPStore[input.Email] = OTPData{
		Code: otp,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}

	body := `
		<h2>Mã OTP đổi mật khẩu</h2>
		<p>OTP của bạn là: <b>` + otp + `</b></p>
		<p>OTP có hiệu lực trong 5 phút</p>
	`

	err := utils.SendMail(
		input.Email,
		"OTP đổi mật khẩu",
		body,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gửi mail thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đã gửi OTP qua email",
	})
}

type ResetPasswordInput struct {
	Email       string `json:"email"`
	OTP         string `json:"otp"`
	MatKhauMoi  string `json:"mat_khau_moi"`
}

func ResetPassword(c *gin.Context) {

	var input ResetPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	otpData, exists := OTPStore[input.Email]

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "OTP không tồn tại",
		})
		return
	}

	// kiểm tra hết hạn
	if time.Now().After(otpData.ExpiredAt) {

		delete(OTPStore, input.Email)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "OTP đã hết hạn",
		})
		return
	}

	// kiểm tra otp
	if otpData.Code != input.OTP {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "OTP không đúng",
		})
		return
	}

	// hash mật khẩu mới
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.MatKhauMoi),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể hash mật khẩu",
		})
		return
	}

	// update password
	if err := config.DB.
		Model(&models.NguoiDung{}).
		Where("email = ?", input.Email).
		Update("mat_khau", string(hashedPassword)).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Đổi mật khẩu thất bại",
		})
		return
	}

	// xóa otp sau khi dùng
	delete(OTPStore, input.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Đổi mật khẩu thành công",
	})
}

// Handler admin
func AdminDashboard(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome Admin Dashboard"})
}

func GetProfile(c *gin.Context) {
	c.JSON(200, gin.H{"message": "User profile"})
}
