package dto

type ChiTietHoaDon struct {
	MaChiTiet uint `json:"ma_chi_tiet"`

	MaHoaDon uint `json:"ma_hoa_don"`

	MaMonAn uint   `json:"ma_mon_an"`
	MonAn   *MonAn `json:"mon_an,omitempty"`

	SoLuong   int     `json:"so_luong"`
	DonGia    float64 `json:"don_gia"`
	ThanhTien float64 `json:"thanh_tien"`

	Options []ChiTietHoaDonOption `json:"options,omitempty"`
}