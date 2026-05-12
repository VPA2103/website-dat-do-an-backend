package models

import "time"

type GioHang struct {
	MaGioHang uint `gorm:"primaryKey;autoIncrement" json:"ma_gio_hang"`

	MaNguoiDung uint  `gorm:"uniqueIndex:idx_user_mon" json:"ma_nguoi_dung"`
	MaMonAn     uint  `gorm:"uniqueIndex:idx_user_mon" json:"ma_mon_an"`
	MonAn       []MonAn `gorm:"foreignKey:MaMonAn;references:MaMonAn" json:"mon_an"`
	SoLuong     int   `json:"so_luong"`
	GiaTien   int     `json:"gia_tien"`
	CreatedAt   time.Time
}
