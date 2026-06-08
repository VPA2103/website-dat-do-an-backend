package dto

import "time"

type DiaChi struct {
	ID uint `json:"id"`

	HoTen string `json:"ho_ten"`
	SDT   string `json:"sdt"`

	DiaChi string `json:"dia_chi"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	MacDinh bool `json:"mac_dinh"`

	MaNguoiDung uint `json:"ma_nguoi_dung"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DiaChiInput struct {
	HoTen    string  `json:"ho_ten" binding:"required"`
	SDT      string  `json:"sdt" binding:"required"`
	DiaChi   string  `json:"dia_chi" binding:"required"`
	Latitude float64 `json:"latitude" `
	Longitude float64 `json:"longitude" `
	MacDinh  bool    `json:"mac_dinh"`
}