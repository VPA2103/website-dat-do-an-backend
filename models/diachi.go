package models

import "time"

type DiaChi struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	HoTen string `json:"ho_ten"`
	SDT   string `json:"sdt"`

	DiaChi string `json:"dia_chi"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	MacDinh bool `gorm:"default:false" json:"mac_dinh"`

	MaNguoiDung uint      `json:"ma_nguoi_dung"`
	NguoiDung *NguoiDung `gorm:"foreignKey:MaNguoiDung" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
