package models

type MonAn struct {
	MaMonAn     uint    `gorm:"primaryKey;autoIncrement" json:"ma_mon_an"`
	MaLoaiMonAn uint    `json:"ma_loai_mon_an" form:"ma_loai_mon_an" gorm:"index"`
	TenMonAn    string  `json:"ten_mon_an" form:"ten_mon_an"`
	GiaTien     float64 `json:"gia_tien" form:"gia_tien"`
	TrangThai   uint    `json:"trang_thai" form:"trang_thai"`
	MoTa        string  `json:"mo_ta" form:"mo_ta"`

	AnhMonAn    []HinhAnh    `gorm:"polymorphic:Owner;polymorphicValue:mon_an" json:"anh_mon_an,omitempty"`
	NhomOptions []NhomOption `gorm:"foreignKey:MaMonAn" json:"nhom_options"`
	Document string `gorm:"type:text" json:"document,omitempty"`
	SearchText string `gorm:"type:text;index" json:"search_text,omitempty"`
	Tags string `gorm:"type:text" json:"tags,omitempty"`
	HasEmbedding bool `gorm:"default:false" json:"has_embedding"`
}
