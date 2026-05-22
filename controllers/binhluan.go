package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/internal/dto"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/models"
)

type BinhLuanController struct {
	Hub *websocket.Hub
}

func NewBinhLuanController(hub *websocket.Hub) *BinhLuanController {
	return &BinhLuanController{Hub: hub}
}

type CreateBinhLuanInput struct {
	MaMonAn uint   `json:"ma_mon_an" binding:"required"`
	NoiDung string `json:"noi_dung" binding:"required"`

	// 👇 optional
	ParentID *uint `json:"parent_id"`
}

type UpdateBinhLuanInput struct {
	NoiDung string `json:"noi_dung"`
}

func (ctrl *BinhLuanController) CreateBinhLuan(c *gin.Context) {
	var input CreateBinhLuanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	maNguoiDungAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Không tìm thấy người dùng trong token"})
		return
	}
	maNguoiDung := maNguoiDungAny.(uint)

	var monAn models.MonAn
	if err := config.DB.First(&monAn, input.MaMonAn).Error; err != nil {
		c.JSON(404, gin.H{"error": "Món ăn không tồn tại"})
		return
	}

	binhLuan := models.BinhLuan{
		MaNguoiDung: maNguoiDung,
		MaMonAn:     input.MaMonAn,
		NoiDung:     input.NoiDung,
		ParentID:    input.ParentID,
	}

	if input.ParentID != nil {
		var parent models.BinhLuan

		if err := config.DB.First(&parent, *input.ParentID).Error; err != nil {
			c.JSON(404, gin.H{
				"error": "Bình luận cha không tồn tại",
			})
			return
		}
	}

	if err := config.DB.Create(&binhLuan).Error; err != nil {
		c.JSON(500, gin.H{"error": "Không thể tạo bình luận"})
		return
	}

	config.DB.Preload("NguoiDung").Preload("Replies.NguoiDung").First(&binhLuan, binhLuan.ID)

	// ✅ Broadcast realtime
	ctrl.Hub.Broadcast(dto.WSMessage{
		Type:    "new_binh_luan",
		Payload: binhLuan,
	})

	var comments []models.BinhLuan

	config.DB.
		Where("ma_mon_an = ? AND parent_id IS NULL", input.MaMonAn).
		Preload("NguoiDung").
		Preload("Replies").
		Preload("Replies.NguoiDung").
		Order("created_at desc").
		Find(&comments)

	c.JSON(200, gin.H{
		"message": "Tạo bình luận thành công",
		"data":    comments,
	})
}

func (ctrl *BinhLuanController) GetBinhLuanByMonAn(c *gin.Context) {
	maMon := c.Param("ma_mon_an")

	var binhLuans []models.BinhLuan

	err := config.DB.
		Where("ma_mon_an = ? AND parent_id IS NULL", maMon).
		Preload("NguoiDung").
		Preload("Replies").
		Preload("Replies.NguoiDung").
		Order("created_at desc").
		Find(&binhLuans).Error

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": binhLuans,
	})
}

func (ctrl *BinhLuanController) UpdateBinhLuan(c *gin.Context) {
	id := c.Param("id")

	maNguoiDungAny, _ := c.Get("user_id")
	maNguoiDung := maNguoiDungAny.(uint)

	var binhLuan models.BinhLuan
	if err := config.DB.First(&binhLuan, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy bình luận"})
		return
	}

	if binhLuan.MaNguoiDung != maNguoiDung {
		c.JSON(403, gin.H{"error": "Không có quyền sửa"})
		return
	}

	var input UpdateBinhLuanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	if err := config.DB.Model(&binhLuan).
		Update("noi_dung", input.NoiDung).Error; err != nil {
		c.JSON(500, gin.H{"error": "Không thể cập nhật"})
		return
	}

	config.DB.Preload("NguoiDung").First(&binhLuan, id)

	ctrl.Hub.Broadcast(dto.WSMessage{
		Type:    "update_binh_luan",
		Payload: binhLuan,
	})

	c.JSON(200, gin.H{"data": binhLuan})
}

func (ctrl *BinhLuanController) GetBinhLuanByID(c *gin.Context) {
	id := c.Param("id")

	maNguoiDungAny, _ := c.Get("user_id")
	maNguoiDung := maNguoiDungAny.(uint)

	var binhLuan models.BinhLuan
	if err := config.DB.First(&binhLuan, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy bình luận"})
		return
	}

	if binhLuan.MaNguoiDung != maNguoiDung {
		c.JSON(403, gin.H{"error": "Không có quyền sửa bình luận này"})
		return
	}

	var input UpdateBinhLuanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	if err := config.DB.Model(&binhLuan).Update("noi_dung", input.NoiDung).Error; err != nil {
		c.JSON(500, gin.H{"error": "Không thể cập nhật"})
		return
	}

	// ✅ Broadcast realtime
	ctrl.Hub.Broadcast( dto.WSMessage{
		Type:    "update_binh_luan",
		Payload: binhLuan,
	})

	c.JSON(200, gin.H{"data": binhLuan})
}

func (ctrl *BinhLuanController) DeleteBinhLuan(c *gin.Context) {
	id := c.Param("id")

	maNguoiDungAny, _ := c.Get("user_id")
	maNguoiDung := maNguoiDungAny.(uint)

	var binhLuan models.BinhLuan
	if err := config.DB.First(&binhLuan, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy bình luận"})
		return
	}

	if binhLuan.MaNguoiDung != maNguoiDung {
		c.JSON(403, gin.H{"error": "Không có quyền xóa bình luận này"})
		return
	}

	if err := config.DB.Delete(&binhLuan, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Không thể xóa"})
		return
	}

	// ✅ Broadcast realtime
	ctrl.Hub.Broadcast( dto.WSMessage{
		Type:    "delete_binh_luan",
		Payload: gin.H{"id": id},
	})

	c.JSON(200, gin.H{"message": "Đã xóa"})
}
