package models

import "time"

type DanhGia struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	MaHoaDon uint `gorm:"column:ma_hoa_don"`

	MaMonAn uint `gorm:"column:ma_mon_an"`
	MonAn   MonAn `gorm:"foreignKey:MaMonAn;references:MaMonAn" json:"mon_an"`

	MaNguoiDung uint `gorm:"column:ma_nguoi_dung"`
	NguoiDung   NguoiDung `gorm:"foreignKey:MaNguoiDung;references:MaNguoiDung" json:"nguoi_dung"`
	TrangThai          string  `gorm:"type:varchar(30);default:'hien'" json:"trang_thai"`

	SoSao   int    `json:"so_sao"`
	NoiDung string `json:"noi_dung"`

	CreatedAt time.Time `gorm:"column:ngay_danh_gia" json:"ngay_danh_gia"`
}
