package dto

import "time"

type LienHe struct {
	MaLienHe  uint      `json:"id" gorm:"primaryKey"`
	SDT       string    `json:"sdt" form:"sdt" gorm:"unique"`
	Email     string    `json:"email" form:"email" gorm:"type:varchar(150);not null;index"`
	TieuDe    string    `json:"tieu_de" form:"tieu_de" gorm:"type:varchar(200);not null"`
	NoiDung   string    `json:"noi_dung" form:"noi_dung" gorm:"type:text;not null"`
	HoTen     string    `json:"ho_ten" form:"ho_ten" gorm:"type:varchar(100);not null"`
	TrangThai string    `json:"trang_thai" form:"trang_thai" gorm:"type:varchar(50);default:'chua_xu_ly'"`
	NgayTao   time.Time `json:"ngay_tao" gorm:"autoCreateTime"`
}
