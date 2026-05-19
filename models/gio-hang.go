package models

import "time"

type GioHang struct {
	MaGioHang uint `gorm:"primaryKey;autoIncrement" json:"ma_gio_hang"`

	MaNguoiDung uint    `gorm:"uniqueIndex:idx_user_mon" json:"ma_nguoi_dung"`
	MaMonAn     uint    `gorm:"uniqueIndex:idx_user_mon" json:"ma_mon_an"`
	MonAn       []MonAn `gorm:"foreignKey:MaMonAn;references:MaMonAn" json:"mon_an"`
	SoLuong     int     `json:"so_luong"`
	GiaTien     int     `json:"gia_tien"`
	CreatedAt   time.Time

	Options []GioHangOption `gorm:"foreignKey:MaGioHang"`
}

type GioHangOption struct {
	MaGioHangOption uint `gorm:"primaryKey;autoIncrement"`

	MaGioHang    uint `json:"ma_gio_hang"`
	MaNhomOption uint `json:"ma_nhom_option"`
	MaOptionItem uint `json:"ma_option_item"`

	TenNhomOption string `json:"ten_nhom_option"`
	TenOption     string `json:"ten_option"`
	GiaThem       int    `json:"gia_them"`
}
