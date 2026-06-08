package dto

import "time"

type DanhGia struct {
	ID uint `json:"id"`

	MaHoaDon    uint `json:"ma_hoa_don"`
	MaNguoiDung uint `json:"ma_nguoi_dung"`
	MaMonAn     uint `json:"ma_mon_an"`

	SoSao   int    `json:"so_sao"`
	NoiDung string `json:"noi_dung"`

	NguoiDung *NguoiDung `json:"nguoi_dung,omitempty"`

	CreatedAt time.Time `json:"ngay_danh_gia"`
}