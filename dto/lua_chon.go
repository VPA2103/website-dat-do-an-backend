package dto

type CreateNhomOptionInput struct {
	MaMonAn         uint                   `json:"ma_mon_an"`
	TenNhom         string                 `json:"ten_nhom"`
	BatBuoc         bool                   `json:"bat_buoc"`
	ChonNhieu       bool                   `json:"chon_nhieu"`
	SoLuongToiDa    int                    `json:"so_luong_toi_da"`
	SoLuongToiThieu int                    `json:"so_luong_toi_thieu"`
	OptionItems     []CreateOptionItemInput `json:"option_items"`
}

type CreateOptionItemInput struct {
	TenOption string  `json:"ten_option"`
	GiaThem   float64 `json:"gia_them"`
}

type NhomOption struct {
	MaNhomOption    uint         `json:"ma_nhom_option"`
	MaMonAn         uint         `json:"ma_mon_an"`
	MonAn           MonAn        `json:"-"`
	TenNhom         string       `json:"ten_nhom"`
	BatBuoc         bool         `json:"bat_buoc"`
	ChonNhieu       bool         `json:"chon_nhieu"`
	SoLuongToiDa    int          `json:"so_luong_toi_da"`
	SoLuongToiThieu int          `json:"so_luong_toi_thieu"`
	TrangThai       uint         `json:"trang_thai"`
	OptionItems     []OptionItem `json:"OptionItems"`
}

type OptionItem struct {
	MaOptionItem uint       `json:"ma_option_item"`
	MaNhomOption uint       `json:"ma_nhom_option"`
	NhomOption   NhomOption `json:"-"`
	TenOption    string     `json:"ten_option"`
	GiaThem      float64    `json:"gia_them"`
	TrangThai    uint       `json:"trang_thai"`
}