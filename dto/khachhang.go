package dto

type KhachHang struct {
	MaKH         uint   `gorm:"primaryKey;autoIncrement" json:"ma_kh"`
	HoTen        string `json:"ho_ten"`
	GioiTinh     string `json:"gioi_tinh"`
	NgaySinh     string `json:"ngay_sinh"`
	DiaChi       string `json:"dia_chi"`
	Email        string `json:"email"`
	MatKhau      string `json:"-"`
	AnhKhachHang string `json:"anh_khach_hang"`
	SDT          string `json:"sdt"`
}
