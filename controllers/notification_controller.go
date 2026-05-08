package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

func GetNotifications(c *gin.Context) {
	userID := c.Query("user_id")

	var notifications []models.ThongBao
	config.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications)

	c.JSON(http.StatusOK, notifications)
}

// Đánh dấu đã đọc
func MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	config.DB.Model(&models.ThongBao{}).
		Where("id = ?", id).
		Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{"message": "Đã đọc"})
}
