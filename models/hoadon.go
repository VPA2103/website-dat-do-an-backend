package models

import "time"

type HoaDon struct {
	MaHD        uint      `gorm:"primaryKey;size:10;autoIncrement" json:"ma_hd"`
	MaNguoiDung uint      `json:"ma_nguoi_dung"`
	HoTen       string    `json:"ho_ten"`
	SDT         string    `json:"sdt"`
	DiaChi      string    `json:"dia_chi"`
	GhiChu      string    `json:"ghi_chu"`
	Ngay        time.Time `json:"ngay"`

	TongTien  float64 `json:"tong_tien"` // sau giảm
	TamTinh   float64 `json:"tam_tinh"`  // trước giảm
	TienGiam  float64 `json:"tien_giam"`
	TrangThai string  `gorm:"type:varchar(30);default:'cho_xac_nhan'" json:"trang_thai"`

	// MaNVOrder      *uint           `gorm:"size:10"`
	// NhanVienOrder  *NhanVien       `gorm:"foreignKey:MaNVOrder;references:MaNV"`
	GiamGiaID      *uint           `json:"giam_gia_id"`
	GiamGia        GiamGia         `gorm:"foreignKey:GiamGiaID;references:ID" json:"giam_gia"`
	ChiTietHoaDons []ChiTietHoaDon `gorm:"foreignKey:MaHoaDon" json:"chi_tiet_hoa_dons"`
}
