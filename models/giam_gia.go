package models

import "time"

type GiamGia struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	Code string `gorm:"size:50;unique;not null" json:"code"`

	TenChuongTrinh string `gorm:"size:255" json:"ten_chuong_trinh"`

	LoaiGiamGia string `gorm:"type:varchar(20);not null" json:"loai_giam_gia"`

	GiaTriGiam float64 `gorm:"type:decimal(10,2);not null" json:"gia_tri_giam"`

	DonToiThieu float64 `gorm:"type:decimal(10,2);default:0" json:"don_toi_thieu"`

	GiamToiDa float64 `gorm:"type:decimal(10,2)" json:"giam_toi_da"`

	GioiHanSuDung *int `json:"gioi_han_su_dung"`

	SoLanDaDung int `gorm:"default:0" json:"so_lan_da_dung"`

	NgayBatDau  time.Time `json:"ngay_bat_dau"`
	NgayKetThuc time.Time `json:"ngay_ket_thuc"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	AnhGiamGia []HinhAnh `gorm:"polymorphic:Owner;polymorphicValue:giam_gia;constraint:OnDelete:CASCADE" json:"anh_giam_gia,omitempty"`
}
