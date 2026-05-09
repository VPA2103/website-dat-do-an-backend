package services

import (
	"time"

	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"gorm.io/gorm"
)

func GetOrCreateHoaDon(maBan uint) (*models.HoaDon, error) {

	var hoaDon models.HoaDon

	err := config.DB.Where("ma_ban = ? AND trang_thai = ?", maBan, 0).
		First(&hoaDon).Error

	if err != nil {

		if err == gorm.ErrRecordNotFound {

			hoaDon = models.HoaDon{
				MaBan:     maBan,
				Ngay:      time.Now(),
				TrangThai: "pending",
			}

			if err := config.DB.Create(&hoaDon).Error; err != nil {
				return nil, err
			}

		} else {
			return nil, err
		}
	}

	return &hoaDon, nil
}
func CloseHoaDon(maBan uint) error {

	var hoaDon models.HoaDon

	// tìm hóa đơn đang mở
	if err := config.DB.
		Where("ma_ban = ? AND trang_thai = ?", maBan, 0).
		First(&hoaDon).Error; err != nil {
		return err
	}

	// cập nhật trạng thái
	if err := config.DB.Model(&hoaDon).
		Update("trang_thai", 1).Error; err != nil {
		return err
	}

	return nil
}

func UpdateTongTien(maHD uint) error {

	var tong float64

	if err := config.DB.Model(&models.ChiTietHoaDon{}).
		Where("ma_hd = ?", maHD).
		Select("SUM(thanh_tien)").
		Scan(&tong).Error; err != nil {
		return err
	}

	if err := config.DB.Model(&models.HoaDon{}).
		Where("ma_hd = ?", maHD).
		Update("tong_tien", tong).Error; err != nil {
		return err
	}

	return nil
}

func AddMon(maBan uint, maMon uint, soLuong int) error {

	if soLuong <= 0 {
		return nil
	}

	// lấy hoặc tạo hóa đơn
	hoaDon, err := GetOrCreateHoaDon(maBan)
	if err != nil {
		return err
	}

	var chiTiet models.ChiTietHoaDon

	err = config.DB.
		Where("ma_hd = ? AND ma_mon_an = ?", hoaDon.MaHD, maMon).
		First(&chiTiet).Error

	// nếu món đã tồn tại trong hóa đơn
	if err == nil {

		chiTiet.SoLuong += soLuong
		chiTiet.ThanhTien = float64(chiTiet.SoLuong) * chiTiet.DonGia

		if err := config.DB.Save(&chiTiet).Error; err != nil {
			return err
		}

		// nếu chưa có món trong hóa đơn
	} else if err == gorm.ErrRecordNotFound {

		var mon models.MonAn

		if err := config.DB.
			First(&mon, "ma_mon_an = ?", maMon).Error; err != nil {
			return err
		}

		chiTiet = models.ChiTietHoaDon{
			MaHoaDon:  hoaDon.MaHD,
			MaMonAn:   maMon,
			SoLuong:   soLuong,
			DonGia:    mon.GiaTien,
			ThanhTien: float64(soLuong) * mon.GiaTien,
			TrangThai: "dang_goi",
		}

		if err := config.DB.Create(&chiTiet).Error; err != nil {
			return err
		}

	} else {
		return err
	}

	// cập nhật tổng tiền hóa đơn
	if err := UpdateTongTien(hoaDon.MaHD); err != nil {
		return err
	}

	return nil
}
