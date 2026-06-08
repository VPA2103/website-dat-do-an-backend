package dto

type MonAn struct {
	MaMonAn     uint    `json:"ma_mon_an"`
	MaLoaiMonAn uint    `json:"ma_loai_mon_an"`
	TenMonAn    string  `json:"ten_mon_an"`
	GiaTien     float64 `json:"gia_tien"`
	TrangThai   uint    `json:"trang_thai"`
	MoTa        string  `json:"mo_ta"`

	AnhMonAn    []HinhAnh    `json:"anh_mon_an,omitempty"`
	NhomOptions []NhomOption `json:"nhom_options,omitempty"`
}

type CreateMonAnRequest struct {
	MaLoaiMonAn uint    `json:"ma_loai_mon_an" binding:"required"`
	TenMonAn    string  `json:"ten_mon_an" binding:"required"`
	GiaTien     float64 `json:"gia_tien" binding:"required"`
	MoTa        string  `json:"mo_ta"`
}
