package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
	"golang.org/x/crypto/bcrypt"
)

// 🧱 Thêm nhân viên
func CreateNhanVien(c *gin.Context) {
	var nv models.NguoiDung

	// ✅ Lấy dữ liệu từ form-data
	if err := c.ShouldBind(&nv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu form không hợp lệ: " + err.Error()})
		return
	}

	// ✅ Kiểm tra loại nhân viên chỉ được phép là "user" hoặc "admin"
	if nv.LoaiNguoiDung != "user" && nv.LoaiNguoiDung != "admin"&& nv.LoaiNguoiDung != "shipper" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loại nhân viên không hợp lệ. Chỉ chấp nhận 'user' hoặc 'admin'."})
		return
	}

	// ✅ Mặc định ngày vào làm
	if nv.NgayVaoLam == "" {
		nv.NgayVaoLam = time.Now().Format("02-01-2006 15:04:05")
	}

	// ✅ Kiểm tra mật khẩu
	if nv.MatKhau == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mật khẩu không được để trống"})
		return
	}

	// ✅ Hash mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(nv.MatKhau),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Không thể mã hóa mật khẩu",
		})
		return
	}

	// Gán lại mật khẩu đã mã hóa
	nv.MatKhau = string(hashedPassword)

	// ✅ Lưu nhân viên trước để có MaNV (ID)
	if err := config.DB.Create(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo nhân viên: " + err.Error()})
		return
	}

	// ✅ Upload ảnh (nếu có)
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()

			uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{
				Folder: "nhanvien",
			})
			if err == nil {
				img := models.HinhAnh{
					OwnerID:   nv.MaNguoiDung,
					OwnerType: "nguoi_dung",
					Url:       uploadResult.SecureURL,
				}
				config.DB.Create(&img)
			}
		}
	}

	// ✅ Preload ảnh khi trả về
	config.DB.Preload("AnhNhanVien").First(&nv, nv.MaNguoiDung)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo nhân viên thành công",
		"data":    nv,
	})
}

// 📋 Lấy danh sách nhân viên
func GetAllNhanVien(c *gin.Context) {
	var nhanViens []models.NguoiDung
	if err := config.DB.Preload("AnhNhanVien").Find(&nhanViens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nhanViens)
}

// 🔍 Lấy 1 nhân viên theo ID
func GetNhanVienByID(c *gin.Context) {
	id := c.Param("id")
	var nv models.NguoiDung
	if err := config.DB.Preload("AnhNhanVien").Find(&nv, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nv)
}

// ✏️ Cập nhật nhân viên
func UpdateNhanVien(c *gin.Context) {
	id := c.Param("id")
	var nv models.NguoiDung

	// Tìm nhân viên theo ID
	if err := config.DB.First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}

	matKhau := c.PostForm("mat_khau")
	gioiTinh := c.PostForm("gioi_tinh")
	hoTen := c.PostForm("ho_ten")
	ngaySinh := c.PostForm("ngay_sinh")
	sdt := c.PostForm("sdt")

	loaiNhanVien := c.PostForm("loai_nhan_vien")
	email := c.PostForm("email")

	// Cập nhật từng trường nếu có dữ liệu
	if matKhau != "" {

		// Mã hóa mật khẩu mới
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(matKhau), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Mã hóa mật khẩu thất bại",
			})
			return
		}

		nv.MatKhau = string(hashedPassword)
	}

	if hoTen != "" {
		nv.HoTen = hoTen
	}
	if ngaySinh != "" {
		nv.NgaySinh = ngaySinh
	}
	if sdt != "" {
		nv.SDT = sdt
	}

	if email != "" {
		nv.Email = email
	}
	if loaiNhanVien != "" {
		nv.LoaiNguoiDung = loaiNhanVien
	}
	if gioiTinh != "" {
		nv.GioiTinh = gioiTinh
	}

	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, _ := file.Open()
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{Folder: "nhanvien"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload ảnh thất bại: " + err.Error()})
			return
		}

		// Xóa ảnh cũ
		config.DB.Where("owner_id = ? AND owner_type = ?", nv.MaNguoiDung, "nguoi_dung").Delete(&models.HinhAnh{})

		// Lưu ảnh mới
		newImg := models.HinhAnh{
			OwnerID:   nv.MaNguoiDung,
			OwnerType: "nguoi_dung",
			Url:       uploadResult.SecureURL,
		}
		config.DB.Create(&newImg)
	}

	// ✅ Lưu thay đổi
	if err := config.DB.Save(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật nhân viên: " + err.Error()})
		return
	}

	// ✅ Lấy lại thông tin mới
	// Trả về kết quả
	config.DB.Preload("AnhNhanVien").First(&nv, nv.MaNguoiDung)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật thành công",
		"data":    nv,
	})
}

// 🗑️ Xóa nhân viên
func DeleteNhanVien(c *gin.Context) {
	id := c.Param("id")
	var nv models.NguoiDung
	if err := config.DB.First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}

	if err := config.DB.Delete(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã xóa nhân viên thành công"})
}

func UpdateThongTinCaNhan(c *gin.Context) {
	id := c.Param("id")

	// ✅ Lấy user hiện tại từ middleware (Auth)
	currentUserID := c.GetUint("user_id")
	currentRole := c.GetString("role")

	// ✅ Nếu không phải admin và ID khác chính mình → cấm
	if currentRole != "admin" && fmt.Sprint(currentUserID) != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền chỉnh sửa thông tin người khác"})
		return
	}

	var nv models.NguoiDung
	if err := config.DB.
		Preload("AnhNhanVien").
		Preload("DiaChis").
		First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy nhân viên"})
		return
	}

	// ✅ Lấy dữ liệu form
	hoTen := c.PostForm("ho_ten")
	gioiTinh := c.PostForm("gioi_tinh")
	ngaySinh := c.PostForm("ngay_sinh")
	sdt := c.PostForm("sdt")

	email := c.PostForm("email")

	oldPassword := c.PostForm("mat_khau_cu")
	newPassword := c.PostForm("mat_khau_moi")
	confirmPassword := c.PostForm("xac_nhan_mat_khau_moi")

	// ✅ Cập nhật thông tin cơ bản
	if hoTen != "" {
		nv.HoTen = hoTen
	}
	if gioiTinh != "" {
		nv.GioiTinh = gioiTinh
	}
	if ngaySinh != "" {
		nv.NgaySinh = ngaySinh
	}
	if sdt != "" {
		nv.SDT = sdt
	}

	if email != "" {
		nv.Email = email
	}

	// ✅ Đổi mật khẩu (nếu có nhập đủ 3 trường)
	if oldPassword != "" || newPassword != "" || confirmPassword != "" {
		if oldPassword == "" || newPassword == "" || confirmPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cần nhập đủ mật khẩu cũ, mật khẩu mới và xác nhận mật khẩu mới"})
			return
		}

		// Chỉ người tự đổi mật khẩu mới cần check password cũ
		if currentRole != "admin" {
			if err := bcrypt.CompareHashAndPassword([]byte(nv.MatKhau), []byte(oldPassword)); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Mật khẩu cũ không đúng"})
				return
			}
		}

		if newPassword != confirmPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Xác nhận mật khẩu mới không khớp"})
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		nv.MatKhau = string(hashedPassword)
	}

	// ✅ Upload ảnh (nếu có)
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		src, _ := file.Open()
		defer src.Close()

		uploadResult, err := config.CLD.Upload.Upload(c, src, uploader.UploadParams{Folder: "nhanvien"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload ảnh thất bại: " + err.Error()})
			return
		}

		config.DB.Where("owner_id = ? AND owner_type = ?", nv.MaNguoiDung, "nguoi_dung").Delete(&models.HinhAnh{})

		newImg := models.HinhAnh{
			OwnerID:   nv.MaNguoiDung,
			OwnerType: "nguoi_dung",
			Url:       uploadResult.SecureURL,
		}
		config.DB.Create(&newImg)
	}

	// ✅ Lưu thay đổi
	if err := config.DB.Save(&nv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật thông tin cá nhân: " + err.Error()})
		return
	}

	config.DB.Preload("AnhNhanVien").First(&nv, nv.MaNguoiDung)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật thông tin cá nhân thành công",
		"data":    nv,
	})
}
