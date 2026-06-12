package controllers

import (
	//"fmt"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	//"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	
)

type MonAnChiTietResponse struct {
	MonAn     models.MonAn      `json:"mon_an"`
	DanhGias  []models.DanhGia  `json:"danh_gias"`
	BinhLuans []models.BinhLuan `json:"binh_luans"`
}

func (h *ChatHandler) CreateMonAn(c *gin.Context) {
	var monan models.MonAn

	// 1. bind
	if err := c.ShouldBind(&monan); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 2. validate
	if monan.TenMonAn == "" || monan.MoTa == "" {
		c.JSON(400, gin.H{"error": "Thiếu dữ liệu"})
		return
	}

	// 3. save món ăn
	if err := config.DB.Create(&monan).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// =========================
	// 🧠 BUILD DOCUMENT
	// =========================
	document := fmt.Sprintf(
		"Món: %s\nMô tả: %s\nGiá: %.0f\nLoại: %d",
		monan.TenMonAn,
		monan.MoTa,
		monan.GiaTien,
		monan.MaLoaiMonAn,
	)

	// monan.Document = document
	// monan.SearchText = monan.TenMonAn + " " + monan.MoTa
	// monan.HasEmbedding = false
	config.DB.Save(&monan)

	// =========================
	// 🔥 EMBEDDING
	// =========================
	embedding, err := h.llm.Embed(c.Request.Context(), document)
	if err != nil {
		log.Println("embed error:", err)
	}

	// =========================
	// 📦 METADATA → JSON
	// =========================
	metaJSON, _ := json.Marshal(map[string]any{
		"id":    monan.MaMonAn,
		"name":  monan.TenMonAn,
		"price": monan.GiaTien,
	})

	// =========================
	// 📦 VECTOR STRING FORMAT
	// =========================
	vectorStr := vectorToString(embedding)

	// =========================
	// INSERT VECTOR DB
	// =========================

	if len(embedding) > 0 {

		embeddingID := fmt.Sprintf("mon_an_%d", monan.MaMonAn)

		result := config.DB.Exec(`
		INSERT INTO menu_embeddings (id, document, metadata, embedding)
		VALUES ($1, $2, $3, $4)
	`,
			embeddingID, // 👈 FIX Ở ĐÂY
			document,
			string(metaJSON),
			vectorStr,
		)

		if result.Error != nil {
			log.Println("vector insert error:", result.Error)
		}
	}

	// =========================
	// 🖼️ UPLOAD IMAGE
	// =========================
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()

			uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
				Folder: "monan",
			})

			if err == nil {
				img := models.HinhAnh{
					OwnerID:   monan.MaMonAn,
					OwnerType: "mon_an",
					Url:       uploadResult.SecureURL,
				}
				config.DB.Create(&img)
			}
		}
	}

	// =========================
	// RESPONSE
	// =========================
	config.DB.Preload("AnhMonAn").First(&monan, monan.MaMonAn)

	c.JSON(201, gin.H{
		"message": "Tạo món ăn + embedding thành công",
		"data":    monan,
	})
}

// ======================= GET ALL =======================
func GetAllMonAn(c *gin.Context) {
	var list []models.MonAn
	config.DB.Preload("AnhMonAn").Find(&list)

	c.JSON(http.StatusOK, gin.H{"data": list})
}

func GetMonAnByID(c *gin.Context) {
	id := c.Param("id")
	var monan models.MonAn

	if err := config.DB.
		Preload("AnhMonAn").
		Preload("NhomOptions").
		Preload("NhomOptions.OptionItems").First(&monan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy món ăn"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": monan})
}

// ======================= UPDATE =======================
func (h *ChatHandler) UpdateMonAn(c *gin.Context) {
	id := c.Param("id")
	var monan models.MonAn

	// 1️⃣ Find
	if err := config.DB.First(&monan, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy món ăn"})
		return
	}

	// 2️⃣ Bind
	var input models.MonAn
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if input.MoTa == "" {
		c.JSON(400, gin.H{"error": "Mô tả món ăn không được để trống"})
		return
	}

	// 3️⃣ Update DOMAIN
	config.DB.Model(&monan).Updates(map[string]any{
		"ten_mon_an":     input.TenMonAn,
		"mo_ta":          input.MoTa,
		"gia_tien":       input.GiaTien,
		"ma_loai_mon_an": input.MaLoaiMonAn,
		"trang_thai":     input.TrangThai,
	})

	// =========================
	// 🧠 BUILD DOCUMENT (RAG)
	// =========================
	document := fmt.Sprintf(
		"Món: %s\nMô tả: %s\nGiá: %.0f\nLoại: %d",
		input.TenMonAn,
		input.MoTa,
		input.GiaTien,
		input.MaLoaiMonAn,
	)

	// =========================
	// 🔥 EMBEDDING
	// =========================
	embedding, err := h.llm.Embed(c.Request.Context(), document)
	if err != nil {
		log.Println("embed error:", err)
	}

	// =========================
	// 📦 METADATA
	// =========================
	metaJSON, _ := json.Marshal(map[string]any{
		"type":  "mon_an",
		"id":    monan.MaMonAn,
		"name":  input.TenMonAn,
		"price": input.GiaTien,
	})

	// =========================
	// 📦 VECTOR STRING
	// =========================
	vectorStr := vectorToString(embedding)

	// =========================
	// UPDATE VECTOR DB
	// =========================
	if len(embedding) > 0 {

		embeddingID := fmt.Sprintf("mon_an_%d", monan.MaMonAn)

		result := config.DB.Exec(`
		UPDATE menu_embeddings
		SET document = $1,
		    metadata = $2,
		    embedding = $3
		WHERE id = $4
	`,
			document,
			string(metaJSON),
			vectorStr,
			embeddingID, // 👈 FIX Ở ĐÂY
		)

		if result.Error != nil {
			log.Println("vector update error:", result.Error)
		}
	}

	// =========================
	// 🖼️ UPLOAD IMAGE
	// =========================
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, _ := file.Open()
		defer src.Close()

		upload, _ := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
			Folder: "monan",
		})

		config.DB.Where(
			"owner_id = ? AND owner_type = ?",
			monan.MaMonAn,
			"mon_an",
		).Delete(&models.HinhAnh{})

		config.DB.Create(&models.HinhAnh{
			Url:       upload.SecureURL,
			OwnerID:   monan.MaMonAn,
			OwnerType: "mon_an",
		})
	}

	// =========================
	// RESPONSE
	// =========================
	config.DB.Preload("AnhMonAn").First(&monan, monan.MaMonAn)

	c.JSON(200, gin.H{
		"message": "Cập nhật món ăn + embedding thành công",
		"data":    monan,
	})
}

// ======================= DELETE =======================
func DeleteMonAn(c *gin.Context) {
	id := c.Param("id")
	var monan models.MonAn

	// 1️⃣ Find món ăn
	if err := config.DB.First(&monan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy món ăn"})
		return
	}

	// =========================
	// 🗑️ DELETE EMBEDDING (RAG)
	// =========================
	embeddingID := fmt.Sprintf("mon_an_%d", monan.MaMonAn)

	if err := config.DB.Exec(
		`DELETE FROM menu_embeddings WHERE id = $1`,
		embeddingID,
	).Error; err != nil {
		log.Println("delete embedding error:", err)
	}

	// =========================
	// 🗑️ DELETE IMAGE
	// =========================
	config.DB.Where(
		"owner_id = ? AND owner_type = ?",
		monan.MaMonAn,
		"mon_an",
	).Delete(&models.HinhAnh{})

	// =========================
	// 🗑️ DELETE DOMAIN
	// =========================
	config.DB.Delete(&monan)

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa món ăn + embedding thành công",
	})
}

func GetMonAnDetail(c *gin.Context) {

	id := c.Param("id")

	var monan models.MonAn

	err := config.DB.
		Preload("AnhMonAn").
		Preload("NhomOptions").
		Preload("NhomOptions.OptionItems").
		First(&monan, id).Error

	if err != nil {
		c.JSON(404, gin.H{
			"error": "Không tìm thấy món ăn",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": monan,
	})
}

func vectorToString(vec []float32) string {
	if len(vec) == 0 {
		return "[]"
	}

	var b strings.Builder
	b.WriteString("[")

	for i, v := range vec {
		b.WriteString(fmt.Sprintf("%f", v))
		if i < len(vec)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("]")
	return b.String()
}

func SearchMonAn(c *gin.Context) {
	keyword := c.Query("q")

	var list []models.MonAn

	query := config.DB.
		Preload("AnhMonAn").
		Where("trang_thai = ?", 1)

	if keyword != "" {
		query = query.Where(
			"LOWER(ten_mon_an) LIKE LOWER(?) OR LOWER(mo_ta) LIKE LOWER(?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
	}

	if err := query.Find(&list).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": list,
	})
}

func GetMonAnCoBinhLuanVaDanhGia(c *gin.Context) {
	var monAns []models.MonAn

	err := config.DB.
		Preload("AnhMonAn").
		Find(&monAns).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Lỗi truy vấn"})
		return
	}

	type MonAnQuanLy struct {
		MaMonAn     uint        `json:"ma_mon_an"`
		TenMonAn    string      `json:"ten_mon_an"`
		AnhMonAn    interface{} `json:"anh_mon_an"`
		SoBinhLuan  int64       `json:"so_binh_luan"`
		SoDanhGia   int64       `json:"so_danh_gia"`
	}

	var result []MonAnQuanLy

	for _, m := range monAns {
		var soBL int64
		var soDG int64

		config.DB.Model(&models.BinhLuan{}).
			Where("ma_mon_an = ?", m.MaMonAn).
			Count(&soBL)

		config.DB.Model(&models.DanhGia{}).
			Where("ma_mon_an = ?", m.MaMonAn).
			Count(&soDG)

		if soBL > 0 || soDG > 0 {
			result = append(result, MonAnQuanLy{
				MaMonAn:    m.MaMonAn,
				TenMonAn:   m.TenMonAn,
				AnhMonAn:   m.AnhMonAn,
				SoBinhLuan: soBL,
				SoDanhGia:  soDG,
			})
		}
	}

	c.JSON(200, gin.H{
		"data": result, // ✅ ARRAY
	})
}
func GetMonAnCoBinhLuanVaDanhGiaCuaNguoiDung(c *gin.Context) {

	// 🔐 Lấy user đăng nhập
	maNguoiDungAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Chưa đăng nhập"})
		return
	}
	maNguoiDung := maNguoiDungAny.(uint)

	var monAns []models.MonAn

	err := config.DB.
		Preload("AnhMonAn").
		Find(&monAns).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Lỗi truy vấn"})
		return
	}

	type MonAnQuanLy struct {
		MaMonAn     uint        `json:"ma_mon_an"`
		TenMonAn    string      `json:"ten_mon_an"`
		AnhMonAn    interface{} `json:"anh_mon_an"`
		SoBinhLuan  int64       `json:"so_binh_luan"`
		SoDanhGia   int64       `json:"so_danh_gia"`
	}

	var result []MonAnQuanLy

	for _, m := range monAns {
		var soBL int64
		var soDG int64

		// 🗨️ Bình luận của user
		config.DB.Model(&models.BinhLuan{}).
			Where("ma_mon_an = ? AND ma_nguoi_dung = ?", m.MaMonAn, maNguoiDung).
			Count(&soBL)

		// ⭐ Đánh giá của user
		config.DB.Model(&models.DanhGia{}).
			Where("ma_mon_an = ? AND ma_nguoi_dung = ?", m.MaMonAn, maNguoiDung).
			Count(&soDG)

		// ✅ Chỉ lấy món user đã từng tương tác
		if soBL > 0 || soDG > 0 {
			result = append(result, MonAnQuanLy{
				MaMonAn:    m.MaMonAn,
				TenMonAn:   m.TenMonAn,
				AnhMonAn:   m.AnhMonAn,
				SoBinhLuan: soBL,
				SoDanhGia:  soDG,
			})
		}
	}

	c.JSON(200, gin.H{
		"data": result,
	})
}


