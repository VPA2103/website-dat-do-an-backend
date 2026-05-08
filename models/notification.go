package models

import "time"

type ThongBao struct {
	ID        uint   `gorm:"primaryKey"`
	MaNguoiDung    uint   `gorm:"index"`
	Type      string `gorm:"size:50"` // order, system, chat,...
	Title     string
	Content   string
	TrangThai    bool `gorm:"index"`
	NgayTao time.Time
}
