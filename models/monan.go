package models

type MonAn struct {
	MaMonAn     uint    `gorm:"primaryKey;autoIncrement" form:"ma_mon_an" json:"ma_mon_an"`
	MaLoaiMonAn uint    `form:"ma_loai_mon_an" json:"ma_loai_mon_an"`
	TenMonAn    string  `form:"ten_mon_an" json:"ten_mon_an"`
	GiaTien     float64 `form:"gia_tien" json:"gia_tien"`
	TrangThai   uint    `form:"trang_thai" json:"trang_thai"`
	MoTa        string  `form:"mo_ta" json:"mo_ta"`

	AnhMonAn    []HinhAnh    `gorm:"polymorphic:Owner;polymorphicValue:mon_an" json:"anh_mon_an,omitempty"`
	NhomOptions []NhomOption `gorm:"foreignKey:MaMonAn" json:"nhom_options"`
}
