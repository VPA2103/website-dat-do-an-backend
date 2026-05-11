package models

type NguoiDung struct {
	MaNguoiDung uint `gorm:"primaryKey;autoIncrement" json:"ma_nguoi_dung"`

	HoTen         string `json:"ho_ten" form:"ho_ten"`
	GioiTinh      string `json:"gioi_tinh" form:"gioi_tinh"`
	NgaySinh      string `json:"ngay_sinh" form:"ngay_sinh"`
	SDT           string `json:"sdt" form:"sdt"`
	NgayVaoLam    string `json:"ngay_vao_lam" form:"ngay_vao_lam"`
	Email         string `json:"email" form:"email"`
	MatKhau       string `json:"mat_khau" form:"mat_khau"`
	LoaiNguoiDung string `gorm:"type:text;not null" json:"loai_nguoi_dung" form:"loai_nguoi_dung"`

	// Quan hệ
	DiaChis []DiaChi `gorm:"foreignKey:MaNguoiDung;references:MaNguoiDung" json:"dia_chis,omitempty"`

	DatBans     []DatBan  `gorm:"foreignKey:IDNhanVienXacNhan" json:"dat_bans,omitempty"`
	AnhNhanVien []HinhAnh `gorm:"polymorphic:Owner;polymorphicValue:nguoi_dung" json:"anh_nguoi_dung,omitempty"`

	YeuThichs []YeuThich `gorm:"foreignKey:MaNguoiDung;references:MaNguoiDung" json:"yeu_thichs,omitempty"`
	DanhGias  []DanhGia  `gorm:"foreignKey:MaNguoiDung;references:MaNguoiDung" json:"danh_gias,omitempty"`
	BinhLuans []BinhLuan `gorm:"foreignKey:MaNguoiDung;references:MaNguoiDung" json:"binh_luans,omitempty"`
}
