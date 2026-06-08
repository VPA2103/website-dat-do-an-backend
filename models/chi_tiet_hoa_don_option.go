package models

type ChiTietHoaDonOption struct {
	MaCTHDOption uint `gorm:"primaryKey;autoIncrement" json:"ma_cthd_option"`

	MaChiTiet uint          `json:"ma_chi_tiet"`
	ChiTiet   ChiTietHoaDon `gorm:"foreignKey:MaChiTiet;" json:"-"`

	MaOptionItem uint       `json:"ma_option_item"`
	OptionItem   OptionItem `gorm:"foreignKey:MaOptionItem;references:MaOptionItem" json:"option_item"`

	TenOption string  `json:"ten_option"`
	GiaThem  float64 `json:"gia_them"`
}