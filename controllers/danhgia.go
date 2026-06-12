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

type ThongKeDanhGiaNgaymodels struct {
	Ngay      string `json:"ngay"`
	SoDanhGia int64  `json:"so_danh_gia"`
}

func NewDanhGiaController(hub *websocket.Hub) *DanhGiaController {
	return &DanhGiaController{Hub: hub}
}

type DanhGiaTheoMonResponse struct {
	MonAn    models.MonAn     `json:"mon_an"`
	DanhGias []models.DanhGia `json:"danh_gias"`
}

type DanhGiaInput struct {
	MaHoaDon    uint   `json:"ma_hoa_don"`
	MaNguoiDung uint   `json:"ma_nguoi_dung"`
	MaMonAn     uint   `json:"ma_mon_an"`
	SoSao       int    `json:"so_sao"`
	NoiDung     string `json:"noi_dung"`
	TrangThai string `json:"trang_thai"`
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
	TrangThai string `json:"trang_thai"`
}

func (ctrl *DanhGiaController) CreateDanhGia(c *gin.Context) {
	var input DanhGiaInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 🔥 CHECK TRÙNG ĐÁNH GIÁ
	var count int64
	config.DB.Model(&models.DanhGia{}).
		Where("ma_hoa_don = ? AND ma_nguoi_dung = ? AND ma_mon_an = ?",
			input.MaHoaDon, input.MaNguoiDung, input.MaMonAn).
		Count(&count)

	if count > 0 {
		c.JSON(400, gin.H{
			"error": "Bạn đã đánh giá món này rồi",
		})
		return
	}

	dg := models.DanhGia{
		MaHoaDon:    input.MaHoaDon,
		MaNguoiDung: input.MaNguoiDung,
		MaMonAn:     input.MaMonAn,
		SoSao:       input.SoSao,
		NoiDung:     input.NoiDung,
		TrangThai:   "hien",
	}

	config.DB.Create(&dg)

	c.JSON(200, dg)
}

func (ctrl *DanhGiaController) AnDanhGia(c *gin.Context) {
	id := c.Param("id")

	var danhGia models.DanhGia

	// 🔍 kiểm tra tồn tại
	if err := config.DB.First(&danhGia, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Không tìm thấy đánh giá",
		})
		return
	}

	// 🔒 cập nhật trạng thái = ẩn
	if err := config.DB.Model(&danhGia).
		Update("trang_thai", "an").Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Không thể ẩn đánh giá",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Ẩn đánh giá thành công",
	})
}
func (ctrl *DanhGiaController) HienDanhGia(c *gin.Context) {
	id := c.Param("id")

	var danhGia models.DanhGia

	// 🔍 kiểm tra tồn tại
	if err := config.DB.First(&danhGia, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Không tìm thấy đánh giá",
		})
		return
	}

	// 🔒 cập nhật trạng thái = ẩn
	if err := config.DB.Model(&danhGia).
		Update("trang_thai", "hien").Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Không thể hiện đánh giá",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Hiện đánh giá thành công",
	})
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

//cua giao dien menu
func GetDanhGiaByMonAn(c *gin.Context) {
	maMon := c.Param("id")

	var data []models.DanhGia

	config.DB.
		Preload("NguoiDung").
		Preload("NguoiDung.AnhNhanVien").
		Where("ma_mon_an = ? AND trang_thai = ?", maMon, "hien").
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
	maHD, err1 := strconv.Atoi(c.Query("ma_hoa_don"))
	maUser, err2 := strconv.Atoi(c.Query("ma_nguoi_dung"))
	maMon, err3 := strconv.Atoi(c.Query("ma_mon_an"))

	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(400, gin.H{"error": "invalid params"})
		return
	}

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
		"data": ThongKeDanhGiaNgaymodels{
			Ngay:      ngay,
			SoDanhGia: soDanhGia,
		},
	})
}
func (ctrl *DanhGiaController) GetAllDanhGiaByNguoiDung(c *gin.Context) {

	// 1️⃣ Lấy user_id từ middleware
	maNguoiDungAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Chưa đăng nhập"})
		return
	}
	maNguoiDung := maNguoiDungAny.(uint)

	// 2️⃣ Lấy toàn bộ đánh giá của user
	var danhGias []models.DanhGia
	err := config.DB.
		Where("ma_nguoi_dung = ?", maNguoiDung).
		Preload("NguoiDung").
		Preload("NguoiDung.AnhNhanVien").
		Preload("MonAn").
		Preload("MonAn.AnhMonAn").
		Find(&danhGias).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Không thể lấy đánh giá"})
		return
	}

	// 3️⃣ Gom theo món ăn
	mapTheoMon := make(map[uint][]models.DanhGia)

	for _, dg := range danhGias {
		mapTheoMon[dg.MaMonAn] = append(mapTheoMon[dg.MaMonAn], dg)
	}

	// 4️⃣ Build response
	var result []DanhGiaTheoMonResponse

	for _, listDG := range mapTheoMon {
		if len(listDG) == 0 {
			continue
		}

		result = append(result, DanhGiaTheoMonResponse{
			MonAn:    listDG[0].MonAn, // cùng món
			DanhGias: listDG,
		})
	}

	// 5️⃣ Response
	c.JSON(200, gin.H{
		"data": result,
	})
}
func GetAllDanhGiaByMonAn(c *gin.Context) {
	maMon := c.Param("id")

	var data []models.DanhGia

	config.DB.
		Preload("NguoiDung").
		Preload("NguoiDung.AnhNhanVien").
		Where("ma_mon_an = ? ", maMon).
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
			TrangThai: d.TrangThai,
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
