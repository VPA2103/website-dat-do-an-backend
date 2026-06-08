package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/dto"
	"github.com/vpa/quanlynhahang-backend/internal/usecase"
	"github.com/vpa/quanlynhahang-backend/models"
)

type ContactHandler struct {
	NotiUC *usecase.NotificationUseCase
}

func (h *ContactHandler) GuiLienHe(c *gin.Context) {
	var lienHe models.LienHe

	if err := c.ShouldBind(&lienHe); err != nil {
		fmt.Println("❌ Bind error:", err.Error())

		c.JSON(400, gin.H{
			"message": "Dữ liệu không hợp lệ",
			"error":   err.Error(),
		})
		return
	}

	if lienHe.HoTen == "" || lienHe.Email == "" || lienHe.TieuDe == "" || lienHe.NoiDung == "" || lienHe.SDT == "" {
		c.JSON(400, gin.H{"message": "Vui lòng nhập đầy đủ thông tin"})
		return
	}

	if err := config.DB.Create(&lienHe).Error; err != nil {
		c.JSON(500, gin.H{"message": "Lưu liên hệ thất bại"})
		return
	}
	// thông báo realtime
	// 2. Gửi thông báo qua usecase
	h.NotiUC.Notify(dto.WSMessage{
		Type:    "notify",
		Content: "Có khách vừa gửi liên hệ",
		Role:    "admin",
	})

	c.JSON(200, gin.H{
		"message": "Gửi thành công",
		"data":    lienHe,
	})
}

func AdminGetAllLienHe(c *gin.Context) {
	// 👉 Nếu bạn đã có middleware check admin
	// thì KHÔNG cần đoạn check quyền ở đây

	var danhSachLienHe []models.LienHe

	if err := config.DB.
		Order("ngay_tao DESC").
		Find(&danhSachLienHe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Không thể lấy danh sách liên hệ",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lấy danh sách liên hệ thành công",
		"data":    danhSachLienHe,
	})
}

func DeleteLienHe(c *gin.Context) {
	id := c.Param("id")
	var lienhe models.LienHe

	// Kiểm tra tồn tại
	if err := config.DB.First(&lienhe, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Không tìm thấy liên hệ",
		})
		return
	}

	// Xóa
	if err := config.DB.Delete(&lienhe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Xóa liên hệ thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa thành công",
	})
}
