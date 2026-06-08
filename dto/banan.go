package dto

type BanAn struct {
	MaBan     uint     `gorm:"primaryKey;size:10;autoIncrement" json:"ma_ban" form:"ma_ban"`
	TenBan    string   `json:"ten_ban" form:"ten_ban"`
	SoChoNgoi int      `json:"so_cho_ngoi" form:"so_cho_ngoi"`
	TrangThai int      `json:"trang_thai" form:"trang_thai"`
	AnhBan    []HinhAnh `gorm:"polymorphic:Owner;polymorphicValue:ban_an" json:"anh_ban,omitempty"`
	// Anh_QR    string   `json:"anh_qr" form:"anh_qr"`
}
