package models

import "time"

type ThanhToan struct {
	MaThanhToan uint `gorm:"primaryKey;size:10;autoIncrement" json:"ma_thanh_toan"`

	MaHD   uint   `gorm:"size:10"`
	HoaDon HoaDon `gorm:"foreignKey:MaHD;references:MaHD" json:"hoa_don"`

	SoTien            float64
	HinhThucThanhToan string
	NgayThanhToan     time.Time
	GioThanhToan      time.Time

	MaNVThanhToan     string    `gorm:"size:10"` // foreign key
	NhanVienThanhToan NguoiDung `gorm:"foreignKey:MaNVThanhToan;references:MaNguoiDung"`
}
