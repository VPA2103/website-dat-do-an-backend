package models

type DatBan struct {
	MaDatBan           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	SDT          string `json:"sdt" form:"sdt" binding:"required"`
	TenKhachHang string `json:"ten_khach_hang" form:"ten_khach_hang" binding:"required"`
	Email        string `json:"email" form:"email" binding:"required"`
	GhiChu       string `json:"ghi_chu" form:"ghi_chu"`
	MaBanAn      uint   `json:"ma_ban_an" form:"ma_ban_an" binding:"required"`
	Ngay         string `json:"ngay" form:"ngay" binding:"required"`
	Gio          string `json:"gio" form:"gio" binding:"required"`

	TrangThai string `gorm:"type:varchar(50);default:'dang_xu_ly'" json:"trang_thai"`

	// 🔑 Khóa ngoại – cho phép NULL
	IDNhanVienXacNhan *uint     `json:"id_nhan_vien_xac_nhan"`
	NhanVienXacNhan   *NguoiDung `gorm:"foreignKey:IDNhanVienXacNhan;references:MaNguoiDung" json:"nhan_vien_xac_nhan,omitempty"`
}
