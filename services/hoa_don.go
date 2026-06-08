package services

import (
	"fmt"
	"os"

	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"github.com/vpa/quanlynhahang-backend/utils"
)

func GetHoaDonByBan(maBan uint) (*models.HoaDon, error) {

	var hoaDon models.HoaDon

	err := config.DB.
		Where("ma_ban = ? AND trang_thai < ?", maBan, 2).
		First(&hoaDon).Error

	if err != nil {
		return nil, err
	}

	return &hoaDon, nil
}

func GetMonDaGoi(maBan uint) ([]models.ChiTietHoaDon, error) {

	var hoaDon models.HoaDon

	err := config.DB.
		Where("ma_ban = ? AND trang_thai = 0", maBan).
		First(&hoaDon).Error

	if err != nil {
		return nil, err
	}

	var ds []models.ChiTietHoaDon

	err = config.DB.
		Preload("MonAn").
		Where("ma_hd = ?", hoaDon.MaHoaDon).
		Find(&ds).Error

	return ds, err
}

func TaoQRThanhToan(maBan uint) (string, float64, error) {

	var hoaDon models.HoaDon

	err := config.DB.
		Where("ma_ban = ? AND trang_thai = 0", maBan).
		First(&hoaDon).Error

	if err != nil {
		return "", 0, err
	}

	qr, err := utils.GenerateQRPayment(
		hoaDon.TongTien,
		os.Getenv("VIETQR_BANK_BIN"),
		os.Getenv("VIETQR_ACCOUNT_NO"),
		fmt.Sprintf("HOADON_%d", hoaDon.MaHoaDon),
	)

	if err != nil {
		return "", 0, err
	}

	// chuyển trạng thái chờ thanh toán
	config.DB.Model(&hoaDon).
		Update("trang_thai", 1)

	return qr, hoaDon.TongTien, nil
}
