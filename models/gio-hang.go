package models

type GioHang struct {
	MaGioHang   uint `gorm:"primaryKey;autoIncrement"`
	MaNguoiDung uint `gorm:"index"`
	MaMonAn     uint `gorm:"index"`

	MonAn MonAn `gorm:"foreignKey:MaMonAn;references:MaMonAn"`

	SoLuong int
	GiaTien int

	Options []GioHangOption `gorm:"foreignKey:MaGioHang;constraint:OnDelete:CASCADE"`
}

type GioHangOption struct {
	MaGioHangOption uint `gorm:"primaryKey;autoIncrement"`

	MaGioHang    uint `json:"ma_gio_hang"`
	MaNhomOption uint `json:"ma_nhom_option"`
	MaOptionItem uint `json:"ma_option_item"`

	TenNhomOption string     `json:"ten_nhom_option"`
	TenOption     string     `json:"ten_option"`
	GiaThem       int        `json:"gia_them"`
	OptionItem    OptionItem `gorm:"foreignKey:MaOptionItem;references:MaOptionItem"`
}
