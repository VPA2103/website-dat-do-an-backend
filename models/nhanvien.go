package models

type NguoiDung struct {
	MaNV         uint   `gorm:"primaryKey;autoIncrement" json:"ma_nv"`
	HoTen        string `json:"ho_ten" form:"ho_ten"`
	GioiTinh     string `json:"gioi_tinh" form:"gioi_tinh"`
	NgaySinh     string `json:"ngay_sinh" form:"ngay_sinh"`
	SDT          string `json:"sdt" form:"sdt"`
	DiaChi       string `json:"dia_chi" form:"dia_chi"`
	NgayVaoLam   string `json:"ngay_vao_lam" form:"ngay_vao_lam"`
	Email        string `json:"email" form:"email" `
	MatKhau      string `json:"mat_khau" form:"mat_khau"`
	LoaiNhanVien string `gorm:"type:text;not null" json:"loai_nhan_vien" form:"loai_nhan_vien"`
	// ✅ Khai báo quan hệ 1-1 rõ ràng
	DatBans     []DatBan `gorm:"foreignKey:IDNhanVienXacNhan" json:"dat_bans,omitempty"`
	AnhNhanVien []HinhAnh `gorm:"polymorphic:Owner;polymorphicValue:nhan_vien" json:"anh_nhan_vien,omitempty"`
}
