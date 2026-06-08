package models

import "time"

type DanhGia struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	MaHoaDon    uint `gorm:"column:ma_hoa_don"`
	MaMonAn     uint `gorm:"column:ma_mon_an"`
	
	SoSao     int       `json:"so_sao"`
	NoiDung   string    `json:"noi_dung"`

	MaNguoiDung uint `gorm:"column:ma_nguoi_dung"`
	NguoiDung NguoiDung `gorm:"foreignKey:MaNguoiDung;" json:"nguoi_dung"`

	CreatedAt time.Time `gorm:"column:ngay_danh_gia" json:"ngay_danh_gia"`
}
