package models

import "time"

type HoaDon struct {
	MaHD      uint   `gorm:"primaryKey;size:10;autoIncrement"`
	HoTen     string `json:"ho_ten"`
	SDT       string `json:"sdt"`
	DiaChi    string `json:"dia_chi"`
	GhiChu    string `json:"ghi_chu"`
	Ngay      time.Time
	TongTien  float64
	TrangThai string `gorm:"type:varchar(30);default:'cho_xac_nhan'" json:"trang_thai"`

	// MaNVOrder      *uint           `gorm:"size:10"`
	// NhanVienOrder  *NhanVien       `gorm:"foreignKey:MaNVOrder;references:MaNV"`
	ChiTietHoaDons []ChiTietHoaDon `gorm:"foreignKey:MaHoaDon"`
}
