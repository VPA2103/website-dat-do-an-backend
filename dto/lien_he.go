package dto

import "time"

type LienHe struct {
	MaLienHe  uint      `json:"id"`
	SDT       string    `json:"sdt" form:"sdt"`
	Email     string    `json:"email" form:"email"`
	TieuDe    string    `json:"tieu_de" form:"tieu_de"`
	NoiDung   string    `json:"noi_dung" form:"noi_dung"`
	HoTen     string    `json:"ho_ten" form:"ho_ten"`
	TrangThai string    `json:"trang_thai" form:"trang_thai"`
	NgayTao   time.Time `json:"ngay_tao"`
}