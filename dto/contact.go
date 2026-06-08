package dto

type Contact struct {
	MaLienHe     string `json:"ma_lien_he"`
	TenKhachHang string `json:"ten_khach_hang"`
	SDT          string `json:"sdt"`
	Email        string `json:"email"`
	NoiDung      string `json:"noi_dung"`
}