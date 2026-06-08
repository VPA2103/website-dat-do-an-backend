package dto

type ChiTietHoaDonOption struct {
	MaCTHDOption uint `json:"ma_cthd_option"`

	MaChiTiet uint `json:"ma_chi_tiet"`

	MaOptionItem uint         `json:"ma_option_item"`
	OptionItem   *OptionItem  `json:"option_item,omitempty"`

	TenOption string  `json:"ten_option"`
	GiaThem   float64 `json:"gia_them"`
}