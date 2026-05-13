package models

type ChiTietHoaDon struct {
	MaChiTiet uint `gorm:"primaryKey;autoIncrement" json:"ma_chi_tiet"`

	MaHoaDon uint  `json:"ma_hoa_don"`

	MonAn MonAn `gorm:"foreignKey:MaMonAn;references:MaMonAn" json:"mon_an"`
	MaMonAn  uint  `json:"ma_mon_an"`

	SoLuong   int     `json:"so_luong"`
	DonGia    float64 `json:"don_gia"`
	ThanhTien float64 `json:"thanh_tien"`
	TrangThai string  `json:"trang_thai"`
	GhiChu    string  `gorm:"type:text" json:"ghi_chu"`
}
