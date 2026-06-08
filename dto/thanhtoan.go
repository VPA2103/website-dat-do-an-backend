package dto

import "time"

type ThanhToan struct {
	MaThanhToan uint `gorm:"primaryKey;autoIncrement" json:"ma_thanh_toan"`

	MaHD uint `gorm:"not null" json:"ma_hd"`

	HoaDon HoaDon `gorm:"foreignKey:MaHD;references:MaHD" json:"hoa_don"`

	SoTien float64 `gorm:"not null" json:"so_tien"`

	HinhThucThanhToan string `gorm:"type:varchar(50)" json:"hinh_thuc_thanh_toan"`

	NgayThanhToan time.Time `json:"ngay_thanh_toan"`

	MaNVThanhToan string `gorm:"size:10" json:"ma_nv_thanh_toan"`

	NhanVienThanhToan NguoiDung `gorm:"foreignKey:MaNVThanhToan;references:MaNguoiDung" json:"nhan_vien_thanh_toan"`
}
