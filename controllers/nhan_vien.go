package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/dto"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
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
	if nv.LoaiNguoiDung != "user" && nv.LoaiNguoiDung != "admin" && nv.LoaiNguoiDung != "shipper" {
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
	trangThai := c.PostForm("trang_thai")
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
	if trangThai != "" {
		nv.TrangThai = trangThai
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

	// Tìm nhân viên
	if err := config.DB.First(&nv, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy nhân viên",
		})
		return
	}

	// Cập nhật trạng thái thành khóa
	if err := config.DB.Model(&nv).Update("trang_thai", "khoa").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đã khóa tài khoản nhân viên",
	})
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
func GetShippers(c *gin.Context) {
	var shippers []models.NguoiDung

	err := config.DB.
		Preload("AnhNhanVien").
		Where("loai_nguoi_dung = ?", "shipper").
		Find(&shippers).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Lỗi lấy danh sách shipper",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, shippers)
}

func AssignShipper(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {

		var input struct {
			MaHoaDon  uint `json:"ma_hoa_don"`
			MaShipper uint `json:"ma_shipper"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "invalid input"})
			return
		}

		var hoaDon models.HoaDon

		// 1️⃣ kiểm tra hóa đơn
		if err := config.DB.
			Where("ma_hoa_don = ?", input.MaHoaDon).
			First(&hoaDon).Error; err != nil {

			c.JSON(404, gin.H{"error": "Hóa đơn không tồn tại"})
			return
		}

		// 2️⃣ chỉ assign khi đã xác nhận
		if hoaDon.TrangThai != "da_xac_nhan" {
			c.JSON(400, gin.H{"error": "Hóa đơn chưa sẵn sàng giao"})
			return
		}

		// 3️⃣ update shipper + trạng thái
		if err := config.DB.Model(&hoaDon).Updates(map[string]interface{}{
			"ma_shipper": input.MaShipper,
			"trang_thai": "da_giao_shipper",
		}).Error; err != nil {

			c.JSON(500, gin.H{"error": "Không thể assign shipper"})
			return
		}

		// 4️⃣ LOAD LẠI FULL DATA (RẤT QUAN TRỌNG)
		if err := config.DB.
			Preload("Shipper").
			Preload("ChiTietHoaDons").
			Preload("ChiTietHoaDons.MonAn").
			Preload("ChiTietHoaDons.Options").
			Preload("ChiTietHoaDons.Options.OptionItem").
			Preload("ChiTietHoaDons.Options.OptionItem.NhomOption").
			First(&hoaDon, hoaDon.MaHoaDon).Error; err != nil {

			c.JSON(500, gin.H{"error": "Không load được chi tiết hóa đơn"})
			return
		}

		// 5️⃣ REALTIME CHO SHIPPER
		hub.SendToUser(input.MaShipper, dto.WSMessage{
			Type:    "shipper_new_order",
			Payload: hoaDon,
		})

		//realtime cho user
		hub.SendToUser(hoaDon.MaNguoiDung, dto.WSMessage{
			Type: "assign_shipper_user",
			Payload: gin.H{
				"ma_hoa_don": hoaDon.MaHoaDon,
				"shipper":    hoaDon.Shipper, // ⭐ chỉ cần shipper
			},
		})

		// 6️⃣ REALTIME CHO ADMIN (cập nhật list)
		hub.Broadcast(dto.WSMessage{
			Type:    "admin_assign_shipper",
			Payload: hoaDon,
		})

		c.JSON(200, gin.H{
			"message": "Assign shipper thành công",
			"data":    hoaDon,
		})
	}
}

// func NotifyShipper(hub *websocket.Hub, shipperID uint, ship models.ShipOrder) {

// 	msg := dto.WSMessage{
// 		Type:    "NEW_ORDER",
// 		Payload: ship,
// 	}

// 	hub.SendToUser(shipperID, msg)
// }

// func AcceptShipOrder(c *gin.Context) {
// 	var input struct {
// 		ShipOrderID uint `json:"ship_order_id"`
// 	}

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(400, gin.H{"error": "invalid input"})
// 		return
// 	}

// 	var ship models.ShipOrder

// 	// 1. lấy đơn
// 	if err := config.DB.First(&ship, input.ShipOrderID).Error; err != nil {
// 		c.JSON(404, gin.H{"error": "not found"})
// 		return
// 	}

// 	// 2. check trạng thái
// 	if ship.Status != "pending" {
// 		c.JSON(400, gin.H{"error": "order already taken"})
// 		return
// 	}

// 	now := time.Now()

// 	// 3. update ship order
// 	config.DB.Model(&ship).Updates(map[string]interface{}{
// 		"status":      "accepted",
// 		"accepted_at": now,
// 	})

// 	// 4. update hóa đơn
// 	config.DB.Model(&models.HoaDon{}).
// 		Where("ma_hoa_don = ?", ship.MaHoaDon).
// 		Updates(map[string]interface{}{
// 			"ma_shipper": ship.MaShipper,
// 			"trang_thai": "dang_giao",
// 		})

// 	c.JSON(200, gin.H{
// 		"message": "Đã xác nhận đơn hàng",
// 	})
// }
