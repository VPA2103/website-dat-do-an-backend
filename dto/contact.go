package dto

type Contact struct {
	MaLienHe     string `gorm:"primaryKey;size:10;autoIncrement"`
	TenKhachHang string
	SDT          string
	Email        string
	NoiDung      string
}
