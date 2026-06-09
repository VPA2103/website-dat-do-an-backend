package models

import "time"

type HoaDon struct {
	MaHoaDon uint `gorm:"column:ma_hoa_don;primaryKey" json:"ma_hd"`
	MaNguoiDung uint      `json:"ma_nguoi_dung"`
	HoTen       string    `json:"ho_ten"`
	SDT         string    `json:"sdt"`
	DiaChi      string    `json:"dia_chi"`
	GhiChu      string    `json:"ghi_chu"`
	Ngay        time.Time `json:"ngay"`

	TongTien           float64 `json:"tong_tien"` // sau giảm
	TamTinh            float64 `json:"tam_tinh"`  // trước giảm
	TienGiam           float64 `json:"tien_giam"`
	TrangThai          string  `gorm:"type:varchar(30);default:'cho_xac_nhan'" json:"trang_thai"`
	TrangThaiThanhToan string  `json:"trang_thai_thanh_toan"`

	// MaNVOrder      *uint           `gorm:"size:10"`
	// NhanVienOrder  *NhanVien       `gorm:"foreignKey:MaNVOrder;references:MaNV"`
	GiamGiaID      *uint           `json:"giam_gia_id"`
	GiamGia        GiamGia         `gorm:"foreignKey:GiamGiaID;references:ID" json:"giam_gia"`
	ChiTietHoaDons []ChiTietHoaDon `gorm:"foreignKey:MaHoaDon" json:"chi_tiet_hoa_dons"`
	ThanhToans     *ThanhToan      `gorm:"foreignKey:MaHoaDon;references:MaHoaDon" json:"thanh_toans"`
}

type ShipOrder struct {
	ID uint `gorm:"primaryKey"`

	MaHoaDon uint `gorm:"index"`

	MaShipper *uint `gorm:"index"`

	Status string // pending | accepted | rejected | delivering | done

	CreatedAt  time.Time
	AcceptedAt *time.Time
}
