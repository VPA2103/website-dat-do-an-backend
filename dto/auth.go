package dto

import "time"

type RegisterInput struct {
	HoTen       string `json:"ho_ten"`
	Email       string `json:"email"`
	MatKhau     string `json:"mat_khau"`
	SoDienThoai string `json:"so_dien_thoai"`
}

type RegisterOTPData struct {
	Code      string
	ExpiredAt time.Time
	UserData  RegisterInput
}

type VerifyRegisterOTPInput struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

var RegisterOTPStore = map[string]RegisterOTPData{}