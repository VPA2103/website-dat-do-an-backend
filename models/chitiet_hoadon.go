package models

type ChiTietHoaDon struct {
	MaChiTiet uint `gorm:"primaryKey;autoIncrement"`
	MaHoaDon      uint `gorm:"size:10"`
	MaMonAn   uint `gorm:"size:10"`
	SoLuong   int
	DonGia    float64
	ThanhTien float64
	TrangThai string
	GhiChu    string `gorm:"type:text"`
}
