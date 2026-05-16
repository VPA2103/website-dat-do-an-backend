package models

import "time"

type DiaChi struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	MaNguoiDung uint `json:"ma_nguoi_dung"`

	HoTen string `json:"ho_ten"`
	SDT   string `json:"sdt"`

	DiaChi string `json:"dia_chi"`

	MacDinh bool `gorm:"default:false" json:"mac_dinh"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
