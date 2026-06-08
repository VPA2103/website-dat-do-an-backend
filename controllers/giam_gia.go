package controllers

import (
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/dto"
)

type GiamGiaInput struct {
	Code           string  `form:"code"`
	TenChuongTrinh string  `form:"ten_chuong_trinh"`
	LoaiGiamGia    string  `form:"loai_giam_gia"`
	GiaTriGiam     float64 `form:"gia_tri_giam"`
	DonToiThieu    float64 `form:"don_toi_thieu"`
	GiamToiDa      float64 `form:"giam_toi_da"`

	GioiHanSuDung *int `form:"gioi_han_su_dung"`

	NgayBatDau  string `form:"ngay_bat_dau"`
	NgayKetThuc string `form:"ngay_ket_thuc"`

	IsActive bool `form:"is_active"`
}

func CreateGiamGia(c *gin.Context) {
	var input GiamGiaInput

	// bind form-data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ: " + err.Error(),
		})
		return
	}

	// validate loại giảm giá
	if input.LoaiGiamGia != "percent" && input.LoaiGiamGia != "fixed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "loai_giam_gia phải là percent hoặc fixed",
		})
		return
	}

	// parse time
	ngayBatDau, err := time.Parse(time.RFC3339, input.NgayBatDau)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ngày bắt đầu không hợp lệ",
		})
		return
	}

	ngayKetThuc, err := time.Parse(time.RFC3339, input.NgayKetThuc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ngày kết thúc không hợp lệ",
		})
		return
	}

	giamGia := dto.GiamGia{
		Code:           input.Code,
		TenChuongTrinh: input.TenChuongTrinh,
		LoaiGiamGia:    input.LoaiGiamGia,
		GiaTriGiam:     input.GiaTriGiam,
		DonToiThieu:    input.DonToiThieu,
		GiamToiDa:      input.GiamToiDa,
		GioiHanSuDung:  input.GioiHanSuDung,
		NgayBatDau:     ngayBatDau,
		NgayKetThuc:    ngayKetThuc,
		IsActive:       input.IsActive,
	}

	// lưu trước để có ID
	if err := config.DB.Create(&giamGia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể tạo mã giảm giá",
		})
		return
	}

	// upload nhiều ảnh
	form, _ := c.MultipartForm()

	if form != nil {
		files := form.File["images"]

		for _, file := range files {

			src, err := file.Open()
			if err != nil {
				continue
			}

			uploadResult, err := config.CLD.Upload.Upload(
				c,
				src,
				uploader.UploadParams{
					Folder: "giamgia",
				},
			)

			src.Close()

			if err != nil {
				continue
			}

			image := dto.HinhAnh{
				Url:       uploadResult.SecureURL,
				OwnerID:   giamGia.ID,
				OwnerType: "giam_gia",
			}

			config.DB.Create(&image)
		}
	}

	// preload ảnh
	config.DB.Preload("AnhGiamGia").First(&giamGia, giamGia.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo mã giảm giá thành công",
		"data":    giamGia,
	})
}

func GetAllGiamGia(c *gin.Context) {
	var giamGia []dto.GiamGia

	if err := config.DB.Find(&giamGia).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Không thể lấy danh sách",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": giamGia,
	})
}

func GetGiamGiaById(c *gin.Context) {
	id := c.Param("id")

	var giamGia dto.GiamGia

	if err := config.DB.First(&giamGia, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy mã giảm giá",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": giamGia,
	})
}

func UpdateGiamGia(c *gin.Context) {
	id := c.Param("id")

	var giamGia dto.GiamGia

	if err := config.DB.Preload("AnhGiamGia").
		First(&giamGia, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy mã giảm giá",
		})
		return
	}

	var input GiamGiaInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu không hợp lệ",
		})
		return
	}

	// parse date
	ngayBatDau, _ := time.Parse(time.RFC3339, input.NgayBatDau)
	ngayKetThuc, _ := time.Parse(time.RFC3339, input.NgayKetThuc)

	giamGia.Code = input.Code
	giamGia.TenChuongTrinh = input.TenChuongTrinh
	giamGia.LoaiGiamGia = input.LoaiGiamGia
	giamGia.GiaTriGiam = input.GiaTriGiam
	giamGia.DonToiThieu = input.DonToiThieu
	giamGia.GiamToiDa = input.GiamToiDa
	giamGia.GioiHanSuDung = input.GioiHanSuDung
	giamGia.NgayBatDau = ngayBatDau
	giamGia.NgayKetThuc = ngayKetThuc
	giamGia.IsActive = input.IsActive

	if err := config.DB.Save(&giamGia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể cập nhật",
		})
		return
	}

	// upload thêm ảnh mới
	form, _ := c.MultipartForm()

	if form != nil {
		files := form.File["images"]

		for _, file := range files {

			src, err := file.Open()
			if err != nil {
				continue
			}

			uploadResult, err := config.CLD.Upload.Upload(
				c,
				src,
				uploader.UploadParams{
					Folder: "giamgia",
				},
			)

			src.Close()

			if err != nil {
				continue
			}

			image := dto.HinhAnh{
				Url:       uploadResult.SecureURL,
				OwnerID:   giamGia.ID,
				OwnerType: "giam_gia",
			}

			config.DB.Create(&image)
		}
	}

	config.DB.Preload("AnhGiamGia").First(&giamGia, giamGia.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật thành công",
		"data":    giamGia,
	})
}

// DELETE
func DeleteGiamGia(c *gin.Context) {
	id := c.Param("id")

	var giamGia dto.GiamGia

	if err := config.DB.First(&giamGia, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy mã giảm giá",
		})
		return
	}

	if err := config.DB.Delete(&giamGia).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Không thể xoá",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xoá thành công",
	})
}
