package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/dto"
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

	// Gửi mail chào mừng (bất đồng bộ, không block response)
	go func() {
		mailInfo := utils.DangKyMailInfo{
			TenKhachHang: newKH.HoTen,
			Email:        newKH.Email,
			MaNguoiDung:  newKH.MaNguoiDung,
		}
		if err := utils.SendMailSauKhiDangKy(newKH.Email, mailInfo); err != nil {
			log.Printf("Gửi mail đăng ký thất bại cho %s: %v", newKH.Email, err)
		}
	}()

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
		Code:      otp,
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
	Email      string `json:"email"`
	OTP        string `json:"otp"`
	MatKhauMoi string `json:"mat_khau_moi"`
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

func SendRegisterOTP(c *gin.Context) {

	var input dto.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	// Kiểm tra email tồn tại chưa
	var user models.NguoiDung

	err := config.DB.
		Where("email = ?", input.Email).
		First(&user).Error

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email đã tồn tại",
		})
		return
	}

	otp := GenerateOTP()

	dto.RegisterOTPStore[input.Email] = dto.RegisterOTPData{
		Code:      otp,
		ExpiredAt: time.Now().Add(5 * time.Minute),
		UserData:  input,
	}

	body := `<!DOCTYPE html>
				<html lang="vi">
				<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
				<body style="margin:0;padding:20px;background:#f4ede0;font-family:Georgia,'Times New Roman',serif;">
				<div style="max-width:520px;margin:0 auto;background:#ffffff;border-radius:8px;overflow:hidden;">

					<!-- Header -->
					<div style="background:#1a1a1a;padding:32px 32px 24px;text-align:center;">
					<div style="font-size:22px;letter-spacing:4px;color:#e8d5b0;">✦ NHÀ HÀNG ✦</div>
					<div style="font-size:11px;letter-spacing:6px;color:#8a7a5a;margin-top:4px;font-family:'Courier New',monospace;">SAIGON KITCHEN</div>
					</div>

					<!-- Title -->
					<div style="background:#f7f0e3;padding:24px 32px 16px;text-align:center;border-bottom:1px solid #e0d0b0;">
					<div style="font-size:11px;letter-spacing:5px;color:#8a7a5a;font-family:'Courier New',monospace;margin-bottom:8px;">XÁC THỰC</div>
					<div style="font-size:22px;color:#2a1f0a;letter-spacing:1px;">Mã OTP của bạn</div>
					<div style="width:40px;height:1px;background:#c4a55a;margin:12px auto 0;"></div>
					</div>

					<!-- Body -->
					<div style="padding:24px 32px;background:#fdfaf4;">
					<p style="font-size:14px;color:#4a3c20;line-height:1.8;margin:0 0 16px;">
						Kính gửi quý khách,
					</p>
					<p style="font-size:14px;color:#4a3c20;line-height:1.8;margin:0 0 24px;">
						Chúng tôi nhận được yêu cầu xác thực tài khoản của bạn tại <strong>Saigon Kitchen</strong>. Vui lòng sử dụng mã OTP dưới đây:
					</p>

					<!-- OTP Box -->
					<div style="text-align:center;background:#fff;border:0.5px solid #e0d0b0;border-radius:8px;padding:28px 14px;margin-bottom:24px;">
						<div style="font-size:11px;letter-spacing:4px;color:#8a7a5a;font-family:'Courier New',monospace;margin-bottom:12px;">MÃ XÁC THỰC</div>
						<div style="font-size:36px;font-weight:700;letter-spacing:10px;color:#1a1a1a;font-family:'Courier New',monospace;">` + otp + `</div>
						<div style="width:40px;height:1px;background:#c4a55a;margin:16px auto 0;"></div>
						<div style="font-size:12px;color:#8a7a5a;margin-top:10px;font-family:'Courier New',monospace;">Hiệu lực trong 5 phút</div>
					</div>

					<!-- Ghi chú -->
					<div style="background:#f7f0e3;border-left:3px solid #c4a55a;padding:12px 14px;margin-bottom:24px;">
						<p style="font-size:13px;color:#5a4520;margin:0;line-height:1.6;">
						Vui lòng <strong>không chia sẻ</strong> mã OTP này với bất kỳ ai, kể cả nhân viên của chúng tôi.
						</p>
					</div>

					<p style="font-size:13px;color:#6a5a3a;line-height:1.8;margin:0;">
						Nếu bạn không thực hiện yêu cầu này, vui lòng bỏ qua email hoặc liên hệ ngay với chúng tôi để được hỗ trợ.
					</p>
					</div>

					<!-- Footer -->
					<div style="padding:16px 32px;background:#1a1a1a;text-align:center;">
					<p style="font-size:11px;color:#6a5a3a;margin:0;letter-spacing:1px;font-family:'Courier New',monospace;">
						123 Đường ABC, Q.1, TP.HCM &nbsp;|&nbsp; 028-xxxx-xxxx
					</p>
					</div>

				</div>
				</body>
				</html>`

	err = utils.SendMail(
		input.Email,
		"Xác nhận đăng ký tài khoản",
		body,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gửi mail thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đã gửi OTP xác nhận",
	})
}

func VerifyRegisterOTP(c *gin.Context) {

	var input dto.VerifyRegisterOTPInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	data, exists := dto.RegisterOTPStore[input.Email]

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "OTP không tồn tại",
		})
		return
	}

	// Check hết hạn
	if time.Now().After(data.ExpiredAt) {

		delete(dto.RegisterOTPStore, input.Email)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "OTP đã hết hạn",
		})
		return
	}

	// Check OTP
	if data.Code != input.OTP {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "OTP không đúng",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(data.UserData.MatKhau),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lỗi hash password",
		})
		return
	}

	// Tạo user
	user := models.NguoiDung{
		HoTen:   data.UserData.HoTen,
		Email:   data.UserData.Email,
		MatKhau: string(hashedPassword),
		SDT:     data.UserData.SoDienThoai,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể tạo tài khoản",
		})
		return
	}

	// Xóa OTP
	delete(dto.RegisterOTPStore, input.Email)

	// Gửi mail chào mừng
	utils.SendMailSauKhiDangKy(user.Email, utils.DangKyMailInfo{
		TenKhachHang: user.HoTen,
		MaNguoiDung:  user.MaNguoiDung,
		Email:        user.Email,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng ký thành công",
	})
}

// Handler admin
func AdminDashboard(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome Admin Dashboard"})
}

func GetProfile(c *gin.Context) {
	c.JSON(200, gin.H{"message": "User profile"})
}
