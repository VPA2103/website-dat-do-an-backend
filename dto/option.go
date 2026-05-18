package dto

type CreateNhomOptionInput struct {
	MaMonAn uint `json:"ma_mon_an"`

	TenNhom string `json:"ten_nhom"`

	BatBuoc bool `json:"bat_buoc"`

	ChonNhieu bool `json:"chon_nhieu"`

	SoLuongToiDa int `json:"so_luong_toi_da"`

	SoLuongToiThieu int `json:"so_luong_toi_thieu"`

	OptionItems []CreateOptionItemInput `json:"option_items"`
}

type CreateOptionItemInput struct {
	TenOption string `json:"ten_option"`

	GiaThem float64 `json:"gia_them"`
}
