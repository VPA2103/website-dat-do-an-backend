package models

import "time"

type ThanhToan struct {
	MaThanhToan uint `gorm:"primaryKey;size:10;autoIncrement" json:"ma_thanh_toan"`

	MaHD   uint   `gorm:"size:10" json:"ma_hd"`
	HoaDon HoaDon `gorm:"foreignKey:MaHD;references:MaHD" json:"hoa_don"`

	SoTien            float64     `json:"so_tien"`
	HinhThucThanhToan string    `json:"hinh_thuc_thanh_toan"`
	NgayThanhToan     time.Time `json:"ngay_thanh_toan"`


	MaNVThanhToan     string    `gorm:"size:10"` // foreign key
	NhanVienThanhToan NguoiDung `gorm:"foreignKey:MaNVThanhToan;references:MaNguoiDung"`
}
