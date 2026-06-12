package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	//"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/xuri/excelize/v2"

	//"github.com/vpa/quanlynhahang-backend/dto"
	//"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/models"
	//"github.com/vpa/quanlynhahang-backend/utils"
	//"gorm.io/gorm"
	//"github.com/xuri/excelize/v2"
)

type MonAnBanChayDTO struct {
	MaMonAn  uint   `json:"ma_mon_an"`
	TenMonAn string `json:"ten_mon_an"`
	SoLuong  int64  `json:"so_luong"`
}

type TopMonBanChayDTO struct {
	MaMonAn  uint   `json:"ma_mon_an"`
	TenMonAn string `json:"ten_mon_an"`
	TongBan  int    `json:"tong_ban"`
}

type DoanhThuDTO struct {
	Ngay              string  `json:"ngay,omitempty"`
	Thang             int     `json:"thang,omitempty"`
	Nam               int     `json:"nam,omitempty"`
	DoanhThu          float64 `json:"doanh_thu"`
	SoDon             int64   `json:"so_don"`
	DoanhThuTrungBinh float64 `json:"doanh_thu_trung_binh"`
}

func GetDoanhThuTheoNgay(c *gin.Context) {

	ngay := c.Query("ngay")
	if ngay == "" {
		ngay = time.Now().Format("2006-01-02")
	}

	start := ngay + " 00:00:00"
	end := ngay + " 23:59:59"

	var result DoanhThuDTO

	err := config.DB.
		Model(&models.HoaDon{}).
		Select(`
			COALESCE(SUM(tong_tien), 0) AS doanh_thu,
			COUNT(*) AS so_don
		`).
		Where(`
			ngay BETWEEN ? AND ?
			AND trang_thai = ?
			AND trang_thai_thanh_toan = ?
		`, start, end, "da_giao", "da_thanh_toan").
		Scan(&result).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể tính doanh thu",
		})
		return
	}

	result.Ngay = ngay

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func GetDoanhThuTheoThang(c *gin.Context) {

	thang, _ := strconv.Atoi(c.DefaultQuery("thang", fmt.Sprint(int(time.Now().Month()))))
	nam, _ := strconv.Atoi(c.DefaultQuery("nam", fmt.Sprint(time.Now().Year())))

	var result DoanhThuDTO

	err := config.DB.
		Model(&models.HoaDon{}).
		Select(`
			COALESCE(SUM(tong_tien), 0) AS doanh_thu,
			COUNT(ma_hoa_don) AS so_don
		`).
		Where(`
			EXTRACT(MONTH FROM ngay) = ?
			AND EXTRACT(YEAR FROM ngay) = ?
			AND trang_thai = 'da_giao'
			AND trang_thai_thanh_toan = 'da_thanh_toan'
		`, thang, nam).
		Scan(&result).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Không thể tính doanh thu tháng"})
		return
	}

	result.Thang = thang
	result.Nam = nam

	c.JSON(200, gin.H{"data": result})
}

func GetDoanhThuTheoNam(c *gin.Context) {

	nam, _ := strconv.Atoi(c.DefaultQuery("nam", fmt.Sprint(time.Now().Year())))

	var result DoanhThuDTO

	err := config.DB.
		Model(&models.HoaDon{}).
		Select(`
			COALESCE(SUM(tong_tien), 0) AS doanh_thu,
			COUNT(ma_hoa_don) AS so_don
		`).
		Where(`
			EXTRACT(YEAR FROM ngay) = ?
			AND trang_thai = 'da_giao'
			AND trang_thai_thanh_toan = 'da_thanh_toan'
		`, nam).
		Scan(&result).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Không thể tính doanh thu năm"})
		return
	}

	result.Nam = nam
	c.JSON(200, gin.H{"data": result})
}
func ExportDoanhThuNgay(c *gin.Context) {

	ngay := c.Query("ngay")
	if ngay == "" {
		ngay = time.Now().Format("2006-01-02")
	}

	var hoaDons []models.HoaDon

	config.DB.
		Preload("ChiTietHoaDons.MonAn").
		Preload("ChiTietHoaDons.Options.OptionItem").
		Where("DATE(ngay)=? AND trang_thai=? AND trang_thai_thanh_toan=?",
			ngay, "da_giao", "da_thanh_toan").
		Find(&hoaDons)

	f := excelize.NewFile()
	sheet := "DoanhThu"
	f.NewSheet(sheet)

	// ======================
	// STYLE
	// ======================
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 18},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D9E1F2"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	// ======================
	// TITLE
	// ======================
	f.SetCellValue(sheet, "A1", "BÁO CÁO DOANH THU NGÀY "+ngay)
	f.MergeCell(sheet, "A1", "F1")
	f.SetCellStyle(sheet, "A1", "F1", titleStyle)

	row := 3
	var grandTotal float64

	// ======================
	// HEADER
	// ======================
	headers := []string{"MAHD", "Họ tên", "SĐT", "Tạm tính", "Giảm", "Tổng"}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, row)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	row++

	// ======================
	// DATA
	// ======================
	for _, hd := range hoaDons {

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), hd.MaHoaDon)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), hd.HoTen)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), hd.SDT)

		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), formatMoneyVND(hd.TamTinh))
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), formatMoneyVND(hd.TienGiam))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), formatMoneyVND(hd.TongTien))

		row++

		// ======================
		// CHI TIẾT MÓN
		// ======================
		for _, ct := range hd.ChiTietHoaDons {

			f.SetCellValue(sheet, fmt.Sprintf("B%d", row),
				"   "+ct.MonAn.TenMonAn)

			f.SetCellValue(sheet, fmt.Sprintf("C%d", row),
				fmt.Sprintf("SL: %d", ct.SoLuong))

			f.SetCellValue(sheet, fmt.Sprintf("F%d", row),
				formatMoneyVND(ct.ThanhTien))

			row++

			// OPTIONS
			for _, op := range ct.Options {

				name := op.TenOption
				if name == "" {
					name = op.OptionItem.TenOption
				}

				f.SetCellValue(sheet, fmt.Sprintf("B%d", row),
					"      + "+name)

				f.SetCellValue(sheet, fmt.Sprintf("F%d", row),
					formatMoneyVND(op.GiaThem))

				row++
			}
		}

		grandTotal += hd.TongTien
		row++
	}

	// ======================
	// TOTAL
	// ======================
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), "TỔNG DOANH THU:")
	f.SetCellValue(sheet, fmt.Sprintf("F%d", row), formatMoneyVND(grandTotal))

	// ======================
	// COLUMN WIDTH
	// ======================
	f.SetColWidth(sheet, "A", "A", 10)
	f.SetColWidth(sheet, "B", "B", 35)
	f.SetColWidth(sheet, "C", "C", 15)
	f.SetColWidth(sheet, "D", "F", 18)

	f.DeleteSheet("Sheet1")

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=doanh_thu.xlsx")

	_ = f.Write(c.Writer)
}

func GetDanhSachNgayDoanhThu(c *gin.Context) {

	var days []string

	err := config.DB.
		Model(&models.HoaDon{}).
		Select("DISTINCT CAST(ngay AS DATE)").
		Where("trang_thai = ? AND trang_thai_thanh_toan = ?", "da_giao", "da_thanh_toan").
		Order("ngay DESC").
		Pluck("ngay", &days).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "fail"})
		return
	}

	c.JSON(200, gin.H{
		"data": days,
	})
}

func TopMonBanChay(c *gin.Context) {
	type Result struct {
		TenMonAn string `json:"ten_mon_an"`
		TongBan  int    `json:"tong_ban"`
	}

	var result []Result

	err := config.DB.Raw(`
		SELECT 
			m.ten_mon_an,
			SUM(ct.so_luong) AS tong_ban
		FROM chi_tiet_hoa_dons ct
		JOIN mon_ans m ON m.ma_mon_an = ct.ma_mon_an
		JOIN hoa_dons h ON h.ma_hoa_don = ct.ma_hoa_don
		WHERE h.trang_thai = 'da_giao'
		  AND h.trang_thai_thanh_toan = 'da_thanh_toan'
		GROUP BY m.ten_mon_an
		ORDER BY tong_ban DESC
		LIMIT 9
	`).Scan(&result).Error

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Lỗi lấy top món bán chạy",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, result)
}

func GetTiLeHoanThanhHomNay(c *gin.Context) {
	today := time.Now().Format("2006-01-02")

	var tongDon int64
	var donHoanThanh int64

	config.DB.Model(&models.HoaDon{}).
		Where("CAST(ngay AS DATE) = ?", today).
		Count(&tongDon)

	config.DB.Model(&models.HoaDon{}).
		Where(`
			CAST(ngay AS DATE) = ?
			AND trang_thai = 'da_giao'
			AND trang_thai_thanh_toan = 'da_thanh_toan'
		`, today).
		Count(&donHoanThanh)

	tiLe := 0.0
	if tongDon > 0 {
		tiLe = float64(donHoanThanh) / float64(tongDon) * 100
	}

	c.JSON(200, gin.H{
		"data": gin.H{
			"tong_don":   tongDon,
			"hoan_thanh": donHoanThanh,
			"ti_le":      math.Round(tiLe),
		},
	})
}

func GetTopMonAnBanChay(c *gin.Context) {

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	var result []MonAnBanChayDTO

	err := config.DB.Raw(`
	SELECT 
		ma.ma_mon_an,
		ma.ten_mon_an,
		SUM(cthd.so_luong) AS so_luong
	FROM chi_tiet_hoa_dons cthd
	JOIN hoa_dons hd ON hd.ma_hoa_don = cthd.ma_hoa_don
	JOIN mon_ans ma ON ma.ma_mon_an = cthd.ma_mon_an
	WHERE hd.trang_thai = 'da_giao'
	  AND hd.trang_thai_thanh_toan = 'da_thanh_toan'
	GROUP BY ma.ma_mon_an, ma.ten_mon_an
	ORDER BY so_luong DESC
	LIMIT ?
`, limit).Scan(&result).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Không thể lấy top món ăn"})
		return
	}

	c.JSON(200, gin.H{"data": result})
}

func formatMoneyVND(n float64) string {
	// ép về int trước cho an toàn
	v := int64(n)

	s := fmt.Sprintf("%d", v)

	// format thủ công dấu chấm
	nStr := ""
	for i, c := range reverse(s) {
		if i != 0 && i%3 == 0 {
			nStr += "."
		}
		nStr += string(c)
	}

	return reverse(nStr) + " ₫"
}
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}