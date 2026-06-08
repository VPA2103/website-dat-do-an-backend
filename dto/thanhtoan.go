package dto

import "time"

type ThanhToan struct {
	MaThanhToan       uint      `json:"ma_thanh_toan"`
	MaHoaDon          uint      `json:"ma_hd"`
	HoaDon            HoaDon    `json:"hoa_don"`
	SoTien            float64   `json:"so_tien"`
	HinhThucThanhToan string    `json:"hinh_thuc_thanh_toan"`
	NgayThanhToan     time.Time `json:"ngay_thanh_toan"`
	MaNVThanhToan     string    `json:"ma_nv_thanh_toan"`
	NhanVienThanhToan NguoiDung `json:"nhan_vien_thanh_toan"`
}