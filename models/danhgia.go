package models

import "time"

type DanhGia struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	MaHoaDon    uint `gorm:"column:ma_hoa_don"`
	MaNguoiDung uint `gorm:"column:ma_nguoi_dung"`
	MaMonAn     uint `gorm:"column:ma_mon_an"`

	SoSao     int       `json:"so_sao"`
	NoiDung   string    `json:"noi_dung"`
	NguoiDung NguoiDung `gorm:"foreignKey:MaNguoiDung;references:MaNguoiDung" json:"nguoi_dung"`
	CreatedAt time.Time `gorm:"column:ngay_danh_gia" json:"ngay_danh_gia"`
}
