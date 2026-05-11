package models

import "time"

type BinhLuan struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	MaNguoiDung uint `gorm:"index" json:"ma_nguoi_dung"`
	MaMonAn     uint `gorm:"index" json:"ma_mon_an"`

	NoiDung string `json:"noi_dung"`
	MaCha   *uint  `gorm:"index" json:"ma_cha"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	BinhLuanCha *BinhLuan  `gorm:"foreignKey:MaCha" json:"binh_luan_cha,omitempty"`
	BinhLuans   []BinhLuan `gorm:"foreignKey:MaCha" json:"binh_luans,omitempty"`

	NguoiDung NguoiDung `gorm:"foreignKey:MaNguoiDung;" json:"nguoi_dung"`
}
