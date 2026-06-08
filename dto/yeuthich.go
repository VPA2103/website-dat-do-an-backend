package dto

import "time"

type YeuThich struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	MaNguoiDung uint  `gorm:"uniqueIndex:idx_user_mon"`
	MaMonAn     uint  `gorm:"uniqueIndex:idx_user_mon"`
	MonAn       MonAn `gorm:"foreignKey:MaMonAn;references:MaMonAn" json:"mon_an"`
	CreatedAt   time.Time
}
