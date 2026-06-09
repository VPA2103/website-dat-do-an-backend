package dto

import "time"

type HoaDon struct {
	MaHoaDon        uint      `json:"ma_hoa_don"`
	MaNguoiDung     uint      `json:"ma_nguoi_dung"`
	HoTen           string    `json:"ho_ten"`
	SDT             string    `json:"sdt"`
	DiaChi          string    `json:"dia_chi"`
	GhiChu          string    `json:"ghi_chu"`
	Ngay            time.Time `json:"ngay"`

	TongTien           float64 `json:"tong_tien"`
	TamTinh            float64 `json:"tam_tinh"`
	TienGiam           float64 `json:"tien_giam"`
	TrangThai          string  `json:"trang_thai"`
	TrangThaiThanhToan string  `json:"trang_thai_thanh_toan"`

	GiamGiaID      *uint           `json:"giam_gia_id"`
	GiamGia        GiamGia         `json:"giam_gia"`
	ChiTietHoaDons []ChiTietHoaDon `json:"chi_tiet_hoa_dons"`
	ThanhToans     *ThanhToan      `json:"thanh_toans"`
}