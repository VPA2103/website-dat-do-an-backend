package models

import "time"

type YeuThich struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	MaNguoiDung uint  `gorm:"uniqueIndex:idx_user_mon"`
	MaMonAn     uint  `gorm:"uniqueIndex:idx_user_mon"`
	MonAn       MonAn `gorm:"foreignKey:MaMonAn;" json:"mon_an"`
	CreatedAt   time.Time
}
