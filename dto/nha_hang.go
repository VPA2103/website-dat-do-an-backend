package dto

type NhaHang struct {
	MaNhaHang   uint      `gorm:"primaryKey;size:10;autoIncrement" json:"ma_nha_hang" form:"ma_nha_hang"`
	TenNhaHang  string    `json:"ten_nha_hang" form:"ten_nha_hang"`
	MaNguoiDung uint      `json:"ma_nguoi_dung" form:"ma_nguoi_dung"`
	TrangThai   int       `json:"trang_thai" form:"trang_thai"`
	DiaChi   string       `json:"dia_chi" form:"dia_chi"`
	SoTaiKhoan   int       `json:"so_tai_khoan" form:"so_tai_khoan"`
	NganHang   string       `json:"ngan_hang" form:"ngan_hang"`
	TenNguoiNhan   string       `json:"ten_nguoi_nhan" form:"ten_nguoi_nhan"`
	// Ảnh nhà hàng (polymorphic giống BanAn)
	AnhNhaHang []HinhAnh `gorm:"polymorphic:Owner;polymorphicValue:nha_hang" json:"anh_nha_hang,omitempty"`
}