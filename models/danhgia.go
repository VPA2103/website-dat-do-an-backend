package models

import "time"

type DanhGia struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	MaNguoiDung uint `gorm:"uniqueIndex:idx_user_mon"`
	MaMonAn     uint `gorm:"uniqueIndex:idx_user_mon"`

	SoSao   int `gorm:"check:so_sao >= 1 AND so_sao <= 5"`
	NoiDung string

	CreatedAt time.Time
	UpdatedAt time.Time
}
