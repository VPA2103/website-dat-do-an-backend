package models

import "time"

type HoaDon struct {
	MaHD      uint `gorm:"primaryKey;size:10;autoIncrement"`
	MaBan     uint `gorm:"size:10"`
	Ngay   time.Time
	TongTien  float64
	TrangThai string

	// MaNVOrder      *uint           `gorm:"size:10"`
	// NhanVienOrder  *NhanVien       `gorm:"foreignKey:MaNVOrder;references:MaNV"`
	ChiTietHoaDons []ChiTietHoaDon `gorm:"foreignKey:MaHoaDon"`
}
