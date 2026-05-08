package models

import "time"

type ThanhToan struct {
	MaThanhToan       string `gorm:"primaryKey;size:10;autoIncrement" json:"ma_thanh_toan"`
	MaHD              string `gorm:"size:10"`
	SoTien            float64
	HinhThucThanhToan string
	NgayThanhToan     time.Time
	GioThanhToan      time.Time

	MaNVThanhToan     string   `gorm:"size:10"` // foreign key
	NhanVienThanhToan NguoiDung `gorm:"foreignKey:MaNVThanhToan;references:MaNV"`
}
