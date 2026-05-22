package models

type NhomOption struct {
	MaNhomOption uint `gorm:"primaryKey;autoIncrement" json:"ma_nhom_option"`

	MaMonAn uint  `json:"ma_mon_an"`
	MonAn   MonAn `gorm:"foreignKey:MaMonAn;references:MaMonAn" json:"-"`

	TenNhom string `json:"ten_nhom"`

	BatBuoc bool `json:"bat_buoc"` // bắt buộc chọn

	ChonNhieu bool `json:"chon_nhieu"` // cho chọn nhiều

	SoLuongToiDa    int `json:"so_luong_toi_da"`
	SoLuongToiThieu int `json:"so_luong_toi_thieu"`

	TrangThai uint `json:"trang_thai"`

	OptionItems []OptionItem `gorm:"foreignKey:MaNhomOption;constraint:OnDelete:CASCADE" json:"OptionItems"`
}

type OptionItem struct {
	MaOptionItem uint `gorm:"primaryKey;autoIncrement" json:"ma_option_item"`

	MaNhomOption uint       `json:"ma_nhom_option"`
	NhomOption   NhomOption `gorm:"foreignKey:MaNhomOption;references:MaNhomOption"`

	TenOption string `json:"ten_option"`

	GiaThem float64 `json:"gia_them"`

	TrangThai uint `json:"trang_thai"`
}

