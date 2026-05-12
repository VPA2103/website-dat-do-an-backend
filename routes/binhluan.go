package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/controllers"
	"github.com/vpa/quanlynhahang-backend/internal/websocket"
	"github.com/vpa/quanlynhahang-backend/middleware"
)

func BinhLuanRoutes(r *gin.Engine, hub *websocket.Hub) {
	ctrl := controllers.NewBinhLuanController(hub)

	binhluan := r.Group("/binh-luan")
	{
		binhluan.POST("", middleware.AuthMiddleware(), ctrl.CreateBinhLuan)
		binhluan.GET("/mon-an/:ma_mon_an", ctrl.GetBinhLuanByMonAn)
		binhluan.GET("/:id", middleware.AuthMiddleware(), ctrl.GetBinhLuanByID)
		binhluan.PUT("/:id", middleware.AuthMiddleware(), ctrl.UpdateBinhLuan)
		binhluan.DELETE("/:id", middleware.AuthMiddleware(), ctrl.DeleteBinhLuan)
	}
}
