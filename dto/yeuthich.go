package dto

import "time"

type YeuThich struct {
	ID          uint      `json:"id"`
	MaNguoiDung uint      `json:"ma_nguoi_dung"`
	MaMonAn     uint      `json:"ma_mon_an"`
	MonAn       MonAn     `json:"mon_an"`
	CreatedAt   time.Time `json:"created_at"`
}