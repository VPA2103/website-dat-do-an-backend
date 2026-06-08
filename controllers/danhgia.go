package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/models"
)

type DanhGiaController struct {
	Hub *websocket.Hub
}

type ThongKeDanhGiaNgayDTO struct {
	Ngay        string `json:"ngay"`
	SoDanhGia   int64  `json:"so_danh_gia"`
}

func NewDanhGiaController(hub *websocket.Hub) *DanhGiaController {
	return &DanhGiaController{Hub: hub}
}

type DanhGiaInput struct {
	MaHoaDon    uint   `json:"ma_hoa_don"`
	MaNguoiDung uint   `json:"ma_nguoi_dung"`
	MaMonAn     uint   `json:"ma_mon_an"`
	SoSao       int    `json:"so_sao"`
	NoiDung     string `json:"noi_dung"`
}
type NguoiDungMini struct {
	MaNguoiDung uint   `json:"ma_nguoi_dung"`
	HoTen       string `json:"ho_ten"`
	Anh         string `json:"anh"`
}
type DanhGiaResponse struct {
	ID        uint          `json:"id"`
	MaMonAn   uint          `json:"ma_mon_an"`
	SoSao     int           `json:"so_sao"`
	NoiDung   string        `json:"noi_dung"`
	NguoiDung NguoiDungMini `json:"nguoi_dung"`
}

func (ctrl *DanhGiaController) CreateDanhGia(c *gin.Context) {
	var input DanhGiaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	dg := models.DanhGia{
		MaHoaDon:    input.MaHoaDon,
		MaNguoiDung: input.MaNguoiDung,
		MaMonAn:     input.MaMonAn,
		SoSao:       input.SoSao,
		NoiDung:     input.NoiDung,
	}

	config.DB.Create(&dg)

	c.JSON(200, dg)
}

func GetRatingByMon(c *gin.Context) {
	rows, err := config.DB.Raw(`
		SELECT 
			ma_mon_an,
			AVG(CAST(so_sao AS FLOAT)) AS avg_sao,
			COUNT(*) AS tong_danh_gia
		FROM danh_gia
		GROUP BY ma_mon_an
	`).Rows()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type Result struct {
		MaMonAn     uint    `json:"ma_mon_an"`
		AvgSao      float64 `json:"avg_sao"`
		TongDanhGia int     `json:"tong_danh_gia"`
	}

	var result []Result

	for rows.Next() {
		var item Result
		rows.Scan(&item.MaMonAn, &item.AvgSao, &item.TongDanhGia)
		result = append(result, item)
	}

	c.JSON(200, result)
}

func GetDanhGiaByMonAn(c *gin.Context) {
	maMon := c.Param("id")

	var data []models.DanhGia

	config.DB.
		Preload("NguoiDung").
		Preload("NguoiDung.AnhNhanVien").
		Where("ma_mon_an = ?", maMon).
		Find(&data)

	var res []DanhGiaResponse

	for _, d := range data {

		anh := ""
		if len(d.NguoiDung.AnhNhanVien) > 0 {
			anh = d.NguoiDung.AnhNhanVien[0].Url
		}

		res = append(res, DanhGiaResponse{
			ID:      d.ID,
			MaMonAn: d.MaMonAn,
			SoSao:   d.SoSao,
			NoiDung: d.NoiDung,
			NguoiDung: NguoiDungMini{
				MaNguoiDung: d.NguoiDung.MaNguoiDung,
				HoTen:       d.NguoiDung.HoTen,
				Anh:         anh,
			},
		})
	}

	c.JSON(200, gin.H{
		"data": res,
	})
}

func (ctrl *DanhGiaController) UpdateDanhGia(c *gin.Context) {
	id := c.Param("id")

	var dg models.DanhGia
	if err := config.DB.First(&dg, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy"})
		return
	}

	var input DanhGiaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	dg.SoSao = input.SoSao
	dg.NoiDung = input.NoiDung

	if err := config.DB.Save(&dg).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dg)
}

func (ctrl *DanhGiaController) DeleteDanhGia(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.DanhGia{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Đã xoá đánh giá"})
}

func CheckDanhGia(c *gin.Context) {
	maHD, _ := strconv.Atoi(c.Query("ma_hd"))
	maUser, _ := strconv.Atoi(c.Query("ma_nguoi_dung"))
	maMon, _ := strconv.Atoi(c.Query("ma_mon_an"))

	var count int64

	config.DB.Model(&models.DanhGia{}).
		Where("ma_hoa_don = ? AND ma_nguoi_dung = ? AND ma_mon_an = ?",
			maHD, maUser, maMon).
		Count(&count)

	c.JSON(200, gin.H{
		"da_danh_gia": count > 0,
	})
}

func (ctrl *DanhGiaController) GetSoLuongDanhGiaHomNay(c *gin.Context) {

	ngay := time.Now().Format("2006-01-02")

	var soDanhGia int64

	err := config.DB.
		Model(&models.DanhGia{}).
		Where(`
			CAST(ngay_danh_gia AS DATE) = ?
		`, ngay).
		Count(&soDanhGia).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Không thể thống kê đánh giá hôm nay"})
		return
	}

	c.JSON(200, gin.H{
		"data": ThongKeDanhGiaNgayDTO{
			Ngay:      ngay,
			SoDanhGia: soDanhGia,
		},
	})
}
