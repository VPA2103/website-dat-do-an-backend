package dto

type NguoiDung struct {
	MaNguoiDung   uint   `json:"ma_nguoi_dung"`
	HoTen         string `json:"ho_ten" form:"ho_ten"`
	GioiTinh      string `json:"gioi_tinh" form:"gioi_tinh"`
	NgaySinh      string `json:"ngay_sinh" form:"ngay_sinh"`
	SDT           string `json:"sdt" form:"sdt"`
	NgayVaoLam    string `json:"ngay_vao_lam" form:"ngay_vao_lam"`
	Email         string `json:"email" form:"email"`
	MatKhau       string `json:"mat_khau" form:"mat_khau"`
	LoaiNguoiDung string `json:"loai_nguoi_dung" form:"loai_nguoi_dung"`
	TrangThai     string `json:"trang_thai" form:"trang_thai"`

	// Quan hệ
	DiaChis     []DiaChi  `json:"dia_chis,omitempty"`
	DatBans     []DatBan  `json:"dat_bans,omitempty"`
	AnhNhanVien []HinhAnh `json:"anh_nguoi_dung,omitempty"`
	YeuThichs   []YeuThich `json:"yeu_thichs,omitempty"`
	DanhGias    []DanhGia  `json:"danh_gias,omitempty"`
	BinhLuans   []BinhLuan `json:"binh_luans,omitempty"`
}