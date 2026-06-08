package dto

type BanAn struct {
	MaBan     uint   `gorm:"primaryKey;autoIncrement"`
	TenBan    string
	SoChoNgoi int
	TrangThai int

	AnhBan []HinhAnh `gorm:"polymorphic:Owner;polymorphicValue:ban_an"`
}
