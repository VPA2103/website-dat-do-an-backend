package models

import "time"

type BinhLuan struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	MaNguoiDung uint `gorm:"index" json:"ma_nguoi_dung"`
	MaMonAn     uint `gorm:"index" json:"ma_mon_an"`

	NoiDung string `json:"noi_dung"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	NguoiDung NguoiDung `gorm:"foreignKey:MaNguoiDung;references:MaNguoiDung" json:"nguoi_dung"`
}
