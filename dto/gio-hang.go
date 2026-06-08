package dto

type GioHang struct {
	MaGioHang   uint
	MaNguoiDung uint
	MaMonAn     uint

	MonAn MonAn

	SoLuong int
	GiaTien int

	Options []GioHangOption
}

type GioHangOption struct {
	MaGioHangOption uint

	MaGioHang    uint `json:"ma_gio_hang"`
	MaNhomOption uint `json:"ma_nhom_option"`
	MaOptionItem uint `json:"ma_option_item"`

	TenNhomOption string     `json:"ten_nhom_option"`
	TenOption     string     `json:"ten_option"`
	GiaThem       int        `json:"gia_them"`
	OptionItem    OptionItem
}