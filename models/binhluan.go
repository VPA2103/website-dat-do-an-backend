package models

import "time"

type BinhLuan struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	MaNguoiDung uint `gorm:"index" json:"ma_nguoi_dung"`
	MaMonAn     uint `gorm:"index" json:"ma_mon_an"`

	// ✅ SỬA Ở ĐÂY
	MonAn MonAn `gorm:"foreignKey:MaMonAn;references:MaMonAn" json:"mon_an"`

	// 👇 Bình luận cha
	ParentID *uint `gorm:"index" json:"parent_id"`

	NoiDung string `json:"noi_dung"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// ✅ SỬA Ở ĐÂY
	NguoiDung NguoiDung `gorm:"foreignKey:MaNguoiDung;references:MaNguoiDung" json:"nguoi_dung"`

	// 👇 Reply con
	Replies []BinhLuan `gorm:"foreignKey:ParentID" json:"replies"`
}
