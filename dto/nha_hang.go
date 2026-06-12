package dto

type NhaHang struct {
	MaNhaHang    uint      `json:"ma_nha_hang" form:"ma_nha_hang"`
	TenNhaHang   string    `json:"ten_nha_hang" form:"ten_nha_hang"`
	TrangThai     int       `json:"trang_thai" form:"trang_thai"`
	DiaChi        string    `json:"dia_chi" form:"dia_chi"`
	SoTaiKhoan    string       `json:"so_tai_khoan" form:"so_tai_khoan"`
	NganHang      string    `json:"ngan_hang" form:"ngan_hang"`
	TenNguoiNhan  string    `json:"ten_nguoi_nhan" form:"ten_nguoi_nhan"`

	// 👉 Mới thêm
	GioMoCua     string    `json:"gio_mo_cua" form:"gio_mo_cua"`
	GioDongCua   string    `json:"gio_dong_cua" form:"gio_dong_cua"`
	MoTa         string    `json:"mo_ta" form:"mo_ta"`

	// Ảnh nhà hàng
	AnhNhaHang   []HinhAnh `json:"anh_nha_hang,omitempty"`
}