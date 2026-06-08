package dto

import "time"

type GiamGia struct {
	ID              uint      `json:"id"`
	Code            string    `json:"code"`
	TenChuongTrinh  string    `json:"ten_chuong_trinh"`
	LoaiGiamGia     string    `json:"loai_giam_gia"`
	GiaTriGiam      float64   `json:"gia_tri_giam"`
	DonToiThieu     float64   `json:"don_toi_thieu"`
	GiamToiDa       float64   `json:"giam_toi_da"`
	GioiHanSuDung   *int      `json:"gioi_han_su_dung"`
	SoLanDaDung     int       `json:"so_lan_da_dung"`
	NgayBatDau      time.Time `json:"ngay_bat_dau"`
	NgayKetThuc     time.Time `json:"ngay_ket_thuc"`
	IsActive        bool      `json:"is_active"`
	AnhGiamGia      []HinhAnh `json:"anh_giam_gia,omitempty"`
}